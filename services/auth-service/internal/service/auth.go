package service

import (
	"context"
	"fmt"

	"rhythmify/services/auth-service/internal/models"
	"rhythmify/services/auth-service/internal/repository"
	"rhythmify/shared/jwt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, *jwt.TokenPair, error) {
	// Check if email already exists
	emailExists, err := s.userRepo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check email: %w", err)
	}
	if emailExists {
		return nil, nil, fmt.Errorf("email already exists")
	}

	// Check if username already exists
	usernameExists, err := s.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check username: %w", err)
	}
	if usernameExists {
		return nil, nil, fmt.Errorf("username already exists")
	}

	// Create user object
	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user.ToResponse(), tokens, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, *jwt.TokenPair, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	tokens, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user.ToResponse(), tokens, nil
}

// RefreshToken generates new tokens using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// Validate and refresh token
	tokens, err := s.jwtManager.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	return tokens, nil
}

// GetProfile returns user profile information
func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user.ToResponse(), nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(ctx context.Context, userID int64, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Get current user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields if provided
	if req.Email != nil {
		// Check if new email already exists (if different from current)
		if *req.Email != user.Email {
			emailExists, err := s.userRepo.CheckEmailExists(ctx, *req.Email)
			if err != nil {
				return nil, fmt.Errorf("failed to check email: %w", err)
			}
			if emailExists {
				return nil, fmt.Errorf("email already exists")
			}
		}
		user.Email = *req.Email
	}

	if req.Username != nil {
		// Check if new username already exists (if different from current)
		if *req.Username != user.Username {
			usernameExists, err := s.userRepo.CheckUsernameExists(ctx, *req.Username)
			if err != nil {
				return nil, fmt.Errorf("failed to check username: %w", err)
			}
			if usernameExists {
				return nil, fmt.Errorf("username already exists")
			}
		}
		user.Username = *req.Username
	}

	// Update user in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user.ToResponse(), nil
}

// LinkTelegram links a Telegram account to a user
func (s *AuthService) LinkTelegram(ctx context.Context, userID int64, req *models.LinkTelegramRequest) error {
	// Check if Telegram ID is already linked to another user
	existingUser, err := s.userRepo.GetByTelegramID(ctx, req.TelegramID)
	if err == nil && existingUser.ID != userID {
		return fmt.Errorf("telegram account already linked to another user")
	}

	// Link Telegram ID to user
	if err := s.userRepo.LinkTelegram(ctx, userID, req.TelegramID); err != nil {
		return fmt.Errorf("failed to link telegram: %w", err)
	}

	return nil
}

// GetUserByTelegramID retrieves a user by their Telegram ID
func (s *AuthService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user.ToResponse(), nil
}

// ValidateToken validates a JWT token and returns user claims
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
