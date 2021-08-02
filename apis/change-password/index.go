package changePassword

import (
	"github.com/gofiber/fiber/v2"

	"adrift-backend/middlewares"
)

func Setup(app *fiber.App) {
	group := app.Group("/api/password")

	group.Post(
		"/",
		middlewares.Authorize,
		changePasswordController,
	)
}
