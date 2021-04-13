package utilities

import (
	JWT "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type GenerateJWTParams struct {
	Client    string
	ExpiresIn int64
	Secret    string
	UserId    string
}

type JWTClaims struct {
	Client string `json:"client"`
	UserId string `json:"userId"`
	JWT.StandardClaims
}

type ResponseParams struct {
	Ctx    *fiber.Ctx
	Data   interface{}
	Info   string
	Status int
}
