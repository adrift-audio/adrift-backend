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

func signUpController(ctx *fiber.Ctx) error {
	var body SignUpBodyStruct
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
	firstName := body.FirstName
	lastName := body.LastName
	password := body.Password
	signedAgreement := body.SignedAgreement

	if client == "" || email == "" || firstName == "" ||
		lastName == "" || password == "" || !signedAgreement {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.MissingData,
			Status: fiber.StatusBadRequest,
		})
	}

	trimmedClient := strings.TrimSpace(client)
	trimmedEmail := strings.TrimSpace(email)
	trimmedFirstName := strings.TrimSpace(firstName)
	trimmedLastName := strings.TrimSpace(lastName)
	trimmedPassword := strings.TrimSpace(password)

	if trimmedClient == "" || trimmedEmail == "" || trimmedFirstName == "" ||
		trimmedLastName == "" || trimmedPassword == "" {
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

	existingRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "email", Value: trimmedEmail}},
	)
	existingUser := &Schemas.User{}
	existingRecord.Decode(existingUser)
	if existingUser.ID != "" {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.EmailAlreadyInUse,
			Status: fiber.StatusBadRequest,
		})
	}

	now := utilities.MakeTimestamp()
	NewUser := new(Schemas.User)
	NewUser.ID = ""
	NewUser.Email = trimmedEmail
	NewUser.FirstName = trimmedFirstName
	NewUser.LastName = trimmedLastName
	NewUser.Role = configuration.Roles.User
	NewUser.SignedAgreement = true
	NewUser.Created = now
	NewUser.Updated = now
	insertionResult, insertionError := UserCollection.InsertOne(ctx.Context(), NewUser)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}
	createdRecord := UserCollection.FindOne(
		ctx.Context(),
		bson.D{{Key: "_id", Value: insertionResult.InsertedID}},
	)
	createdUser := &Schemas.User{}
	createdRecord.Decode(createdUser)

	UserSecretCollection := DB.Instance.Database.Collection(DB.Collections.UserSecret)

	secret, secretError := utilities.MakeHash(
		createdUser.ID + fmt.Sprintf("%v", utilities.MakeTimestamp()),
	)
	if secretError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	NewUserSecret := new(Schemas.UserSecret)
	NewUserSecret.ID = ""
	NewUserSecret.Secret = secret
	NewUserSecret.UserId = createdUser.ID
	NewUserSecret.Created = now
	NewUserSecret.Updated = now
	_, insertionError = UserSecretCollection.InsertOne(ctx.Context(), NewUserSecret)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	PasswordCollection := DB.Instance.Database.Collection(DB.Collections.Password)

	hash, hashError := utilities.MakeHash(trimmedPassword)
	if hashError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	NewPassword := new(Schemas.Password)
	NewPassword.ID = ""
	NewPassword.Hash = hash
	NewPassword.RecoveryCode = ""
	NewPassword.UserId = createdUser.ID
	NewPassword.Created = now
	NewPassword.Updated = now
	_, insertionError = PasswordCollection.InsertOne(ctx.Context(), NewPassword)
	if insertionError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	expiration, expirationError := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION"))
	if expirationError != nil {
		expiration = configuration.DefaultTokenExpiration
	}
	token, tokenError := utilities.GenerateJWT(utilities.GenerateJWTParams{
		Client:    trimmedClient,
		ExpiresIn: int64(expiration),
		Secret:    secret,
		UserId:    createdUser.ID,
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
			createdUser.ID,
		),
		secret,
		configuration.Redis.TTL,
	).Err()
	if redisError != nil {
		return utilities.Response(utilities.ResponseParams{
			Ctx:    ctx,
			Info:   configuration.ResponseMessages.InternalServerError,
			Status: fiber.StatusInternalServerError,
		})
	}

	emailTemplate := utilities.CreateWelcomeTemplate(
		createdUser.FirstName,
		createdUser.LastName,
	)
	utilities.SendEmail(
		createdUser.Email,
		emailTemplate.Subject,
		emailTemplate.Message,
	)

	return utilities.Response(utilities.ResponseParams{
		Ctx: ctx,
		Data: fiber.Map{
			"token": token,
			"user":  createdUser,
		},
	})
}
