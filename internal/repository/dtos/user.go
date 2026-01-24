package dtos

type UpdateUsernameRequest struct {
	Username string `json:"username"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}

type VerifyNewEmailRequest struct {
	Code string `json:"code"`
	Email string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}

type NewEmailAndVerificationCode struct {
	Code string `json:"code"`
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Passowrd string `json:"passowrd"`
	RefreshToken string `json:"refreshToken"`
}