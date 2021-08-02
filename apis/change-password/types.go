package changePassword

type ChangePasswordBodyStruct struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
