package auth

import (
	"github.com/gofiber/fiber/v2"

	"adrift-backend/utilities"
)

func getRecoveryCodeController(ctx *fiber.Ctx) error {
	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
