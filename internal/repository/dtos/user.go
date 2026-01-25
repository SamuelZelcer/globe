package dtos

type UpdateUsernameRequest struct {
	Username string `json:"username"`
	RefreshToken string `json:"refreshToken"`
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}

type VerificationRequest struct {
	Code string `json:"code"`
	RefreshToken string `json:"refreshToken"`
}

type NewEmailAndVerificationCode struct {
	Code string `json:"code"`
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	NewPassowrd string `json:"newPassword"`
	Passowrd string `json:"password"`
	RefreshToken string `json:"refreshToken"`
}

type NewPasswordAndVerificationCode struct {
	Code string `json:"code"`
	NewPassowrd string `json:"newPassword"`
}