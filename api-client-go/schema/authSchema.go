package schema

type EmailAddress struct {
	Email string `json:"email"`
}

type LoginUser struct {
	EmailAddress
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
