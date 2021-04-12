package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"adrift-backend/apis/index"
	"adrift-backend/configuration"
	"adrift-backend/database"
)

func main() {
	env := os.Getenv("ENV")
	if env != configuration.Environments.Heroku {
		envError := godotenv.Load()
		if envError != nil {
			log.Fatal(envError)
			return
		}
	}

	databaseError := database.Connect()
	if databaseError != nil {
		log.Fatal(databaseError)
		return
	}

	app := fiber.New()

	// middlewares
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))
	app.Use(limiter.New(limiter.Config{
		Expiration: 30 * time.Second,
		Max:        15,
	}))
	app.Use(logger.New())

	index.Setup(app)

	// get the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "5611"
	}

	// launch the app
	launchError := app.Listen(":" + port)
	if launchError != nil {
		panic(launchError)
	}
}
