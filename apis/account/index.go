package account

import (
	"github.com/gofiber/fiber/v2"

	"adrift-backend/middlewares"
)

func Setup(app *fiber.App) {
	group := app.Group("/api/account")

	group.Get(
		"/",
		middlewares.Authorize,
		getAccountController,
	)
	group.Patch(
		"/",
		middlewares.Authorize,
		updateAccountController,
	)
}
