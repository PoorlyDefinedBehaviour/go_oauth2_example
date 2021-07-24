package authcontroller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	authservice "oauth2_example/src/domain/services/auth"
)

func Signin(ctx echo.Context) error {
	return ctx.Redirect(http.StatusTemporaryRedirect, authservice.GenerateOAuth2URI())
}

func OAuth2Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	requestIdentifier := ctx.QueryParam("state")
	if requestIdentifier == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	userInfo, err := authservice.HandleOAuth2Callback(code, requestIdentifier)
	if err != nil {
		if errors.Is(err, authservice.ErrInvalidRequestIdentifier) {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"message": err.Error()})
		}
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, userInfo)
}
