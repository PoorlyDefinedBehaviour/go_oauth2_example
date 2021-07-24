package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	config "oauth2_example/src/config"
	authrouter "oauth2_example/src/http/routers/auth"
	homerouter "oauth2_example/src/http/routers/home"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	if err := config.ReadConfigsFromEnv(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	authrouter.RegisterRoutes(e)
	homerouter.RegisterRoutes(e)

	if err := e.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
