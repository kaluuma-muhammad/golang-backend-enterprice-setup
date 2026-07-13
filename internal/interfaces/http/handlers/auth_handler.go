package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appAuth "github.com/go-api/internal/application/auth"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/go-api/internal/shared/validator"
)

type AuthHandler struct {
	service   appAuth.AuthServiceContract
	validator *validator.Validator
}

func NewAuthHandler(service appAuth.AuthServiceContract, validator *validator.Validator) *AuthHandler {

	return &AuthHandler{service: service, validator: validator}
}

func getClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	return c.ClientIP()
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

	ip := getClientIP(c)
	result, err := h.service.Register(c.Request.Context(), req, c.Request.UserAgent(), ip)

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

	ip := getClientIP(c)
	result, err := h.service.Login(c.Request.Context(), req, c.Request.UserAgent(), ip)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req appAuth.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.ForgotPassword(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Forgot password email sent successfully"})
}

func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	var req appAuth.ActivateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := middleware.CurrentUser(c)

	err := h.service.ActivateAccount(c.Request.Context(), user.ID, req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Account activated successfully"})
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

func (h *AuthHandler) VerifyResetCode(c *gin.Context) {
	var req appAuth.VerifyResetCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.VerifyResetCode(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req appAuth.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Password reset successfully"})
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
	ip := getClientIP(c)
	err := h.service.Logout(c.Request.Context(), middleware.UserID(c), middleware.SessionID(c), ip, c.Request.UserAgent())

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) LogoutAllSessions(c *gin.Context) {

	ip := getClientIP(c)
	err := h.service.LogoutAllSessions(c.Request.Context(), middleware.UserID(c), ip, c.Request.UserAgent())

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}
