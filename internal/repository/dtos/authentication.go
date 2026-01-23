package dtos

type AuthenticationTokens struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken string `json:"accessToken"`
}

type UpdateAuthTokensRequest struct {
	RefreshToken string `json:"refreshToken"`
}