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
	group.Get(
		"/secret/:id",
		middlewares.AuthorizeMicroservices,
		getSecretController,
	)
	group.Post("/get-code", getRecoveryCodeController)
	group.Post("/recover-account", recoverAccountController)
	group.Post("/sign-in", signInController)
	group.Post("/sign-up", signUpController)
}
