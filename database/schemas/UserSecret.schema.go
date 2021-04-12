package schemas

type UserSecret struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Secret string `json:"secret"`
	UserId string `json:"userId" bson:"userId"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}
