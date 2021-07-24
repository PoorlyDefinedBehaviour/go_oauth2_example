package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	config "oauth2_example/src/config"
	authrouter "oauth2_example/src/http/routers/auth"
	homerouter "oauth2_example/src/http/routers/home"
	inmemorycache "oauth2_example/src/infra/in_memory_cache"
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

	cache := inmemorycache.New()

	cache.SetWithExpiration("key", "value", 1*time.Minute)

	if err := e.Start(":3000"); err != nil {
		log.Fatal(err)
	}
}
