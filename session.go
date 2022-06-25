package token

type Session struct {
	TokenType    string `json:"token_type"` // Bearer
	ExpiresIn    uint   `json:"expires_in"` // The lifetime in seconds of refresh token
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
