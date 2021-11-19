package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func recoverAccountController(ctx *fiber.Ctx) error {
	var body RecoverAccounntBodyStruct
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

	code := body.Code
	password := body.Password
	if code == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	trimmedCode := strings.TrimSpace(code)
	trimmedPassword := strings.TrimSpace(password)
	if trimmedCode == "" || trimmedPassword == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)
	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "recoveryCode", Value: code}},
	)
	passwordRecord := &Schemas.Password{}
	rawPasswordRecord.Decode(passwordRecord)
	if passwordRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidRecoveryCode,
			Status: fiber.StatusUnauthorized,
		})
	}

	now := utilities.MakeTimestamp()
	newHash, hashError := utilities.MakeHash(trimmedPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	_, updateError := PasswordCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: passwordRecord.UserId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "hash",
					Value: newHash,
				},
				{
					Key:   "recoveryCode",
					Value: "",
				},
				{
					Key:   "updated",
					Value: now,
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

	newSecret, secretError := utilities.MakeHash(
		passwordRecord.UserId + fmt.Sprintf("%v", utilities.MakeTimestamp()),
	)
	if secretError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)
	_, updateError = UserSecretCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: passwordRecord.UserId}},
		bson.D{{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "secret",
					Value: newSecret,
				},
				{
					Key:   "updated",
					Value: now,
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

	// delete all of the user-specific keys from Redis via Pipeline
	redisPrefixes := gohelpers.ObjectValues(configuration.Redis.Prefixes)
	var keys []string
	for _, prefix := range redisPrefixes {
		keys = append(keys, utilities.KeyFormatter(
			prefix,
			passwordRecord.UserId,
		))
	}
	pipeline := utilities.RedisClient.Pipeline()
	for _, key := range keys {
		pipeline.Del(ctx.Context(), key)
	}
	pipeline.Exec(ctx.Context())

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
	})
}
