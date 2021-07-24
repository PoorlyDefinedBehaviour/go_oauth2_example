package authcontroller

import (
	"net/http"

	authservice "oauth2_example/src/domain/services/auth"

	"github.com/labstack/echo/v4"
)

func Signin(ctx echo.Context) error {
	return ctx.Redirect(http.StatusTemporaryRedirect, authservice.GenerateOAuth2URL())
}

func OAuth2Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	state := ctx.QueryParam("state")
	if state == "" {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	userInfo, err := authservice.HandleOAuth2Callback(code, state)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, userInfo)
}
