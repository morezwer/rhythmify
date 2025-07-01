package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID         int64     `json:"id" db:"id"`
	Email      string    `json:"email" db:"email"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password_hash"`
	TelegramID *int64    `json:"telegram_id,omitempty" db:"telegram_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserRequest represents request to create a new user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents request to update user profile
type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}

// LinkTelegramRequest represents request to link Telegram account
type LinkTelegramRequest struct {
	TelegramID int64 `json:"telegram_id" binding:"required"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	TelegramID *int64    `json:"telegram_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the user's hashed password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse converts User to UserResponse (removes sensitive data)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		TelegramID: u.TelegramID,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}