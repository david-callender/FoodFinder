package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// GLOABL VAR STORAGE
type Server struct {
	DB *sql.DB
}

// open/create SQLite and ensure schema
func connectDB() (string, error) {
	db := "it worked"
	return db, nil
}

func FindOneUserByID(db *sql.DB, id string) (string, error) {

	return "user data", nil
}

func UpdateOneUserById(db *sql.DB, id string) (string, error) {
	return "update user succesful", nil
}

func generateToken(username string, userid string) (string, string, error) {
	// Purpose: to generate a new pair of access and refresh tokens
	// Arguments: username: string (account username),
	// 			  userid: string (account id in SQL database)
	// Return: access_token: string (access key to store in browser local storage)
	//		   refresh_token: string (this will get stored in the http cookies)

	right_now := time.Now()
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      right_now.Unix(),
		"exp":      right_now.Add(time.Minute * 7).Unix(), // expires in 7 minutes
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      right_now.Unix(),
		"exp":      right_now.Add(time.Hour * 24 * 10).Unix(), // expires in 10 days
	})

	sign_access, err1 := access_token.SignedString([]byte(os.Getenv("access_key")))
	sign_refresh, err2 := refresh_token.SignedString([]byte(os.Getenv("refresh_key")))

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

func RefreshCookieTemplate(c *gin.Context, username string, uid string) (string, error) {

	access, refresh, err := generateToken(username, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return "", err
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		MaxAge:   int((10 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false, // set true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return access, err
}

func (s *Server) Refresh(c *gin.Context) {
	// Method: POST

	jwt_refresh_key := []byte(os.Getenv("refresh_key"))

	refresh_cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "no cookie found!"})
		return
	}
	fmt.Println(refresh_cookie)

	err = verifyToken(refresh_cookie, jwt_refresh_key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token vefication failed"})
		return
	}

	// DO SQL STUFF HERE YO

	// END OF SQL

	username := "test"
	uid := "1"
	access, err := RefreshCookieTemplate(c, username, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (s *Server) Login(c *gin.Context) {
	// Method: POST

	var login_account struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&login_account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "username and password required"})
		return
	}

	// DO SQL STUFF HERE YO

	// END OF SQL

	username := "test"
	uid := "1"
	access, err := RefreshCookieTemplate(c, username, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"detail":       "login success",
		"user":         gin.H{"uuid": uid, "username": login_account.Username},
	})
}

func (s *Server) Register(c *gin.Context) {
	var register_account struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register_account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "username and password required"})
		return
	}

	// DO SQL STUFF HERE YO

	// END OF SQL

	username := "test"
	uid := "1"
	access, err := RefreshCookieTemplate(c, username, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"detail":       "register success",
		"access_token": access,
		"user":         gin.H{"uuid": register_account.Username, "username": register_account.Username},
	})
}

func main() {

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// connect to the database
	_, err := connectDB()
	if err != nil {
		fmt.Println("")
		return
	}

	s := &Server{DB: nil}

	router := gin.Default()

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
