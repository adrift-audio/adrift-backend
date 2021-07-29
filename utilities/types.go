package utilities

import (
	"github.com/gofiber/fiber/v2"
	JWT "github.com/golang-jwt/jwt"
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

type Template struct {
	Message string
	Subject string
}
