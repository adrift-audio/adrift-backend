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
	AccessDenied:        "ACCESS_DENIED",
	EmailAlreadyInUse:   "EMAIL_IS_ALREADY_IN_USE",
	InternalServerError: "INTERNAL_SERVER_ERROR",
	InvalidData:         "INVALID_DATA",
	InvalidEmail:        "INVALID_EMAIL",
	InvalidToken:        "INVALID_TOKEN",
	MissingData:         "MISSING_DATA",
	MissingToken:        "MISSING_TOKEN",
	Ok:                  "OK",
}

var Roles = RolesStruct{
	Admin: "admin",
	User:  "user",
}
