package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"adrift-backend/configuration"
	"adrift-backend/utilities"
)

func AuthorizeMicroservices(ctx *fiber.Ctx) error {
	rawPassphrase := ctx.Get("Passphrase")
	if rawPassphrase == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingPassphrase,
			Status: fiber.StatusUnauthorized,
		})
	}
	trimmedPassphrase := strings.TrimSpace(rawPassphrase)
	if trimmedPassphrase == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingPassphrase,
			Status: fiber.StatusUnauthorized,
		})
	}

	if trimmedPassphrase != os.Getenv("MICROSERVICES_PASSPHRASE") {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	return ctx.Next()
}
