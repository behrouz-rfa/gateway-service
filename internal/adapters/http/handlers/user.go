package handlers

import (
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/behrouz-rfa/gateway-service/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for film requests
type UserHandler struct {
	svc ports.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

// registerRequest represents the request body for creating a user
type registerRequest struct {
	Name     string `json:"name" binding:"omitempty,required" example:"Jon"`
	Email    string `json:"email" binding:"omitempty,required" example:"doa@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/users/register [post]
func (h UserHandler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.UserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	usr, err := h.svc.CreateUser(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, usr)
}

// registerRequest represents the request body for creating a user
type loginRequest struct {
	Email    string `json:"email" binding:"omitempty,required" example:"jondoa"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Login
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body		loginRequest	true	"Register request"
//	@Success		200				{object}	userResponse	"User created"
//	@Failure		400				{object}	errorResponse	"Validation error"
//	@Failure		401				{object}	errorResponse	"Unauthorized error"
//	@Failure		404				{object}	errorResponse	"Data not found error"
//	@Failure		409				{object}	errorResponse	"Data conflict error"
//	@Failure		500				{object}	errorResponse	"Internal server error"
//	@Router			/users/login [post]
func (h UserHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.UserLoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	usr, err := h.svc.Login(ctx, &user)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, usr)
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Page  uint64 `form:"page" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"10"`
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"User ID"
//	@Success		200	{object}	userResponse	"User displayed"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [get]
//	@Security		BearerAuth
func (h UserHandler) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, err := h.svc.GetUser(ctx, userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, user)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"User ID"
//	@Success		200	{object}	response		"User deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/users/{id} [delete]
//	@Security		BearerAuth
func (h UserHandler) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	err := h.svc.DeleteUser(ctx, userID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
