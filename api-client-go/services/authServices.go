package services

import (
	"strings"
	"strconv"
	"errors"
	"math/rand"
	
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) (int, error) {
	token, err := getAuthorizationToken(c)
	if err != nil {
		return -1, err
	}

	decodedToken, err := JwtDecode(
		token,
		ACCESS_TOKEN_SECRET_KEY,
	)
	if err != nil {
		return -1, err
	}

	userId, err := strconv.Atoi(decodedToken)
	if err != nil {
		return -1, err
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

func getAuthorizationToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization header not in header")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != "Bearer" {
		return "", errors.New("Incorrect header format")
	}

	return authParts[1], nil
}

func ValidateLoginUser(
	email string,
	password string,
) (int, error) {
	var hashedPassword string
	var userId int
	emailInDb := false
	rows := driver.Query(
		"SELECT user_id, hashed_password FROM users WHERE email = $1 LIMIT 1",
		email,
	)
	for rows.Next() {
		if err := rows.Scan(
			&userId,
			&hashedPassword,
		); err != nil {
			return -1, err
		}
		emailInDb = true
		break
	}
	if verifyPasswordHash(hashedPassword, password) && emailInDb {
		return userId, nil
	}
	return -1, errors.New("Incorrect username or password")
}

func IsEmailInDB(
	email string,
) (bool, error) {
	return driver.ValueExists(
		"users",
		"email",
		email,
	)
}

func createRandomString(
	length int,
) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LETTERS[rand.Intn(len(LETTERS))]
	}
	return string(b)
}

func SendVerifyMessageToEmailClient(
	email string,
) {
	code := createRandomString(16)
	msg := "{\"" + email + "\":\"" + code + "\"}"

	rabbit.SendMessage(
		msg,
		EMAIL_VERIFICATION_QUEUE,
		10,
	)
}