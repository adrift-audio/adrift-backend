package auth

import (
	"github.com/gofiber/fiber/v2"

	"adrift-backend/middlewares"
)

func Setup(app *fiber.App) {
	group := app.Group("/api/auth")

	group.Get(
		"/complete-logout",
		middlewares.Authorize,
		completeLogoutController,
	)
	group.Post("/sign-in", signInController)
	group.Post("/sign-up", signUpController)
}
