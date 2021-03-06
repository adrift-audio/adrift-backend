package configuration

import "time"

var Clients = ClientsStruct{
	Desktop: "desktop",
	Mobile:  "mobile",
	Web:     "web",
}

var DefaultTokenExpiration = 99999

var Environments = EnvironmentsStruct{
	Development: "development",
	Heroku:      "heroku",
	Production:  "production",
}

var Redis = RedisStruct{
	Prefixes: RedisPrefixes{
		Room:   "room",
		Secret: "secret",
		User:   "user",
	},
	TTL: 24 * time.Hour,
}

var ResponseMessages = ResponseMessagesStruct{
	AccessDenied:         "ACCESS_DENIED",
	EmailAlreadyInUse:    "EMAIL_IS_ALREADY_IN_USE",
	InternalServerError:  "INTERNAL_SERVER_ERROR",
	InvalidData:          "INVALID_DATA",
	InvalidEmail:         "INVALID_EMAIL",
	InvalidRecoveryCode:  "INVALID_RECOVERY_CODE",
	InvalidToken:         "INVALID_TOKEN",
	InvalidUserID:        "INVALID_USER_ID",
	MissingData:          "MISSING_DATA",
	MissingPassphrase:    "MISSING_PASSPHRASE",
	MissingToken:         "MISSING_TOKEN",
	Ok:                   "OK",
	OldPasswordIsInvalid: "OLD_PASSWORD_IS_INVALID",
}

var Roles = RolesStruct{
	Admin: "admin",
	User:  "user",
}
