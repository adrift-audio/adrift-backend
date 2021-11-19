package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	"adrift-backend/utilities"
)

func deleteAccountController(ctx *fiber.Ctx) error {
	userID := ctx.Locals("UserId").(string)

	parsedId, parsingError := primitive.ObjectIDFromHex(userID)
	if parsingError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete User record
	UserCollection := DB.Instance.Database.Collection(DB.Collections.User)
	_, deleteError := UserCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: parsedId}},
	)
	if deleteError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete UserSecret record
	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)
	_, deleteError = UserSecretCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
	)
	if deleteError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	// delete Password record
	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)
	_, deleteError = PasswordCollection.DeleteOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userID}},
	)
	if deleteError != nil {
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
