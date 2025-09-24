package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "modernc.org/sqlite"
)

// Here starts the AUTH functions
// ----------------------PRIVATE KEYS----------------------
var jwt_access_key = []byte("supersecretaccesskey")
var jwt_refresh_key = []byte("supersecretrefreshkey")

// --------------------------------------------------------

type Server struct {
	DB *sql.DB
}

// open/create SQLite and ensure schema
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:app.db?_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	schema := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return db, nil
}

func generateToken(username string, userid string) (string, string, error) {
	// Purpose: to generate a new pair of access and refresh tokens
	// Arguments: username: string (account username),
	// 			  userid: string (account id in SQL database)
	// Return: access_token: string (access key to store in browser local storage)
	//		   refresh_token: string (this will get stored in the http cookies)
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * 7).Unix(), // expires in 7 minutes
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24 * 10).Unix(), // expires in 10 days
	})

	sign_access, err1 := access_token.SignedString(jwt_access_key)
	sign_refresh, err2 := refresh_token.SignedString(jwt_refresh_key)

	if err1 != nil {
		return "", "", err1
	}
	if err2 != nil {
		return "", "", err2
	}

	return sign_access, sign_refresh, err2
}

func verifyToken(tokenString string, secretKey []byte) error {
	// Purpose: to verfiy jwt tokens
	// Arguments: tokenString: string (token to verify),
	// 			  secretKey: string (the key of the token to verify)
	// Return: if the token is valid this function will return nil
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ensure it's really HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// Endpoint functions here

func (s *Server) Refresh(c *gin.Context) {
	// Method: POST

	refresh_cookie, err1 := c.Cookie("refresh_token")
	if err1 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "no cookie found!"})
		return
	}
	fmt.Println(refresh_cookie)
	// verify signature
	err2 := verifyToken(refresh_cookie, jwt_refresh_key)
	if err2 != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid cookie"})
		return
	}

	// parse claims to extract uid/username
	tok, err := jwt.Parse(refresh_cookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwt_refresh_key, nil
	})
	if err != nil || !tok.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid cookie"})
		return
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid claims"})
		return
	}

	uid, _ := claims["sub"].(string)
	username, _ := claims["username"].(string)

	// optional: ensure user still exists
	var exists int
	_ = s.DB.QueryRow(`SELECT 1 FROM users WHERE uuid = ?`, uid).Scan(&exists)
	if exists != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "user not found"})
		return
	}

	// issue new tokens and rotate cookie
	access, new_refresh, err := generateToken(username, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token gen failed"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    new_refresh,
		Path:     "/",
		MaxAge:   int((10 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // set true on HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (s *Server) Login(c *gin.Context) {
	// Method: POST

	var login_account struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&login_account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	// look up user by username
	var storedHash string
	var uid string
	err := s.DB.QueryRow(`SELECT uuid, password_hash FROM users WHERE username = ?`, login_account.Username).Scan(&uid, &storedHash)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// NOTE: your current Register stores plain passwords in password_hash.
	// So compare directly. (Switch to bcrypt later.)
	if storedHash != login_account.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// issue tokens
	access, refresh, err := generateToken(login_account.Username, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	// set refresh cookie (HttpOnly)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		MaxAge:   int((10 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // set true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"user":         gin.H{"uuid": uid, "username": login_account.Username},
	})
}

func (s *Server) Register(c *gin.Context) {
	var register_account struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register_account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	// Insert into DB
	_, err := s.DB.Exec(`INSERT INTO users(uuid, username, password_hash) VALUES(?, ?, ?)`,
		register_account.Username, register_account.Username, register_account.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists or insert failed"})
		return
	}

	// Generate tokens using the stored "uuid" (currently same as username)
	access, refresh, err := generateToken(register_account.Username, register_account.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token gen failed"})
		return
	}

	// Set refresh cookie (HttpOnly)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		MaxAge:   int((10 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // set true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	c.JSON(http.StatusCreated, gin.H{
		"detail":       "register success",
		"access_token": access,
		"user":         gin.H{"uuid": register_account.Username, "username": register_account.Username},
	})
}

func (s *Server) getData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"detail": "hello world"})
}

// used to test functions in deveoplment
func Test() {
	access, refresh, err1 := generateToken("test", "123")

	if err1 != nil {
		fmt.Println("token generation fail:", err1)
		return
	}

	fmt.Println("ACCESS_TOKEN:", access)
	fmt.Println("REFRESH_TOKEN:", refresh)

	err2 := verifyToken(access, jwt_access_key)
	err3 := verifyToken(refresh, jwt_refresh_key)

	if err2 != nil {
		fmt.Println("access token fail:", err1)
	} else {
		fmt.Println("ACCESS IS VAILD!")
	}

	if err3 != nil {
		fmt.Println("refresh token fail:", err2)
	} else {
		fmt.Println("REFRESH IS VAILD!")
	}

	fmt.Println("code is done")

}

func main() {

	// Test()
	// return

	// some database shit not my cup of tea yet
	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := &Server{DB: db}

	router := gin.Default()

	// Method: GET
	// Purpose: testing
	router.GET("/getData", s.getData)

	// Method: POST
	// Purpose: to refresh jwt token for http and browser
	router.POST("/refresh", s.Refresh)

	// Method: POST
	// Purpose: allow users to create accounts
	// Arguments:
	//	username: string,
	//	password: string
	router.POST("/register", s.Register)

	// Method: POST
	// Purpose: users can loginto their accounts
	// Arguments:
	//	username: string,
	//	password: string
	router.POST("/login", s.Login)

	router.Run("localhost:8080")
}
