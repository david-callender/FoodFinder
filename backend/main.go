package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	docclient "github.com/david-callender/FoodFinder/dineocclient"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"	// cors handling later
)

// GLOBAL VAR STORAGE
type Server struct {
	DB *sql.DB
}

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

type mealWithPreference struct {
	Meal         string `json:"meal"`
	Is_preferred bool   `json:"is_preferred"`
	Id           string `json:"id"`
}

// INTERNAL USE FUNCTIONS
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

	creation_time := time.Now()

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      creation_time.Unix(),
		"exp":      creation_time.Add(ACCESS_TOKEN_KEEPALIVE).Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sub":      userid,
		"iat":      creation_time.Unix(),
		"exp":      creation_time.Add(REFRESH_TOKEN_KEEPALIVE).Unix(),
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

func RefreshCookieTemplate(c *gin.Context, username string, uid string) (string, error) {

	access, refresh, err := generateToken(username, uid)
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

func (s *Server) Refresh(c *gin.Context) {
	// Method: POST

	jwt_refresh_key := []byte(os.Getenv("refresh_key"))

	refresh_cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "no cookie found!"})
		return
	}

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
		fmt.Println("database failed to initalize")
		return
	}

	s := &Server{DB: nil}

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
	router.POST("/register", s.Register)

	// Method: POST
	// Purpose: users can login to their accounts
	// Arguments:
	//	username: string,
	//	password: string
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
