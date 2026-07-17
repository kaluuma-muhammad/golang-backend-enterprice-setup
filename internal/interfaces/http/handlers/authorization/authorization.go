package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-api/internal/interfaces/http/response"
	"github.com/google/uuid"
)

func (h *AuthorizationHandler) UserHasRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	ok, err := h.service.UserHasRole(c.Request.Context(), userID, roleID)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"hasRole": ok,
	})
}

func (h *AuthorizationHandler) UserHasPermission(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	permission := c.Param("permission")

	ok, err := h.service.UserHasPermission(c.Request.Context(), userID, permission)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"hasPermission": ok,
	})
}

func (h *AuthorizationHandler) RoleHasPermission(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("roleId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid role id")
		return
	}

	permissionID, err := uuid.Parse(c.Param("permissionId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid permission id")
		return
	}

	ok, err := h.service.RoleHasPermission(c.Request.Context(), roleID, permissionID)
	if err != nil {
		response.Error(c, response.StatusCode(err), err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"hasPermission": ok,
	})
}
