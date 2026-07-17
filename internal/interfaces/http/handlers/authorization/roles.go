package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuthorization "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/application/common"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/google/uuid"
)

func (h *AuthorizationHandler) GetRoles(c *gin.Context) {
	var req common.PaginationRequest

	_ = c.ShouldBindQuery(&req)

	data, meta, err := h.service.GetRoles(c.Request.Context(), req)

	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, data, meta)
}

func (h *AuthorizationHandler) CreateRole(c *gin.Context) {
	var req appAuthorization.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusCreated, role)
}

func (h *AuthorizationHandler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	var req appAuthorization.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = id

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), req)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, role)
}

func (h *AuthorizationHandler) GetRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, role)
}

func (h *AuthorizationHandler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}
