package emailClient

const (
	EMAIL_SUBJECT string = "Points are Bad - Email Verification"
	EMAIL_BASE    string = `Welcome to Points are Bad! Please verify your email using the link below, and create your account:
	`
)

func createEmailBody(token string) string {
	return EMAIL_BASE + SIGNUP_TOKEN_URL + token
}

func createEmailHtml(token string) string {
	return "<p>" + EMAIL_BASE + "<a href=\"" + SIGNUP_TOKEN_URL + token + "\">" + "Verify Email</a></p>"
}