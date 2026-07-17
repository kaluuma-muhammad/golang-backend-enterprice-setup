package authorization

import (
	app "github.com/go-api/internal/application/authorization"
	"github.com/go-api/internal/shared/validator"
)

type AuthorizationHandler struct {
	service   app.AuthorizationServiceContract
	validator *validator.Validator
}

func NewAuthorizationHandler(service app.AuthorizationServiceContract, validator *validator.Validator) *AuthorizationHandler {
	return &AuthorizationHandler{
		service:   service,
		validator: validator,
	}
}
