package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func getRecoveryCodeController(ctx *fiber.Ctx) error {
	var body GetRecoveryCodeBodyStruct
	bodyParsingError := ctx.BodyParser(&body)
	if bodyParsingError != nil {
		if fmt.Sprint(bodyParsingError) == "Unprocessable Entity" {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.MissingData,
				Status: fiber.StatusBadRequest,
			})
		}
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	email := body.Email
	trimmedEmail := strings.TrimSpace(email)

	if email == "" || trimmedEmail == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	emailIsValid := utilities.ValidateEmail(trimmedEmail)
	if !emailIsValid {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidEmail,
			Status: fiber.StatusBadRequest,
		})
	}

	UserCollection := DB.Instance.Database.Collection(DB.Collections.User)
	rawUserRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: trimmedEmail}},
	)
	userRecord := &Schemas.User{}
	rawUserRecord.Decode(userRecord)
	if userRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	code := gohelpers.RandomString(16)
	recoveryLink := os.Getenv("FRONTEND_ENDPOINT") + "/recovery/" + code

	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)
	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "recoveryCode",
					Value: code,
				},
				{
					Key:   "updated",
					Value: utilities.MakeTimestamp(),
				},
			},
		}},
	)
	if updateError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	recoveryTemplate := utilities.CreateAccountRecoveryTemplate(
		userRecord.FirstName,
		userRecord.LastName,
		recoveryLink,
	)
	utilities.SendEmail(
		userRecord.Email,
		recoveryTemplate.Subject,
		recoveryTemplate.Message,
	)

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
