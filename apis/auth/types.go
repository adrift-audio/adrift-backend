package auth

type GetRecoveryCodeBodyStruct struct {
	Email string `json:"email"`
}

type RecoverAccounntBodyStruct struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type SignInBodyStruct struct {
	Client   string `json:"client"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpBodyStruct struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	SignedAgreement bool   `json:"signedAgreement"`
	SignInBodyStruct
}
