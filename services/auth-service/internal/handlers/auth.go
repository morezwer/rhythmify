package handlers

import (
	"github.com/gin-gonic/gin"

	"rhythmify/services/auth-service/internal/middleware"
	"rhythmify/services/auth-service/internal/models"
	"rhythmify/services/auth-service/internal/service"
	"rhythmify/shared/response"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User registration data"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	// Register user
	user, tokens, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to register user")
		return
	}

	// Return success response
	responseData := gin.H{
		"user":   user,
		"tokens": tokens,
	}

	response.Created(c, "User registered successfully", responseData)
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "User login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	// Authenticate user
	user, tokens, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		response.InternalServerError(c, "Failed to login")
		return
	}

	// Return success response
	responseData := gin.H{
		"user":   user,
		"tokens": tokens,
	}

	response.OK(c, "Login successful", responseData)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body gin.H true "Refresh token"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	// Refresh tokens
	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	// Return success response
	response.OK(c, "Token refreshed successfully", gin.H{"tokens": tokens})
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get user profile
	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// Return success response
	response.OK(c, "Profile retrieved successfully", gin.H{"user": user})
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateUserRequest true "User update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req models.UpdateUserRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	// Update user profile
	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to update profile")
		return
	}

	// Return success response
	response.OK(c, "Profile updated successfully", gin.H{"user": user})
}

// LinkTelegram handles linking Telegram account
// @Summary Link Telegram account
// @Description Link a Telegram account to the current user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LinkTelegramRequest true "Telegram link data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/telegram [post]
func (h *AuthHandler) LinkTelegram(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req models.LinkTelegramRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request data: "+err.Error())
		return
	}

	// Link Telegram account
	err := h.authService.LinkTelegram(c.Request.Context(), userID, &req)
	if err != nil {
		if err.Error() == "telegram account already linked to another user" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, "Failed to link Telegram account")
		return
	}

	// Return success response
	response.OK(c, "Telegram account linked successfully", nil)
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check if the auth service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	response.OK(c, "Auth service is healthy", gin.H{
		"service": "auth-service",
		"status":  "healthy",
		"version": "1.0.0",
	})
}

// GetUserByTelegramID handles getting user by Telegram ID (internal endpoint)
// @Summary Get user by Telegram ID
// @Description Get user information by Telegram ID (for internal service communication)
// @Tags internal
// @Accept json
// @Produce json
// @Param telegram_id path int true "Telegram ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /internal/users/telegram/{telegram_id} [get]
func (h *AuthHandler) GetUserByTelegramID(c *gin.Context) {
	// This endpoint is for internal service communication
	// In production, you might want to add IP restrictions or API key authentication

	var req struct {
		TelegramID int64 `uri:"telegram_id" binding:"required"`
	}

	// Bind URI parameter
	if err := c.ShouldBindUri(&req); err != nil {
		response.BadRequest(c, "Invalid Telegram ID")
		return
	}

	// Get user by Telegram ID
	user, err := h.authService.GetUserByTelegramID(c.Request.Context(), req.TelegramID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// Return success response
	response.OK(c, "User found", gin.H{"user": user})
}
