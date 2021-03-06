package configuration

import "time"

type ClientsStruct struct {
	Desktop string
	Mobile  string
	Web     string
}

type EnvironmentsStruct struct {
	Development string
	Heroku      string
	Production  string
}

type RedisPrefixes struct {
	Room   string
	Secret string
	User   string
}

type RedisStruct struct {
	Prefixes RedisPrefixes
	TTL      time.Duration
}
type ResponseMessagesStruct struct {
	AccessDenied         string
	EmailAlreadyInUse    string
	InternalServerError  string
	InvalidData          string
	InvalidEmail         string
	InvalidRecoveryCode  string
	InvalidToken         string
	InvalidUserID        string
	MissingData          string
	MissingPassphrase    string
	MissingToken         string
	Ok                   string
	OldPasswordIsInvalid string
}

type RolesStruct struct {
	Admin string
	User  string
}
