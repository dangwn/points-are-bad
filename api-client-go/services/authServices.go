package services

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
)

func createRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LETTERS[rand.Intn(len(LETTERS))]
	}
	return string(b)
}

func getAuthorizationToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header not in header")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return "", errors.New("incorrect header format")
	}

	return authParts[1], nil
}

func GetCurrentUser(c *gin.Context) (string, error) {
	token, err := getAuthorizationToken(c)
	if err != nil {
		return "", err
	}

	userId, err := JwtDecode(
		token,
		ACCESS_TOKEN_SECRET_KEY,
	)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return userId, nil
}

func IsCurrentUserAdmin(c *gin.Context) bool {
	userId, err := GetCurrentUser(c)
	if err != nil {
		return false
	}

	user, err := GetUserById(userId)
	if err != nil {
		return false
	}
	return user.IsAdmin
}

func IsEmailInDB(email string) (bool, error) {
	return driver.ValueExists(
		"users",
		"email",
		email,
	)
}

func SendVerifyMessageToEmailClient(email string) {
	code := createRandomString(16)
	msg := "{\"verify--" + email + "\":\"" + code + "\"}"

	rabbit.SendMessage(
		msg,
		EMAIL_VERIFICATION_QUEUE,
		10,
	)
}

func ValidateLoginUser(email string, password string) (string, error) {
	var userId string
	var hashedPassword string
	emailInDb := false
	rows, err := driver.Query(
		"SELECT user_id, hashed_password FROM users WHERE email = $1 LIMIT 1",
		email,
	)
	if err != nil{
		return "", err
	}
	
	for rows.Next() {
		if err := rows.Scan(
			&userId,
			&hashedPassword,
		); err != nil {
			return "", err
		}
		emailInDb = true
		break
	}
	if verifyPasswordHash(hashedPassword, password) && emailInDb {
		return userId, nil
	}
	return "", errors.New("incorrect username or password")
}
