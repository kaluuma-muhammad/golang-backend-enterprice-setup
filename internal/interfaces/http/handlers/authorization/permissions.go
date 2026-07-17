package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuthorization "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/google/uuid"
)

func (h *AuthorizationHandler) GetPermissions(c *gin.Context) {
	var req common.PaginationRequest

	_ = c.ShouldBindQuery(&req)

	data, meta, err := h.service.GetPermissions(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, data, meta)
}

func (h *AuthorizationHandler) CreatePermission(c *gin.Context) {
	var req appAuthorization.CreatePermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	permission, err := h.service.CreatePermission(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusCreated, permission)
}

func (h *AuthorizationHandler) GetPermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid permission id")
		return
	}

	permission, err := h.service.GetPermission(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, permission)
}

func (h *AuthorizationHandler) UpdatePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid permission id")
		return
	}

	var req appAuthorization.UpdatePermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = id

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	permission, err := h.service.UpdatePermission(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, permission)
}

func (h *AuthorizationHandler) DeletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid permission id")
		return
	}

	if err := h.service.DeletePermission(c.Request.Context(), id); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Permission deleted successfully",
	})
}
