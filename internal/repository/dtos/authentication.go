package dtos

type AuthenticationTokens struct {
	RefreshToken *string
	AccessToken *string
}

type UpdateAuthTokensRequest struct {
	RefreshToken *string `json:"refreshToken"`
}