package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"rhythmify/services/auth-service/internal/models"
)

// postgresUserRepository implements UserRepository interface
type postgresUserRepository struct {
	db *pgxpool.Pool
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

// Create creates a new user and returns the created user with ID
func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, username, password_hash, telegram_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, user.Email, user.Username, user.Password, user.TelegramID)
	
	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *postgresUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, telegram_id, created_at, updated_at
		FROM users 
		WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, telegram_id, created_at, updated_at
		FROM users 
		WHERE email = $1`

	row := r.db.QueryRow(ctx, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by their username
func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, telegram_id, created_at, updated_at
		FROM users 
		WHERE username = $1`

	row := r.db.QueryRow(ctx, query, username)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// GetByTelegramID retrieves a user by their Telegram ID
func (r *postgresUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, password_hash, telegram_id, created_at, updated_at
		FROM users 
		WHERE telegram_id = $1`

	row := r.db.QueryRow(ctx, query, telegramID)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with telegram_id %d not found", telegramID)
		}
		return nil, fmt.Errorf("failed to get user by telegram_id: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *postgresUserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET email = $2, username = $3, telegram_id = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	row := r.db.QueryRow(ctx, query, user.ID, user.Email, user.Username, user.TelegramID)
	err := row.Scan(&user.UpdatedAt)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user with id %d not found", user.ID)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// LinkTelegram links a Telegram ID to a user
func (r *postgresUserRepository) LinkTelegram(ctx context.Context, userID int64, telegramID int64) error {
	query := `
		UPDATE users 
		SET telegram_id = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, userID, telegramID)
	if err != nil {
		return fmt.Errorf("failed to link telegram: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", userID)
	}

	return nil
}

// Delete soft deletes a user (placeholder for future implementation)
func (r *postgresUserRepository) Delete(ctx context.Context, id int64) error {
	// For now, we'll just return not implemented
	// In the future, you might want to add a deleted_at column for soft deletes
	return fmt.Errorf("delete operation not implemented yet")
}

// CheckEmailExists checks if email already exists
func (r *postgresUserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// CheckUsernameExists checks if username already exists
func (r *postgresUserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}