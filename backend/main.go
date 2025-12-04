package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("starting")

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
	// gin.SetMode(gin.ReleaseMode)

	defer db.Close()

	s := &Server{DB: db, LoggedIn: map[int]bool{}}

	router := gin.Default()

	corsOrigin := os.Getenv("CORS_ORIGIN");

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin}, // Next.js origin
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

	hostname := os.Getenv("HOST_ADDR")

	log.Println("started on " + hostname)

	router.Run(hostname)
}
