package auth

import (
	"github.com/gofiber/fiber/v2"

	"adrift-backend/configuration"
	"adrift-backend/utilities"
)

func signInController(ctx *fiber.Ctx) error {
	var body SignInBodyStruct
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
