package schemas

type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Role            string `json:"role"`
	SignedAgreement bool   `json:"signedAgreement"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
