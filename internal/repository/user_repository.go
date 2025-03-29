package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jefdimar/go-chat-app/internal/database"
	"github.com/jefdimar/go-chat-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidLogin   = errors.New("invalid username or password")
	ErrUserNotFound   = errors.New("user not found")
	ErrInternalServer = errors.New("internal server error")
)

// CreateUser creates a new user in the database
func CreateUser(user *models.User) error {
	// Check if username already exists
	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", user.Username).Scan(&exists)
	if err != nil {
		return ErrInternalServer
	}
	if exists {
		return ErrUserExists
	}

	// Check if email already exists
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		return ErrInternalServer
	}
	if exists {
		return ErrUserExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternalServer
	}

	// Insert the new user
	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, username, email, created_at, updated_at
	`
	now := time.Now()
	err = database.DB.QueryRow(
		query,
		user.Username,
		user.Email,
		string(hashedPassword),
		now,
		now,
	).Scan(&user.ID)

	if err != nil {
		return ErrInternalServer
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	// CLear the password so it's not stored in memory
	user.Password = ""

	return nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}

	return user, nil
}

// AuthenticateUser checks if the provided credentials are valid
func AuthenticateUser(username, password string) (*models.User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, ErrInvalidLogin
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidLogin
	}

	// Clear the password so it's not stored in memory
	user.Password = ""

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServer
	}

	return user, nil
}
