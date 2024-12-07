package user

import (
	"net/http"

	"golang-tes/internal/domain"
	"golang-tes/internal/utils"
	"golang-tes/internal/utils/validator"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type updateProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param request body registerRequest true "User registration details"
// @Success 201 {object} utils.Response "User registered successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 409 {object} utils.Response "Email already exists"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", domain.ErrInvalidInput.Error())
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email", err.Error())
		return
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid password", err.Error())
		return
	}
	if err := validator.ValidateName(req.Name); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid name", err.Error())
		return
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     domain.RoleUser,
	}

	if req.Role != "" {
		if err := validator.ValidateUserRole(req.Role); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid role", err.Error())
			return
		}
		user.Role = req.Role
	}

	err := h.userUsecase.Register(c.Request.Context(), user)
	if err == domain.ErrEmailExists {
		utils.ErrorResponse(c, http.StatusConflict, "Registration failed", err.Error())
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", nil)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body loginRequest true "User login credentials"
// @Success 200 {object} utils.Response{data=map[string]string{token=string}} "Login successful"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Invalid credentials"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", domain.ErrInvalidInput.Error())
		return
	}

	if err := validator.ValidateEmail(req.Email); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid email", err.Error())
		return
	}

	token, err := h.userUsecase.Login(c.Request.Context(), req.Email, req.Password)
	if err == domain.ErrInvalidCredentials {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=domain.User} "Profile retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "user ID not found in context")
		return
	}

	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body updateProfileRequest true "User profile update details"
// @Success 200 {object} utils.Response "Profile updated successfully"
// @Failure 400 {object} utils.Response "Invalid request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "user ID not found in context")
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user := &domain.User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.userUsecase.UpdateProfile(c.Request.Context(), user)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", nil)
}
