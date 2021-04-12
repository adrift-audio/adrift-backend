package auth

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
