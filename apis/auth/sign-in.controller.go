package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/bson"

	"adrift-backend/configuration"
	DB "adrift-backend/database"
	Schemas "adrift-backend/database/schemas"
	"adrift-backend/utilities"
)

func signInController(ctx *fiber.Ctx) error {
	var body SignInBodyStruct
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

	client := body.Client
	email := body.Email
	password := body.Password

	if client == "" || email == "" || password == "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	trimmedClient := strings.TrimSpace(client)
	trimmedEmail := strings.TrimSpace(email)
	trimmedPassword := strings.TrimSpace(password)

	if trimmedClient == "" || trimmedEmail == "" || trimmedPassword == "" {
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

	clients := gohelpers.ObjectValues(configuration.Clients)
	if !gohelpers.IncludesString(clients, trimmedClient) {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InvalidData,
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

	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)

	rawPasswordRecord := PasswordCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
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
		trimmedPassword,
		passwordRecord.Hash,
	)
	if !passwordIsValid || comparisonError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.AccessDenied,
			Status: fiber.StatusUnauthorized,
		})
	}

	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)

	rawUserSecretRecord := UserSecretCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "userId", Value: userRecord.ID}},
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

	expiration, expirationError := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION"))
	if expirationError != nil {
		expiration = configuration.DefaultTokenExpiration
	}
	token, tokenError := utilities.GenerateJWT(utilities.GenerateJWTParams{
		Client:    trimmedClient,
		ExpiresIn: int64(expiration),
		Secret:    userSecretRecord.Secret,
		UserId:    userRecord.ID,
	})
	if tokenError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	redisError := utilities.RedisClient.Set(
		context.Background(),
		utilities.KeyFormatter(
			configuration.Redis.Prefixes.Secret,
			userRecord.ID,
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
			"token": token,
			"user":  userRecord,
		},
	})
}
