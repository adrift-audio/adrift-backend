package auth

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func getSecretController(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)
	rawUserSecretRecord := UserSecretCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: id}},
	)
	userSecretRecord := &Schemas.UserSecret{}
	rawUserSecretRecord.Decode(userSecretRecord)
	if userSecretRecord.ID == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidUserID,
			Status: fiber.StatusUnauthorized,
		})
	}

	redisError := utilities.RedisClient.Set(
		ctx.Context(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.Secret,
			userSecretRecord.UserId,
		),
		userSecretRecord.Secret,
		configuration.Redis.TTL,
	).Err()
	if redisError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"secret": userSecretRecord.Secret,
		},
	})
}
