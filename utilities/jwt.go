package utilities

import (
	"errors"
	"strings"
	"time"

	JWT "github.com/golang-jwt/jwt"

	"adrift-backend/configuration"
)

func DecodePayload(token string) ([]byte, error) {
	partials := strings.Split(token, ".")
	if len(partials) != 3 {
		return nil, errors.New("Invalid token")
	}

	bytePayload, decodeError := JWT.DecodeSegment(partials[1])
	if decodeError != nil {
		return nil, decodeError
	}

	return bytePayload, nil
}

func GenerateJWT(params GenerateJWTParams) (string, error) {
	expiration := params.ExpiresIn * 24 * 60 * 60
	if expiration == 0 {
		expiration = int64(configuration.DefaultTokenExpiration) * 24 * 60 * 60
	}

	claims := JWTClaims{
		params.Client,
		params.UserId,
		JWT.StandardClaims{
			ExpiresAt: time.Now().Unix() + expiration,
		},
	}

	token := JWT.NewWithClaims(JWT.SigningMethodHS256, claims)

	signedToken, signingError := token.SignedString([]byte(params.Secret))
	if signingError != nil {
		return "", signingError
	}

	return signedToken, nil
}

func ParseClaims(token string, secret string) (*JWTClaims, error) {
	decoded, parsingError := JWT.ParseWithClaims(
		token,
		&JWTClaims{},
		func(decoded *JWT.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if parsingError != nil {
		return &JWTClaims{}, parsingError
	}

	if claims, ok := decoded.Claims.(*JWTClaims); ok && decoded.Valid {
		return claims, nil
	}
	return &JWTClaims{}, errors.New(configuration.ResponseMessages.InvalidToken)
}
