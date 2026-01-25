package dtos

type SignUpRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type SignUpVerification struct {
	Code string `json:"code"`
}

type SignInRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}