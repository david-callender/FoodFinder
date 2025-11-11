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

	docclient "github.com/david-callender/FoodFinder/utils/dineocclient"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

// User table in the db
type User struct {
	ID          int
	Email       string
	Password    []byte
	DisplayName string
}

type MealWithPreference struct {
	Meal        string `json:"meal"`
	IsPreferred bool   `json:"isPreferred"`
	Id          string `json:"id"`
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
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Checks if a user exists in the database by an email
func EmailExists(db *pgxpool.Pool, email string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM "Users" WHERE "email"=$1)`,
		email,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Adds a new user to the users table
func AddNewUser(db *pgxpool.Pool, email string, password []byte, displayName string) (int, error) {
	exists, err := EmailExists(db, email)
	if err != nil {
		return -1, fmt.Errorf("failed to check for existing user: %w", err)
	}
	if exists {
		return -1, ErrEmailInUse
	}

	var id int
	err = db.QueryRow(context.Background(), `
        INSERT INTO "Users" ("email", "password", "displayName")
        VALUES ($1, $2, $3)
        RETURNING "id";
    `, email, password, displayName).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert new user: %w", err)
	}
	return id, nil
}

// Finds a user by an email
func GetByEmail(db *pgxpool.Pool, email string) (*User, error) {
	var user User

	user.Email = email

	err := db.QueryRow(context.Background(),
		`SELECT "id", "password" FROM "Users" WHERE "email"=$1`,
		email,
	).Scan(&user.ID, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
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

	sign_access, err := access_token.SignedString([]byte(os.Getenv("access_key")))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	sign_refresh, err := refresh_token.SignedString([]byte(os.Getenv("refresh_key")))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return sign_access, sign_refresh, nil
}

// Verifies a jwt token
func verifyToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
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

	return access, nil
}

// Endpoint functions here

func (s *Server) GetMenu(c *gin.Context) {
	//Method: GET

	day := c.Query("day")
	dining_hall := c.Query("diningHall")
	mealtime := c.Query("mealtime")

	// GetMenuById requires a time.Time so we have to parse the day
	day_as_time, err := time.Parse(time.DateOnly, day)
	if err != nil {
		fmt.Println("/getMenu: invalid date: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid date"})
		return
	}

	// TODO: fetch via SQL query from our database instead of directly using dineocclient
	// TODO: this is technically an API call that requires authentication. Implement this.
	menu, err := docclient.GetMenuById(dining_hall, mealtime, day_as_time)
	if err != nil {
		fmt.Println("/getMenu: failed getting menu data: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "failed getting menu data"})
	}

	// This has to be done because the Meal struct doesn't have preferences and has
	// different field names than what the frontend expects.
	meal_list := make([]MealWithPreference, 0, 20)
	for _, option := range menu.Options {
		meal_list = append(meal_list, MealWithPreference{
			Meal:        option.Name,
			IsPreferred: false,
			Id:          option.Id,
		})
	}

	c.JSON(http.StatusOK, meal_list)
}

func (s *Server) addFoodPreference(c *gin.Context) {
	var foodPreference struct {
		Meal string `json:"meal" binding:"required"`
	}

	err := c.ShouldBindJSON(&foodPreference)

	if err != nil {
		fmt.Println("/addFoodPreference: meal required: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "meal required"})
		return
	}

	fmt.Println("adding preference: ", foodPreference.Meal)
	c.Status(http.StatusOK)
}

func (s *Server) removeFoodPreference(c *gin.Context) {
	var foodPreference struct {
		Meal string `json:"meal" binding:"required"`
	}

	err := c.ShouldBindJSON(&foodPreference)

	if err != nil {
		fmt.Println("/removeFoodPreference: meal required: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "meal required"})
		return
	}

	fmt.Println("removing preference: ", foodPreference.Meal)
	c.Status(http.StatusOK)
}

//----------------------------------------------------
//----------------------------------------------------
//---------------START-OF-API-ENDPOINTS---------------

// Method: POST
func (s *Server) Refresh(c *gin.Context) {
	jwt_refresh_key := []byte(os.Getenv("refresh_key"))

	refresh_cookie, err := c.Cookie("refresh_token")

	if err != nil {
		fmt.Println("/refresh: no refresh token: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "no refresh token"})
		return
	}

	token_data, err := verifyToken(refresh_cookie, jwt_refresh_key)
	if err != nil {
		fmt.Println("/refresh: token verification failed: ", err)
		c.JSON(http.StatusForbidden, gin.H{"detail": "token verification failed"})
		return
	}

	uid_str, err := token_data.GetSubject()

	if err != nil {
		fmt.Println("/refresh: no token subject: ", err)
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload"})
		return
	}

	uid, err := strconv.Atoi(uid_str)

	if err != nil {
		fmt.Println("/refresh: invalid token subject: ", err)
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid token payload"})
		return
	}

	access, err := RefreshCookieTemplate(c, uid)

	if err != nil {
		fmt.Println("/refresh: token generation failed: ", err)
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
		fmt.Println("/login: invalid json: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "email and password required"})
		return
	}

	user_result, err := GetByEmail(s.DB, login_account.Email)
	if err != nil {
		fmt.Println("/login: database error getting user by email: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		return
	}
	if user_result == nil {
		fmt.Println("/login: user doesn't exist: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid credentials"})
		return
	}
	err = CheckPasswordHash(user_result.Password, login_account.Password)
	if err != nil {
		fmt.Println("/login: invalid password: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid credentials"})
		return

	}

	access, err := RefreshCookieTemplate(c, user_result.ID)

	if err != nil {
		fmt.Println("/login: token generation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": access,
		"displayName": user_result.DisplayName,
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
		DisplayName string `json:"displayName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register_account); err != nil {
		fmt.Println("/signup: invalid json: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "email, password, and displayName required"})
		return
	}

	email := register_account.Email

	password, err := HashPassword(register_account.Password)
	if err != nil {
		fmt.Println("/signup: failed to hash password: ", err)
		c.JSON(http.StatusForbidden, gin.H{"detail": "invalid password"})
		return
	}

	uid, err := AddNewUser(s.DB, email, password, register_account.DisplayName)
	if err != nil {
		fmt.Println("/signup: ", err)
		if errors.Is(err, ErrEmailInUse) {
			c.JSON(http.StatusConflict, gin.H{"detail": "email already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		return
	}

	access, err := RefreshCookieTemplate(c, uid)

	if err != nil {
		fmt.Println("/signup: token generation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"accessToken": access,
	})
}

//---------------END-OF-API-ENDPOINTS-----------------
//----------------------------------------------------
//----------------------------------------------------

func main() {
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatalln("Error loading .env file")
		return
	}

	// connect to the database
	db, err := connectDB()
	if err != nil {
		if db != nil {
			db.Close()
		}
		log.Fatalln("failed to initialize database pool: ", err)
		return
	}
	defer db.Close()

	s := &Server{DB: db}

	router := gin.Default()

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
	router.GET("/getMenu", s.GetMenu)
	router.POST("/addFoodPreference", s.addFoodPreference)
	router.POST("/removeFoodPreference", s.removeFoodPreference)

	router.Run("localhost:8080")
}
