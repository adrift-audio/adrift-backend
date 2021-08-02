package changePassword

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func changePasswordController(ctx *fiber.Ctx) error {
	var body ChangePasswordBodyStruct
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

	newPassword := body.NewPassword
	oldPassword := body.OldPassword

	if newPassword == "" || oldPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	trimmedNewPassword := strings.TrimSpace(newPassword)
	trimmedOldPassword := strings.TrimSpace(oldPassword)

	if trimmedOldPassword == "" || trimmedNewPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	userID := ctx.Locals("UserId").(string)
	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)

	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
	)
	passwordRecord := &Schemas.Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	passwordIsValid, comparisonError := utilities.CompareHashes(
		oldPassword,
		passwordRecord.Hash,
	)
	if !passwordIsValid || comparisonError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	hash, hashError := utilities.MakeHash(newPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "hash",
					Value: hash,
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

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
