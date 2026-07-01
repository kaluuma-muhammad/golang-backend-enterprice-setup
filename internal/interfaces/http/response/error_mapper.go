package response

import (
	"net/http"

	"github.com/go-api/internal/application/auth"
	appAuth "github.com/go-api/internal/application/auth"
)

func StatusCode(err error) int {

	switch err {

	case auth.ErrUnauthorized, auth.ErrMissingAuthorizationHeader, auth.ErrInvalidAuthorizationHeader, auth.ErrInvalidRefreshToken:
		return http.StatusUnauthorized

	case appAuth.ErrEmailAlreadyExists:
		return http.StatusConflict

	case appAuth.ErrInvalidCredentials:
		return http.StatusUnauthorized

	case auth.ErrEmailNotVerified:
		return http.StatusForbidden

	default:
		return http.StatusInternalServerError
	}
}
