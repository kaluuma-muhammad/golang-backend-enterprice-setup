package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appAuthorization "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/google/uuid"
)

func (h *AuthorizationHandler) AssignRoleToUser(c *gin.Context) {
	var req appAuthorization.AssignRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AssignRoleToUser(c.Request.Context(), req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Role assigned successfully",
	})
}

func (h *AuthorizationHandler) RemoveRoleFromUser(c *gin.Context) {
	var req appAuthorization.AssignRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.RemoveRoleFromUser(c.Request.Context(), req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Role removed successfully",
	})
}

func (h *AuthorizationHandler) AssignPermissionToRole(c *gin.Context) {
	var req appAuthorization.AssignPermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AssignPermissionToRole(c.Request.Context(), req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Permission assigned successfully",
	})
}

func (h *AuthorizationHandler) RemovePermissionFromRole(c *gin.Context) {
	var req appAuthorization.AssignPermissionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.RemovePermissionFromRole(c.Request.Context(), req); err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"message": "Permission removed successfully",
	})
}

func (h *AuthorizationHandler) ListUserRoles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roles, err := h.service.ListUserRoles(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, roles)
}

func (h *AuthorizationHandler) ListUserPermissions(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	permissions, err := h.service.ListUserPermissions(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, permissions)
}

func (h *AuthorizationHandler) ListRolePermissions(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	permissions, err := h.service.ListRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, permissions)
}
