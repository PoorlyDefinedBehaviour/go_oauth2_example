package main

import (
	"github.com/labstack/echo/v4"
	authrouter "oauth2_example/src/http/routers/auth"
)

func main() {
	e := echo.New()

	authrouter.RegisterRoutes(e)
}
