package homerouter

import (
	"github.com/labstack/echo/v4"
	homecontroller "oauth2_example/src/http/controllers/home"
)

func RegisterRoutes(router *echo.Echo) {
	router.GET("", homecontroller.HomePage)
}
