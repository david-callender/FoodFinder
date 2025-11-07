package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	s := &Server{DB: db, LoggedIn: map[int]bool{}}

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
	router.GET("/getMenu", s.GetMenu)
	router.POST("/addFoodPreference", s.addFoodPreference)
	router.POST("/removeFoodPreference", s.removeFoodPreference)

	router.Run("localhost:8080")
}
