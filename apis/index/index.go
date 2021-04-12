package index

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	group := app.Group("/")

	group.Get("/", indexController)
	group.Get("/api", indexController)
}
