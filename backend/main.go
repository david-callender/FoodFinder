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
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// user table in the db
type Users struct {
	UUID     string
	Email    string
	Password string
}

// GLOABL VAR STORAGE
type Server struct {
	DB *sql.DB
}

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

func connectDB() (*sql.DB, error) {
	// Purpose: start intial connection to database for postgres
	// Arguments: db: *sql.DB (sql database model),
	// Return: db: *sql.DB (database model)
	//		   err: error

	db, err := sql.Open("postgres", os.Getenv("db_url"))
	if err != nil {
		fmt.Println("failed to connect to sql database", err)
		return db, err
	}

	//defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println("failed to ping sql database", err)

		return db, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (uuid UUID PRIMARY KEY, email TEXT, password TEXT)")
	if err != nil {
		return db, err
	}

	return db, nil
}

func EmailExists(db *sql.DB, email string) (bool, error) {
	// Purpose: check if email is exist in users table
	// Arguments: db: *sql.DB (sql database model),
	//			  email: string (email of user)
	// Return: exist: boolean (true is exist false if not)
	//		   err: error
	var exists bool
	err := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)",
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func AddNewUser(db *sql.DB, uuid, email, password string) error {
	// Purpose: add a new user to the users table
	// Arguments: db: *sql.DB (sql database model),
	//			  uuid: uuid (user id stuff)
	//			  email: string (email of user)
	//			  password: string (hash of password)
	// Return: err: error

	email_exist, err := EmailExists(db, email)
	if err != nil {
		return err
	}
	if email_exist {
		return fmt.Errorf("email already in use")
	}

	_, err = db.Exec(
		"INSERT INTO users (uuid, email, password) VALUES ($1, $2, $3)",
		uuid, email, password,
	)
	if err != nil {
		return err
	}
	return nil
}

func FindOneUserByEmail(db *sql.DB, email string) (*Users, error) {
	// Purpose: finds user row by their email
	// Arguments: db: *sql.DB (sql database model),
	//			  email: string (email of user)
	// Return: users: *Users (user data Struct)
	// 		   err: error

	var user Users

	err := db.QueryRow(
		"SELECT uuid, email, password FROM users WHERE email=$1",
		email,
	).Scan(&user.UUID, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		return nil, err // some other error
	}

	return &user, nil
}

func HashPassword(password string) (string, error) {
	// Purpose: hashes users password before storing in db
	// Arguments: password: string (user input password)
	// Return: password_hash: string (hash of password)
	// 		   err: error
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil // store this string in your password column
}

func CheckPasswordHash(hashedPassword, password string) error {
	// Purpose: compares users hash in db to pasword typed in
	// Arguments: password: string (user input password)
	//			  hashed_password: string (hash from db)
	// Return: result: (nil == success, nil != failed)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateToken(email string, userid string) (string, string, error) {
	// Purpose: to generate a new pair of access and refresh tokens
	// Arguments: username: string (account username),
	// 			  userid: string (account id in SQL database)
	// Return: access_token: string (access key to store in browser local storage)
	//		   refresh_token: string (this will get stored in the http cookies)

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

	email, ok := token_data["email"].(string)
	uid, ok := token_data["sub"].(string)

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload (email)"})
		return
	}
	if !ok {
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

	user_result, err := FindOneUserByEmail(s.DB, login_account.Email)
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
	uid := user_result.UUID
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

func (s *Server) Register(c *gin.Context) {
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
	uid := uuid.New().String()

	password, err := HashPassword(register_account.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": "invalid password failed to hash",
		})
		return
	}

	sql_err := AddNewUser(s.DB, uid, email, password)
	if sql_err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"detail": "Email address is already in use!",
		})
		return
	}

	// END OF SQL

	access, err := RefreshCookieTemplate(c, email, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"detail":       "register success",
		"access_token": access,
		"user":         gin.H{"uuid": register_account.Email, "username": register_account.Email},
	})
}

//---------------END-OF-API-ENDPOINTS-----------------
//----------------------------------------------------
//----------------------------------------------------

func main() {

	test_uuid := uuid.New().String()
	fmt.Println(test_uuid)

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// connect to the database
	db, err := connectDB()
	if err != nil {
		fmt.Println("database failed to initalize")
		return
	}

	s := &Server{DB: db}

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
	// Purpose: users can login to their accounts
	// Arguments:
	//	username: string,
	//	password: string
	router.POST("/login", s.Login)

	// Method: POST
	// Purpose: delete refresh token from http cookie which will require another login
	// Arguments: NONE
	router.POST("/logout", s.Logout)

	router.Run("localhost:8080")
}
