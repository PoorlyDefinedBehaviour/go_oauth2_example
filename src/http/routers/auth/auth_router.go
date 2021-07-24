package authrouter

import (
	"github.com/labstack/echo/v4"
	authcontroller "oauth2_example/src/http/controllers/auth"
)

func RegisterRoutes(router *echo.Echo) {
	group := router.Group("/auth")

	group.GET("/signin", authcontroller.Signin)
	group.GET("/oauth2_callback", authcontroller.OAuth2Callback)
}
