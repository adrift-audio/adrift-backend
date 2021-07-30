package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func completeLogoutController(ctx *fiber.Ctx) error {
	userID := ctx.Locals("UserId").(string)

	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)
	rawUserSecretRecord := UserSecretCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
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

	now := utilities.MakeTimestamp()
	newSecret, hashError := utilities.MakeHash(
		userID + fmt.Sprintf("%v", now),
	)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	_, updateError := UserSecretCollection.UpdateOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
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
			userID,
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
