package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/application/security"
	appUser "github.com/go-api/internal/application/user"
	"github.com/go-api/internal/interfaces/http/middleware"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/go-api/internal/shared/validator"
)

type UserHandler struct {
	service         appUser.UserServiceContract
	securityService *security.Service
	validator       *validator.Validator
}

func NewUserHandler(service appUser.UserServiceContract, securityService *security.Service, validator *validator.Validator) *UserHandler {

	return &UserHandler{service: service, securityService: securityService, validator: validator}
}

func (h *UserHandler) GetAuthUser(c *gin.Context) {
	user, err := h.service.GetAuthUser(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserAccount(c *gin.Context) {
	var req appUser.UpdateUserAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.UpdateUserAccount(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req appUser.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.UpdatePassword(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	response.Success(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserAvatar(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}
	defer src.Close()

	req := appUser.UploadAvatarRequest{
		File:        src,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
	}
	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.UpdateUserAvatar(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, res)
}

func (h *UserHandler) GetAuditLogs(c *gin.Context) {
	var req common.PaginationRequest

	_ = c.ShouldBindQuery(&req)

	data, meta, err := h.securityService.GetAuditLogs(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, data, meta)
}

func (h *UserHandler) GetLoginHistory(c *gin.Context) {
	var req common.PaginationRequest

	_ = c.ShouldBindQuery(&req)

	data, meta, err := h.securityService.GetLoginHistory(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, data, meta)
}

func (h *UserHandler) GetUserSessions(c *gin.Context) {
	var req common.PaginationRequest

	_ = c.ShouldBindQuery(&req)

	data, meta, err := h.securityService.GetUserSessions(c.Request.Context(), middleware.UserID(c), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, data, meta)
}

func (h *UserHandler) GetCurrentSessions(c *gin.Context) {

	data, err := h.securityService.GetCurrentSessions(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, data)
}
