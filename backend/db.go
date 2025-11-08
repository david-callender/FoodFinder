package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

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

type Server struct {
	DB       *pgxpool.Pool
	LoggedIn map[int]bool
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

// Checks if a user exists in the database by an email.
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

// Adds a new user to the users table.
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
		RETURNING "id"`, email, password, displayName).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to insert new user: %w", err)
	}
	return id, nil
}

// Finds a user by an email.
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

// Produces a map of preferences to true for a given user.
func GetUserPrefs(db *pgxpool.Pool, uid int) (map[string]bool, error) {
	prefs := make(map[string]bool)

	prefRows, err := db.Query(
		context.Background(),
		`SELECT "preference" FROM "Preferences" WHERE "user"=$1`,
		uid,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return prefs, nil
	} else if err != nil {
		return nil, fmt.Errorf("GetUserPrefs: failed database query: %v", err)
	}

	for prefRows.Next() {
		var pref string
		err = prefRows.Scan(&pref)
		if err != nil {
			return nil, fmt.Errorf("GetUserPrefs: failed reading row: %v", err)
		}

		prefs[pref] = true
	}

	return prefs, err
}

// Hashes a password using bcrypt.
func HashPassword(password string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hashed, nil
}

// Compares users hash in db to typed password. Returns nil on success, or an
// error on fail.
func CheckPasswordHash(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
