package middlewares

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func Authorize(ctx *fiber.Ctx) error {
	// check if token was provided
	rawToken := ctx.Get("Authorization")
	if rawToken == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}
	trimmedToken := strings.TrimSpace(rawToken)
	if trimmedToken == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	// get token payload
	bytePayload, decodeError := utilities.DecodePayload(trimmedToken)
	if decodeError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	// parse payload
	var parsedPayload PayloadContent
	parsingError := json.Unmarshal(bytePayload, &parsedPayload)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidToken,
			Status: fiber.StatusUnauthorized,
		})
	}

	updateExpire := true

	// check user secret in Redis
	key := utilities.KeyFormatter(
		configuration.Redis.Prefixes.Secret+"asd",
		parsedPayload.UserID,
	)
	redisContext := context.Background()
	secret, redisError := utilities.RedisClient.Get(redisContext, key).Result()
	if redisError != nil {
		// if error is not about the missing data
		if redisError != utilities.RedisNil {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.InternalServerError,
				Status: fiber.StatusInternalServerError,
			})
		}

		// user secret was not found in Redis
		UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)
		rawUserSecretRecord := UserSecretCollection.FindOne(
			ctx.Context(),
			bson.D{{Key: "userId", Value: parsedPayload.UserID}},
		)
		userSecretRecord := &Schemas.UserSecret{}
		rawUserSecretRecord.Decode(userSecretRecord)
		if userSecretRecord.ID == "" {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.AccessDenied,
				Status: fiber.StatusUnauthorized,
			})
		}

		// store secret in Redis
		redisSetError := utilities.RedisClient.Set(
			redisContext,
			key,
			userSecretRecord.Secret,
			configuration.Redis.TTL,
		).Err()
		if redisSetError != nil {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.InternalServerError,
				Status: fiber.StatusInternalServerError,
			})
		}

		secret = userSecretRecord.Secret
		updateExpire = false
	}

	// validate token
	validationError := utilities.ValidateToken(trimmedToken, secret)
	if validationError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	// update EXPIRE for the record in Redis if necessary
	if updateExpire {
		expireError := utilities.RedisClient.Expire(
			redisContext,
			key,
			configuration.Redis.TTL,
		).Err()
		if expireError != nil {
			return utilities.Response(utilities.ResponseParams{
				Ctx:    ctx,
				Info:   configuration.ResponseMessages.InternalServerError,
				Status: fiber.StatusInternalServerError,
			})
		}
	}

	// store client and token data in Locals
	ctx.Locals("Client", parsedPayload.Client)
	ctx.Locals("UserId", parsedPayload.UserID)
	return ctx.Next()
}
