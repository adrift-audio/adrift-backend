package changePassword

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	group := app.Group("/api/password")

	group.Post("/", changePasswordController)
}
