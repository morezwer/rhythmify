package repository

import (
	"context"

	"rhythmify/services/auth-service/internal/models"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// Create creates a new user and returns the created user with ID
	Create(ctx context.Context, user *models.User) error
	
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id int64) (*models.User, error)
	
	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	
	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	
	// GetByTelegramID retrieves a user by their Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	
	// Update updates user information
	Update(ctx context.Context, user *models.User) error
	
	// LinkTelegram links a Telegram ID to a user
	LinkTelegram(ctx context.Context, userID int64, telegramID int64) error
	
	// Delete soft deletes a user (if needed in future)
	Delete(ctx context.Context, id int64) error
	
	// CheckEmailExists checks if email already exists
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	
	// CheckUsernameExists checks if username already exists
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
}