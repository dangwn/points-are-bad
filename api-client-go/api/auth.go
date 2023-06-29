package api

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
 * Structs
 */

type EmailAddress struct {
	Email string `json:"email" form:"email"`
}

type LoginUser struct {
	EmailAddress
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

/*
 * Router modules
 */

 
func (r Router) addAuthGroup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login/", login)
	auth.DELETE("/login/", logout)
	auth.POST("/refresh/", refreshAccessToken)
	auth.POST("/verify/", verifyNewUserEmail)
}

func login(c *gin.Context) {
	var loginUser LoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}
	
	userId, err := validateLoginUser(
		loginUser.Email,
		loginUser.Password,
	)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Incorrect email or password",
		})
		return
	}

	accessToken, err := createAccessToken(userId)
	if err != nil {
		abortRouterMethod(c, http.StatusUnauthorized, "Could not create access token", err)
		return
	}

	if err1, err2 := setRefreshTokenCookie(c, userId), setCSRFTokenCookie(c, userId); err1 != nil  || err2 != nil{
		return
	}
	
	c.JSON(http.StatusAccepted, gin.H{
		"access_token": accessToken,
		"token_type": "Bearer",
	})
}

func logout(c *gin.Context) {
	clearAuthCookies(c)
	c.Status(http.StatusNoContent)
}

func refreshAccessToken(c *gin.Context) {
	cookie, err := c.Request.Cookie(REFRESH_TOKEN_NAME)
	if err != nil {
		log.Println("Could not get refresh token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": "Could not get refresh token",
		})
		return
	}

	userId, err := jwtDecode(
		cookie.Value,
		REFRESH_TOKEN_SECRET_KEY,
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": "Could not decode refresh token",
		})
		return
	}

	accessToken, err := jwtEncode(
		userId,
		ACCESS_TOKEN_SECRET_KEY,
		ACCESS_TOKEN_EXPIRE_TIME,
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": "Could not create new access token",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"access_token": accessToken,
		"token_type": "Bearer",
	})
}

func verifyNewUserEmail(c *gin.Context) {
	var newUserEmail EmailAddress
	if err := c.BindQuery(&newUserEmail); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}

	if !verifyEmailFormat(newUserEmail.Email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Email address not valid",
		})
		return
	}

	emailInDb, _ := isEmailInDb(
		newUserEmail.Email,
	)

	if emailInDb {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail":"Email already taken",
		})
	}

	sendVerifyMessageToEmailClient(
		newUserEmail.Email,
	)
	c.Status(http.StatusAccepted)
}


/*
 * Services
 */

func createRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = LETTERS[rand.Intn(len(LETTERS))]
	}
	return string(b)
}

func createVerificationToken() string {
	return createRandomString(14)
}

func clearAuthCookies(c *gin.Context) {
	c.SetCookie(
		REFRESH_TOKEN_NAME,
		"",
		-1, // Unlimited age
		"", // Default path
		FRONTEND_DOMAIN, // Frontend domain
		true, // Secure cookie
		true, // HTTP only
	)
	c.SetCookie(
		CSRF_TOKEN_NAME,
		"",
		-1, // Unlimited age
		"", // Default path
		FRONTEND_DOMAIN, // Frontend domain
		true, // Secure cookie
		true, // HTTP only
	)
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

func getCurrentUser(c *gin.Context) (string, error) {
	token, err := getAuthorizationToken(c)
	if err != nil {
		return "", err
	}

	userId, err := jwtDecode(
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

func isCurrentUserAdmin(c *gin.Context) bool {
	userId, err := getCurrentUser(c)
	if err != nil {
		return false
	}

	user, err := getUserByUserId(userId)
	if err != nil {
		return false
	}
	return user.IsAdmin
}

func isEmailInDb(email string) (bool, error) {
	return driver.ValueExists(
		"users",
		"email",
		email,
	)
}

func sendVerifyMessageToEmailClient(email string) {
	var token string
	// Check to see if token is unique, if it isn't create a new one
	for {
		token = createVerificationToken()
		exists, err := redisKeyExists(token)
		if err != nil {
			log.Println("Error checking if redis token already exists: " + err.Error())
			return
		}
		if !exists {
			break
		}
	}

	// Store the email in redis to be verified by create user function
	status, err := redis.Set(redisContext, token, email, REDIS_DURATION).Result()

	// Send email and token to RabbitMQ for email server
	msg := `{"email":"` + email + `","token":"` + token + `"}`
	rabbit.SendMessage(
		msg,
		EMAIL_VERIFICATION_QUEUE,
		10,
	)
}

func validateLoginUser(email string, password string) (string, error) {
	var userId string
	var hashedPassword string

	if err := driver.QueryRow(
		"SELECT user_id, hashed_password FROM users WHERE email = $1",
		email,
	).Scan(&userId, &hashedPassword); err != nil {
		return "", err
	}

	if !verifyPasswordHash(hashedPassword, password) {
		return "", errors.New("incorrect username or password")
	}

	return userId, nil
}
