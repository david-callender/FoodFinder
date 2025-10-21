package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
  
  "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/cors" // cors handling later
  
  docclient "github.com/david-callender/FoodFinder/dineocclient"
  
  "github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"	
	
)

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

// user table in the db
type User struct {
	ID          int
	Email       string
	Password    []byte
	DisplayName string
}

type mealWithPreference struct {
	Meal         string `json:"meal"`
	Is_preferred bool   `json:"is_preferred"`
	Id           string `json:"id"`
}

// GLOBAL VAR STORAGE
type Server struct {
	DB *pgxpool.Pool
}

var ErrEmailInUse = errors.New("email already in use")

// Connects to the db and returns a connection pool.
func connectDB() (*pgxpool.Pool, error) {
	conStr := os.Getenv("DATABASE_URL")
	if conStr == "" {
		return nil, fmt.Errorf("environment variable 'DATABASE_URL' is not set")
	}

	db, err := pgxpool.New(context.Background(), conStr)
	if err != nil {
		fmt.Println("failed to connect to database", err)
		return db, err
	}

	if err := db.Ping(context.Background()); err != nil {
		fmt.Println("failed to ping database", err)
		return nil, err
	}


// Checks if a user exists in the database by an email
func EmailExists(db *pgxpool.Pool, uid int) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)",
		uid,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Adds a new user to the users table
func AddNewUser(db *pgxpool.Pool, email string, password []byte, displayName string) (int, error) {
	var id int
	err := db.QueryRow(context.Background(), `
        INSERT INTO users (email, password, displayName)
        VALUES ($1, $2, $3)
        ON CONFLICT (email) DO NOTHING
        RETURNING id
    `, email, password, displayName).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, ErrEmailInUse
		}
		return -1, err
	}
	return id, nil
}

// Finds a user by an email
func GetByEmail(db *pgxpool.Pool, email string) (*User, error) {
	var user User

	user.Email = email

	err := db.QueryRow(context.Background(),
		"SELECT id, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err // some other error
	}

	return &user, nil
}

// Hashes a password using bcrypt
func HashPassword(password string) ([]byte, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashed, nil
}

// Compares users hash in db to typed password
// Returns nil on success, err on fail
func CheckPasswordHash(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

// Generates a new pair of access and refresh tokens. Returns (access_token, refresh_token, error)
func generateToken(userid int) (string, string, error) {

	creation_time := time.Now()

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userid,
		"iat": creation_time.Unix(),
		"exp": creation_time.Add(ACCESS_TOKEN_KEEPALIVE).Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userid,
		"iat": creation_time.Unix(),
		"exp": creation_time.Add(REFRESH_TOKEN_KEEPALIVE).Unix(),
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

// Verifies a jwt token
func verifyToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
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
  
  // adds the refresh token to the http cookies and returns the access token
func RefreshCookieTemplate(c *gin.Context, uid int) (string, error) {
	access, refresh, err := generateToken(uid)
  
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

// Endpoint functions here

func (s *Server) GetMenu(c *gin.Context) {
	//Method: GET

	day := c.Query("day")
	dining_hall := c.Query("dining_hall")
	mealtime := c.Query("mealtime")

	// GetMenuById requires a time.Time so we have to parse the day
	day_as_time, err := time.Parse(time.DateOnly, day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid date"})
		return
	}

	// TODO: fetch via SQL query from our database instead of directly using dineocclient
	// TODO: this is technically an API call that requries authentication. Implement this.
	menu, err := docclient.GetMenuById(dining_hall, mealtime, day_as_time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "failed getting menu data"})
	}

	// This has to be done because the Meal struct doesn't have preferences and has
	// different field names than what the frontend expects.
	meal_list := make([]mealWithPreference, 0, 20)
	for _, option := range menu.Options {
		meal_list = append(meal_list, mealWithPreference{
			Meal:         option.Name,
			Is_preferred: false,
			Id:           option.Id,
		})
	}

	c.JSON(http.StatusOK, meal_list)
}


//----------------------------------------------------
//----------------------------------------------------
//---------------START-OF-API-ENDPOINTS---------------

// Method: POST
func (s *Server) Refresh(c *gin.Context) {
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

	uid_str, err := token_data.GetSubject()

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload"})
		return
	}

	uid, err := strconv.Atoi(uid_str)

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload"})
		return
	}

	access, err := RefreshCookieTemplate(c, uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

// Method: POST
func (s *Server) Login(c *gin.Context) {
	var login_account struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&login_account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Email and password required"})
		return
	}

	user_result, err := GetByEmail(s.DB, login_account.Email)
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

	access, err := RefreshCookieTemplate(c, user_result.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		"display_name": user_result.DisplayName,
	})
}

// Method: POST
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

	c.Status(http.StatusOK)
}

func (s *Server) Signup(c *gin.Context) {
	var register_account struct {
		Email       string `json:"email" binding:"required"`
		Password    string `json:"password" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register_account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Email, password, and display name required"})
		return
	}

	email := register_account.Email

	password, err := HashPassword(register_account.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": "invalid password failed to hash",
		})
		return
	}

	uid, err := AddNewUser(s.DB, email, password, register_account.DisplayName)
	if err != nil {
		if errors.Is(err, ErrEmailInUse) {
			c.JSON(http.StatusConflict, gin.H{"detail": "Email address is already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		return
	}

	access, err := RefreshCookieTemplate(c, uid)

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

	router.POST("/refresh", s.Refresh)
	router.POST("/signup", s.Signup)
	router.POST("/login", s.Login)

	// Method: GET
	// Purpose: Fetch a personalized menu with preference data
	// Arguments:
	//	location: string (dineoncampus location ID)
	//	mealtime: string ("breakfast", "lunch", "dinner", or "everyday")
	//	day: string (YYYY-MM-DD)
	router.GET("/getmenu", s.GetMenu)

	router.Run("localhost:8080")
}
