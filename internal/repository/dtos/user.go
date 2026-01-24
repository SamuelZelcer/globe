package dtos

type UpdateUsernameRequest struct {
	Username string `json:"username"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Passowrd string `json:"passowrd"`
}