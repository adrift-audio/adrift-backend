package schemas

type Password struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Hash         string `json:"hash"`
	RecoveryCode string `json:"recoveryCode" bson:"recoveryCode"`
	UserId       string `json:"userId" bson:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
