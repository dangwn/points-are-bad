package api

import (
	"errors"
	"math/rand"
	"net/http"

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

// Authentication router group
func (r Router) addAuthGroup(rg *gin.RouterGroup) {
    auth := rg.Group("/auth")
    auth.POST("/login/", login)
    auth.DELETE("/login/", logout)
    auth.POST("/refresh/", refreshAccessToken)
    auth.POST("/verify/", verifyNewUserEmail)
}

/*
 * Login Endpoint
 * Schema: {email: string, password: string}
 * Verifies user email and password, and returns access tokens + token cookies
 */
func login(c *gin.Context) {
    var loginUser LoginUser
    if err := c.BindJSON(&loginUser); err != nil {
        logMessage := "Error binding user json to variable in login: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Bad request", logMessage)
        return
    }
    
    userId, username, isAdmin, err := validateLoginUser(
        loginUser.Email,
        loginUser.Password,
    )
    if err != nil {
        logMessage := "Error validating login user in login: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Incorrect email or password", logMessage)
        return
    }

    accessToken, err := createAccessToken(userId, username, isAdmin)
    if err != nil {
        logMessage := "Error creating access token in login: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not create access token", logMessage)
        return
    }

    if err := setRefreshTokenCookie(c, userId); err != nil{
        logMessage := "Error setting refresh token in login: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not set refresh token", logMessage) 
        return
    }

    c.JSON(http.StatusAccepted, Token{AccessToken: accessToken, TokenType: "Bearer"})
    // c.JSON(http.StatusAccepted, gin.H{
    //     "access_token": accessToken,
    //     "token_type": "Bearer",
    // })
    Logger.Info("User "+ userId + " successfully logged in")
}

/*
 * Logout Endpoint
 * Verifies current user, then deletes token cookies
 */
func logout(c *gin.Context) {
    userId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not retrieve userId in logout: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not log user out", logMessage)
        return
    }
    
    deleteCookies(c)
    c.Status(http.StatusNoContent)
    Logger.Info("User " + userId + " successfully logged out")
}

/*
 * Token Refresh Endpoint
 * Retrieves refresh and extracts userId from token and returns new tokens and cookies 
 */
func refreshAccessToken(c *gin.Context) {
    // Get refresh token from cookie
    cookie, err := c.Request.Cookie(REFRESH_TOKEN_NAME)
    if err != nil {
        logMessage := "Could not get refresh token cookie in refreshAccessToken: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not get refresh token", logMessage)
        return
    }

    // Get user Id from refresh token
    userId, err := jwtDecode(cookie.Value, REFRESH_TOKEN_SECRET_KEY)
    if err != nil {
        logMessage := "Could not get userId from refresh token in refreshAccessToken: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not decode refresh token", logMessage)
        return
    }

    // Get user's username and admin status
    user, err := getUserByUserId(userId)
    if err != nil {
        logMessage := "Could not retrieve user information in refreshAccessToken: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve user information", logMessage)
        return
    }

    // Create access token
    accessToken, err := createAccessToken(userId, user.Username, user.IsAdmin)
    if err != nil {
        logMessage := "Could not create access token in refreshAccessToken: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not create new access token", logMessage)
        return
    }

    // Set refresh token
    if err := setRefreshTokenCookie(c, userId); err != nil {
        logMessage := "Could not set refresh token cookie in refreshAccessToken: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not set refresh token", logMessage)
        return
    }

    // Return access token
    c.JSON(http.StatusAccepted, Token{AccessToken: accessToken, TokenType: "Bearer"})
    // c.JSON(http.StatusAccepted, gin.H{
    //     "access_token": accessToken,
    //     "token_type": "Bearer",
    // })
    Logger.Info("User " + userId + " refreshed tokens")
}

/*
 * Email Verification Endpoint
 * Schema: {email: string}
 * Verifies email format, and that the email doesn't already exist in database
 * Sends message to email server via RabbitMQ to email new user with verification link
 */
