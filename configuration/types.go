package configuration

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

type ResponseMessagesStruct struct {
	EmailAlreadyInUse   string
	InternalServerError string
	InvalidData         string
	InvalidEmail        string
	InvalidToken        string
	MissingData         string
	Ok                  string
}

type RolesStruct struct {
	Admin string
	User  string
}
