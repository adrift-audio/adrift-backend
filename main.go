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

	"adrift-backend/apis/auth"
	"adrift-backend/apis/index"
	"adrift-backend/configuration"
	"adrift-backend/database"
	"adrift-backend/utilities"
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

	databaseError := database.ConnectMongo()
	if databaseError != nil {
		log.Fatal(databaseError)
		return
	}

	redisError := utilities.ConnectRedis()
	if redisError != nil {
		log.Fatal(redisError)
		return
	}

	app := fiber.New()

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

	auth.Setup(app)
	index.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5611"
	}

	launchError := app.Listen(":" + port)
	if launchError != nil {
		panic(launchError)
	}
}
