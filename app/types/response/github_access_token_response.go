package response

type GithubAccessTokenResponse struct {
	AccessToken           int    `json:"access_token"`
	ExpiresIn             string `json:"expires_in"`
	RefreshToken          int    `json:"refresh_token"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	Scope                 int    `json:"scope"`
	TokenType             string `json:"token_type"`
}
