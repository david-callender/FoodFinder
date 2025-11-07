package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ACCESS_TOKEN_KEEPALIVE = time.Minute * 7
const REFRESH_TOKEN_KEEPALIVE = time.Hour * 24 * 10

// Generates a new pair of access and refresh tokens. Returns (access_token,
// refresh_token).
func generateToken(userid int) (string, string, error) {
	creation_time := time.Now()

	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userid),
		"iat": creation_time.Unix(),
		"exp": creation_time.Add(ACCESS_TOKEN_KEEPALIVE).Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userid),
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

func (s *Server) protectRoute(accessToken string) (int, error) {
	accessKey := os.Getenv("access_key")

	token, err := verifyToken(accessToken, []byte(accessKey))

	if err != nil {
		return 0, fmt.Errorf("failed to verify token: %w", err)
	}

	subjectStr, ok := token["sub"]

	if !ok {
		return 0, fmt.Errorf("no subject in token: %w", err)
	}

	subject, err := strconv.Atoi(subjectStr.(string))

	if err != nil {
		return 0, fmt.Errorf("invalid subject: %w", err)
	}

	loggedIn, ok := s.LoggedIn[subject]

	if !(ok && loggedIn) {
		return 0, fmt.Errorf("not logged in")
	}

	return subject, nil
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