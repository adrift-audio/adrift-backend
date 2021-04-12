package configuration

var Clients = ClientsStruct{
	Desktop: "desktop",
	Mobile:  "mobile",
	Web:     "web",
}

var Environments = EnvironmentsStruct{
	Development: "development",
	Heroku:      "heroku",
	Production:  "production",
}

var ResponseMessages = ResponseMessagesStruct{
	EmailAlreadyInUse:   "EMAIL_IS_ALREADY_IN_USE",
	InternalServerError: "INTERNAL_SERVER_ERROR",
	InvalidData:         "INVALID_DATA",
	InvalidEmail:        "INVALID_EMAIL",
	MissingData:         "MISSING_DATA",
	Ok:                  "OK",
}

var Roles = RolesStruct{
	Admin: "admin",
	User:  "user",
}
