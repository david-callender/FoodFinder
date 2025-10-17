package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/cors" // cors handling later
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// user table in the db
type Users struct {
	ID         string
	Email      string
	Password   string
	DispayName string
}

// GLOABL VAR STORAGE
type Server struct {
	DB *pgxpool.Pool
}

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

var ErrEmailInUse = errors.New("email already in use")

func connectDB() (*pgxpool.Pool, error) {
	dsn := os.Getenv("connection_string")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL (or db_url) is not set")
	}

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Println("failed to connect to database", err)
		return db, err
	}

	if err := db.Ping(context.Background()); err != nil {
		fmt.Println("failed to ping database", err)
		return nil, err
	}

	return db, nil
}

// Purpose: check if email is exist in users table
// Arguments: db: *sql.DB (sql database model),
//			  email: string (email of user)
// Return: exist: boolean (true is exist false if not)
//		   err: error

func EmailExists(db *pgxpool.Pool, email string) (bool, error) {

	var exists bool
	err := db.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)",
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Purpose: add a new user to the users table
// Arguments: db: *sql.DB (sql database model),
//			  uuid: uuid (user id stuff)
//			  email: string (email of user)
//			  password: string (hash of password)
// Return: err: error

func AddNewUser(db *pgxpool.Pool, email, password string) (string, error) {
	var id string
	// Assumes: users(id UUID PRIMARY KEY DEFAULT gen_random_uuid(), email UNIQUE, ...)
	err := db.QueryRow(context.Background(), `
        INSERT INTO users (email, password, phone, displayName)
        VALUES ($1, $2, '', '')
        ON CONFLICT (email) DO NOTHING
        RETURNING id
    `, email, password).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrEmailInUse
		}
		return "", err
	}
	return id, nil
}

// Purpose: finds user row by their email
// Arguments: db: *sql.DB (sql database model),
//			  email: string (email of user)
// Return: users: *Users (user data Struct)
// 		   err: error

func ExistsByEmail(db *pgxpool.Pool, email string) (*Users, error) {

	var user Users

	err := db.QueryRow(context.Background(),
		"SELECT uuid, email, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err // some other error
	}

	return &user, nil
}

// Purpose: hashes users password before storing in db
// Arguments: password: string (user input password)
// Return: password_hash: string (hash of password)
//
//	err: error
func HashPassword(password string) (string, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil // store this string in your password column
}

// Purpose: compares users hash in db to pasword typed in
// Arguments: password: string (user input password)
//
//	hashed_password: string (hash from db)
//
// Return: result: (nil == success, nil != failed)
func CheckPasswordHash(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Purpose: to generate a new pair of access and refresh tokens
// Arguments: username: string (account username),
//
//	userid: string (account id in SQL database)
//
// Return: access_token: string (access key to store in browser local storage)
//
//	refresh_token: string (this will get stored in the http cookies)
func generateToken(email string, userid string) (string, string, error) {

	creation_time := time.Now()

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"sub":   userid,
		"iat":   creation_time.Unix(),
		"exp":   creation_time.Add(ACCESS_TOKEN_KEEPALIVE).Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"sub":   userid,
		"iat":   creation_time.Unix(),
		"exp":   creation_time.Add(REFRESH_TOKEN_KEEPALIVE).Unix(),
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

func verifyToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	// Purpose: to verfiy jwt tokens
	// Arguments: tokenString: string (token to verify),
	// 			  secretKey: string (the key of the token to verify)
	// Return: if the token is valid this function will return nil
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func RefreshCookieTemplate(c *gin.Context, email string, uid string) (string, error) {

	access, refresh, err := generateToken(email, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return "", err
	}

	exp_time := int((REFRESH_TOKEN_KEEPALIVE).Seconds())

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/",
		MaxAge:   exp_time,
		HttpOnly: true,
		Secure:   false, // set true in HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return access, err
}

//----------------------------------------------------
//----------------------------------------------------
//---------------START-OF-API-ENDPOINTS---------------

func (s *Server) Refresh(c *gin.Context) {
	// Method: POST

	jwt_refresh_key := []byte(os.Getenv("refresh_key"))

	refresh_cookie, err := c.Cookie("refresh_token")
	fmt.Printf("token_data raw: %#v\n", refresh_cookie)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "no cookie found!"})
		return
	}

	token_data, err := verifyToken(refresh_cookie, jwt_refresh_key)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"detail": "token vefication failed"})
		return
	}

	// DO SQL STUFF HERE YO

	email, ok1 := token_data["email"].(string)
	uid, ok2 := token_data["sub"].(string)

	if !ok1 {
		fmt.Println(ok1)
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload (email)"})
		return
	}
	if !ok2 {
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload (sub)"})
		return
	}

	// END OF SQL

	access, err := RefreshCookieTemplate(c, email, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (s *Server) Login(c *gin.Context) {
	// Method: POST

	var login_account struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&login_account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Email and password required"})
		return
	}

	// DO SQL STUFF HERE YO

	user_result, err := ExistsByEmail(s.DB, login_account.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		return
	}
	if user_result == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid credentials"})
		return
	}
	err = CheckPasswordHash(user_result.Password, login_account.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid credentials"})
		return

	}

	// END OF SQL

	email := user_result.Email
	uid := user_result.ID
	access, err := RefreshCookieTemplate(c, email, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"detail":       "login success",
		"user":         gin.H{"uuid": uid, "username": login_account.Email},
	})
}

func (s *Server) Logout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // tells browser to delete
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, gin.H{
		"detail": "logout success",
	})
}

func (s *Server) Signup(c *gin.Context) {
	var register_account struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register_account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Email and password required"})
		return
	}

	// DO SQL STUFF HERE YO

	email := register_account.Email

	password, err := HashPassword(register_account.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": "invalid password failed to hash",
		})
		return
	}

	uid, err := AddNewUser(s.DB, email, password)
	if err != nil {
		if errors.Is(err, ErrEmailInUse) {
			c.JSON(http.StatusConflict, gin.H{"detail": "Email address is already in use!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		fmt.Println(err)
		return
	}

	// END OF SQL

	access, err := RefreshCookieTemplate(c, email, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token": access,
	})
}

//---------------END-OF-API-ENDPOINTS-----------------
//----------------------------------------------------
//----------------------------------------------------

func main() {

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// connect to the database
	db, err := connectDB()
	if err != nil {
		fmt.Println("database failed to initialize:", err)
		return
	}
	defer db.Close()

	s := &Server{DB: db}

	router := gin.Default()

	// handle CORS requests for testing. How to avoid? Stolen from Chatgpt.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Next.js origin
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Method: POST
	// Purpose: to refresh jwt token for http and browser
	router.POST("/refresh", s.Refresh)

	// Method: POST
	// Purpose: allow users to create accounts
	// Arguments:
	//	username: string,
	//	password: string
	router.POST("/signup", s.Signup)

	// Method: POST
	// Purpose: users can login to their accounts
	// Arguments:
	//	username: string,
	//	password: string
	router.POST("/login", s.Login)

	router.Run("localhost:8080")
}
