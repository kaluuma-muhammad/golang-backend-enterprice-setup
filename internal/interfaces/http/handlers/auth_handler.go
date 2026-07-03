package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appAuth "github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/go-api/internal/shared/validator"
)

type AuthHandler struct {
	service   appAuth.ServiceContract
	validator *validator.Validator
}

func NewAuthHandler(service appAuth.ServiceContract, validator *validator.Validator) *AuthHandler {

	return &AuthHandler{service: service, validator: validator}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req appAuth.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req appAuth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {

	var req appAuth.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := middleware.CurrentUser(c)

	err := h.service.VerifyEmail(c.Request.Context(), user.ID, req)

	if err != nil {
		response.Error(
			c,
			response.StatusCode(err),
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		gin.H{
			"message": "email verified successfully",
		},
	)
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req appAuth.ResendVerificationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.ResendVerification(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req appAuth.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	err := h.service.Logout(c.Request.Context(), middleware.SessionID(c))

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) LogoutAllSessions(c *gin.Context) {

	err := h.service.LogoutAllSessions(c.Request.Context(), middleware.UserID(c))

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}