func verifyNewUserEmail(c *gin.Context) {
    var newUserEmail EmailAddress
    if err := c.BindJSON(&newUserEmail); err != nil {
        logMessage := "Could not bind query in verifyNewUserEmail: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Bad request", logMessage)
        return
    }

    if !verifyEmailFormat(newUserEmail.Email) {
        logMessage := "Email address format not valid in verifyNewUserEmail"
        abortRouterMethod(c, http.StatusBadRequest, "Email address not valid", logMessage)
        return
    }

    if emailInDb, err := isEmailInDb(newUserEmail.Email); err != nil {
        logMessage := "Could not verify if email is in db in verifyNewUserEmail: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not get refresh token", logMessage)
        return
    } else if emailInDb {
        abortRouterMethod(c, http.StatusUnauthorized, "Email address already taken")
        return
    }

    if err := sendVerifyMessageToEmailClient(newUserEmail.Email); err != nil {
        Logger.Warning("Could not send email verification link in verifyUserEmail" + err.Error())
    } else {
        Logger.Info("Verification link sent")
    }

    // Failures in sending tokens to rabbitmq will not be shown to the user
    c.Status(http.StatusAccepted)
}

/*
 * Services
 */

// Creates a random string of letters of a given length
func createRandomString(length int) string {
    b := make([]rune, length)
    numChars := len(CHARS)
    
    for i := range b {
        b[i] = CHARS[rand.Intn(numChars)]
    }
    return string(b)
}

// Creates a verification token (20 random letters)
func createVerificationToken() string {
    return createRandomString(20)
}

/*
 * Retrieves the current user using the incoming context
 * Decodes the access token and returns the userId of the current user
 */
func getCurrentUser(c *gin.Context) (string, error) {
    token, err := getAuthorizationToken(c)
    if err != nil {
        return "", err
    }

    if userId, err := jwtDecode(token, ACCESS_TOKEN_SECRET_KEY); err != nil {
        return "", err
    } else {
        return userId, nil
    }
}

/*
 * Returns whether the current user is an admin
 * If no current user can be found, it returns false
 */
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

// Returns whether a given email is in the db
func isEmailInDb(email string) (bool, error) {
    return driver.ValueExists("users", "email",    email)
}

/*
 * Sends a message to the email server via RabbitMQ of the form:
 *        {"email":"<email address>","token":"<verification token>"}
 * for the email server to send a verification link to a new user
 */
func sendVerifyMessageToEmailClient(email string) error {
    var token string
    // Check to see if token is unique, if it isn't create a new one
    for {
        token = createVerificationToken()
        exists, err := redisKeyExists(token)
        if err != nil {
            return err
        }
        if !exists {
            break
        }
    }

    // Store the email in redis to be verified by create user function
    if status, err := redis.Set(redisContext, token, email, REDIS_DURATION).Result(); err != nil {
        return err
    } else if status != "OK" {
        return errors.New("redis returned non-OK status code")
    }

    // Send email and token to RabbitMQ for email server
    msg := `{"email":"` + email + `","token":"` + token + `"}`
    rabbit.SendMessage(msg,    EMAIL_VERIFICATION_QUEUE, 10)
    return nil
}

/*
 * Validates a user's login credentials against the db
 * Returns the user's userId string
 */
func validateLoginUser(email, password string) (string, string, bool, error) {
    var userId, username, hashedPassword string
    var isAdmin bool

    if err := driver.QueryRow(
        "SELECT user_id, username, hashed_password, is_admin FROM users WHERE email = $1",
        email,
    ).Scan(&userId, &username, &hashedPassword, &isAdmin); err != nil {
        return "", "", false, err
    }

    if !verifyPasswordHash(hashedPassword, password) {
        return "", "", false, errors.New("incorrect username or password")
    }

    return userId, username, isAdmin, nil
}
