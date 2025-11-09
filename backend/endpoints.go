package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Server) GetMenu(c *gin.Context, db *pgxpool.Pool) {
	//Method: GET

	accessToken := c.Query("accessToken")
	day := c.Query("day")
	dining_hall := c.Query("diningHall")
	mealtime := c.Query("mealtime")

	uid, err := s.protectRoute(accessToken)

	if err != nil {
		fmt.Println("/getMenu: not authenticated: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "unauthenticated"})
		return
	}

	// GetCacheMenu requires a time.Time so we have to parse the day
	day_as_time, err := time.Parse(time.DateOnly, day)
	if err != nil {
		fmt.Println("/getMenu: invalid date: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid date"})
		return
	}

	menu, err := GetCacheMenu(db, dining_hall, mealtime, day_as_time)
	if err != nil {
		fmt.Println("/getMenu: failed getting menu data: ", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"detail": "failed getting menu data"},
		)
		return
	}

	userPrefs, err := GetUserPrefs(db, uid)
	if err != nil {
		fmt.Println("/getMenu: failed getting user preferences: ", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"detail": "failed getting user preferences"},
		)
		return
	}

	for _, option := range menu {
		option.IsPreferred = userPrefs[option.Meal]
	}

	c.JSON(http.StatusOK, menu)
}

func (s *Server) addFoodPreference(c *gin.Context) {
	var foodPreference struct {
		AccessToken string `json:"accessToken" binding:"required"`
		Meal        string `json:"meal" binding:"required"`
	}

	err := c.ShouldBindJSON(&foodPreference)

	if err != nil {
		fmt.Println("/addFoodPreference: accessToken and meal required: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "accessToken and meal required"})
		return
	}

	id, err := s.protectRoute(foodPreference.AccessToken)

	if err != nil {
		fmt.Println("/addFoodPreference: not authenticated: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "unauthenticated"})
		return
	}

	_, err = s.DB.Query(context.Background(), `INSERT INTO "Preferences" ("user", "preference") VALUES ($1, $2)`, id, foodPreference.Meal)

	if err != nil {
		fmt.Println("/addFoodPreference: failed in insert: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) removeFoodPreference(c *gin.Context) {
	var foodPreference struct {
		AccessToken string `json:"accessToken" binding:"required"`
		Meal        string `json:"meal" binding:"required"`
	}

	err := c.ShouldBindJSON(&foodPreference)

	if err != nil {
		fmt.Println("/removeFoodPreference: accessToken and meal required: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "accessToken and meal required"})
		return
	}

	id, err := s.protectRoute(foodPreference.AccessToken)

	if err != nil {
		fmt.Println("/addFoodPreference: not authenticated: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "unauthenticated"})
		return
	}

	_, err = s.DB.Query(context.Background(), `DELETE FROM "Preferences" WHERE "user" = $1 AND "preference" = $2`, id, foodPreference.Meal)

	if err != nil {
		fmt.Println("/addFoodPreference: failed in insert: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "database error"})
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
		c.JSON(http.StatusBadRequest, gin.H{"detail": "no refresh token"})
		return
	}

	token_data, err := verifyToken(refresh_cookie, jwt_refresh_key)
	if err != nil {
		fmt.Println("/refresh: token verification failed: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "token verification failed"})
		return
	}

	uid_str, err := token_data.GetSubject()

	if err != nil {
		fmt.Println(token_data)
		fmt.Println("/refresh: no token subject: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid token payload"})
		return
	}

	uid, err := strconv.Atoi(uid_str)

	if err != nil {
		fmt.Println("/refresh: invalid token subject: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid token payload"})
		return
	}

	access, err := RefreshCookieTemplate(c, uid)

	if err != nil {
		fmt.Println("/refresh: token generation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": access})
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

	_, err = RefreshCookieTemplate(c, user_result.ID)

	if err != nil {
		fmt.Println("/login: token generation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	s.LoggedIn[user_result.ID] = true

	c.JSON(http.StatusOK, gin.H{
		"displayName": user_result.DisplayName,
	})
}

// Method: POST
func (s *Server) Logout(c *gin.Context) {
	jwt_refresh_key := []byte(os.Getenv("refresh_key"))
	refresh_cookie, err := c.Cookie("refresh_token")

	if err != nil {
		fmt.Println("/logout: no refresh token: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid jwt"})
		return
	}

	token_data, err := verifyToken(refresh_cookie, jwt_refresh_key)

	if err != nil {
		fmt.Println("/logout: invalid refresh token: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid jwt"})
		return
	}

	subjectStr, err := token_data.GetSubject()

	if err != nil {
		fmt.Println("/logout: no subject in refresh token: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid jwt"})
		return
	}

	subject, err := strconv.Atoi(subjectStr)

	if err != nil {
		fmt.Println("/logout: invalid subject: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "invalid jwt"})
		return
	}

	s.LoggedIn[subject] = false

	newCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // tells browser to delete
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, newCookie)

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
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid password"})
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

	_, err = RefreshCookieTemplate(c, uid)

	if err != nil {
		fmt.Println("/signup: token generation failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "token generation failed"})
		return
	}

	s.LoggedIn[uid] = true

	c.Status(http.StatusCreated)
}
