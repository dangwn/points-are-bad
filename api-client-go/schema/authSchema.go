package schema

type LoginUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
}

type Email struct {
	Email string `json:"email"`
}