package dtos

type SignUpRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Code string `json:"code"`
}

type SignInRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}