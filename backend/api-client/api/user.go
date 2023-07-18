package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * Structs
 */
type Username struct {
    Username string `json:"username"`
}

type ComparePassword struct {
    CurrentPassword string `json:"old_password"`
    NewPassword     string `json:"new_password"`
}

type SessionUser struct {
    Username    string `json:"username"`
    IsAdmin     bool   `json:"is_admin"`
}

type User struct {
    UserId          string  `json:"user_id"`
    Username        string  `json:"username"`
    Email           string  `json:"email"`
    HashedPassword  string  `json:"hashed_password"`
    IsAdmin         bool    `json:"is_admin"`
}

type NewUser struct {
    Token       string `json:"token"`
    Username    string `json:"username"`
    Password    string `json:"password"`
}

/*
 * Router Methods
 */

// User router group
func (r Router) addUserGroup(rg *gin.RouterGroup) {
    user := rg.Group("/user")
    user.GET("/", displayCurrentUser)
    user.POST("/", createNewUser)
    user.DELETE("/", deleteCurrentUser)
    user.PUT("/username/", editUsername)
    user.PUT("/password/", editPassword)

    if IS_DEV_BUILD {
        user.POST("/testCreateUser/", testCreateUser)
    }
}

/*
 * User Creation Endpoint
 * Schema: {token: string, username: string, password: string}
 * Returns: {access_token: string, token_type: string}
 * Creates a new user in the db and returns the new user's access token
 *     and refresh token cookie
 * The token is taken from the verification link sent to the user
 */
func createNewUser(c *gin.Context) {
    var newUser NewUser
    if err := c.BindJSON(&newUser); err != nil {
        logMessage := "Could not bind json in createNewUser: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not bind json", logMessage)
        return
    }

    email, err := validateVerificationToken(newUser.Token)
    if err != nil {
        logMessage := "Could not validate verification token in createNewUser: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not decode token", logMessage)
        return
    }

    userId, isAdmin, err := addNewUserIntoDb(newUser.Username, email, newUser.Password)
    if err != nil {
        logMessage := "Error adding user to db in createNewUser: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not create user", logMessage)
        return
    }

    accessToken, err := createAccessToken(userId, newUser.Username, isAdmin)
    if err != nil {
        logMessage := "Error creating access token in createNewUser: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not create access token", logMessage)
        return
    }

    if err := setRefreshTokenCookie(c, userId); err != nil {
        logMessage := "Error setting refresh token in createNewUser: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not set refresh token", logMessage)
        return
    }
    
    c.JSON(http.StatusCreated, Token{AccessToken: accessToken, TokenType: "Bearer"})
    Logger.Info("User " + userId + " created")
}

/*
 * [DEPRECEATED] Current User Endpoint
 * Returns: {username: string, is_admin: bool}
 * Do not use this endpoint, user's username and admin status are
 *    now sent via the access token
 */
func displayCurrentUser(c *gin.Context) {
    currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    user, err := getUserByUserId(currentUserId)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Not authorized",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "username": user.Username,
        "is_admin": user.IsAdmin,
    })
}

/*
 * Delete User Endpoint
 * Deletes the current user
 * Will not delete user if they are the only admin user
 */
func deleteCurrentUser(c *gin.Context) {
    currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not get current user in deleteCurrentUser: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve current user", logMessage)
        return
    }

    if err := deleteUserByUserId(currentUserId); err != nil {
        logMessage := "Could not delete user user in deleteCurrentUser: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, err, logMessage)
        return
    }

    deleteCookies(c)
    c.Status(http.StatusNoContent)
    Logger.Info("User " + currentUserId + " deleted")
}

/*
 * Edit Username Endpoint
 * Schema: {username: string}
 * Edits the current user's username
 */
func editUsername(c *gin.Context) {
    var username Username

    currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not get current user in editUsername: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve current user", logMessage)
        return
    }

    if err = c.BindJSON(&username); err != nil {
        logMessage := "Could not bind JSON in editUsername: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve new username", logMessage)
        return
    }


    if err = updateUsernameByUserId(currentUserId, username.Username); err != nil {
        logMessage := "Could not update username in editUsername: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not update username", logMessage)
        return
    }

    c.Status(http.StatusAccepted)
    Logger.Info("User " + currentUserId + " updated their username to " + username.Username)
}

/*
 * Edit Password Endpoint
 * Schema: {current_password: string, new_password: string}
 * Validates the current user's current password, and changes it to a provided new one
 */ 
func editPassword(c *gin.Context) {
    var comparePassword ComparePassword

    currentUserId, err := getCurrentUser(c)
    if err != nil {
        logMessage := "Could not get current user in editPassword: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not retrieve current user", logMessage)
        return
    }

    if err = c.BindJSON(&comparePassword); err != nil {
        logMessage := "Could not bind JSON in editPassword: " + err.Error()
        abortRouterMethod(c, http.StatusBadRequest, "Could not retrieve passwords", logMessage)
        return
    }

    if err = updatePasswordByUserId(currentUserId, comparePassword.CurrentPassword, comparePassword.NewPassword); err != nil {
        logMessage := "Could not update password in editPassword: " + err.Error()
        abortRouterMethod(c, http.StatusUnauthorized, "Could not update password", logMessage)
        return
    }

    c.Status(http.StatusAccepted)
    Logger.Info("User " + currentUserId + " updated their password")
}

/*
 * Services
 */

/*
 * Adds a new user into the db
 * Every user's email is unique
 * By default no new users are admins, except the very first user
 */
func addNewUserIntoDb(username, email, password string) (string, bool, error) {
    // If the new user is the first in the db, they are an admin
    adminInDB, err := driver.ValueExists("users", "is_admin", true)
    if err != nil {
        return "", false, err
    }
    isAdmin := !adminInDB

    // Verify email is unique
    if emailIsUnique, err := verifyEmailIsUnique(email); err != nil {
        return "", false, err
    } else if !emailIsUnique {
        return "", false, errors.New("email already in use")
    }

    // Create user Id and password hash
    hashedPassword, err := hashPassword(password)
    if err != nil {
        return "", false, err
    }
    userId := createUUID()

    // Insert user into db and populate the predictions table with their predictions
    if err := insertUserIntoDb(userId, username, email, hashedPassword, isAdmin); err != nil {
        return "", false, err
    }
    if err := populatePredictionsByUserId(userId); err != nil {
        return "", false, err
    }

    return userId, isAdmin, nil
}

/*
 * Deletes a user by their user Id
 * Will not delete user if they are the only admin user
 */
func deleteUserByUserId(userId string) error {
    var otherAdminExists bool
    adminExistsQuery := `
        SELECT EXISTS(
            SELECT 1 FROM users
            WHERE user_id != $1 AND
                is_admin = true
        )
    `

    if err := driver.QueryRow(adminExistsQuery, userId).Scan(&otherAdminExists); err != nil {
        return err
    } else if !otherAdminExists {
        return errors.New("no other admins exist")
    }

    if result, err := driver.Exec(
        "DELETE FROM users WHERE user_id = $1",
        userId,
    ); err != nil {
        return err
    } else {
        if n, err:= result.RowsAffected(); err != nil {
            return err
        } else if n != 1 {
            return errors.New("incorrect number of rows affected when deleting user")
        }
    }

    return nil
}

// Retrieves a user's information by their user Id
func getUserByUserId(userId string) (User, error) {
    var user User
    userQuery := `
        SELECT user_id, username, email, hashed_password, is_admin 
        FROM users 
        WHERE user_id = $1
    `

    err := driver.QueryRow(userQuery, userId).Scan(
        &user.UserId,
        &user.Username,
        &user.Email,
        &user.HashedPassword,
        &user.IsAdmin,
    )
    return user, err
}

// Retrieves a user's password hash
func getUserPasswordHash(userId string) (string, error) {
    var hash string
    err := driver.QueryRow("SELECT hashed_password FROM users WHERE user_id = $1", userId).Scan(&hash)
    return hash, err
}

// Inserts a new user into the db
func insertUserIntoDb(userId, username, email, hashedPassword string, isAdmin bool) error {
    _, err := driver.Exec(
        "INSERT INTO users VALUES($1, $2, $3, $4, $5, 0, 0, 0, NULL)",
        userId,
        username,
        email,
        hashedPassword,
        isAdmin,
    )
    return err
}

// Updates a user's username
func updateUsernameByUserId(userId, username string) error {
    result, err := driver.Exec(
        "UPDATE users SET username = $1 WHERE user_id = $2",
        username,
        userId,
    )
    if err != nil {
        return err
    } 

    if n, err := result.RowsAffected(); err != nil {
        return err
    } else if n != 1 {
        return errors.New("incorrect number of rows affected when updating username")
    }

    return nil
}

/*
 * Updates a user's hashed password in db
 * User's previous must be supplied and verified against password currently in db
 */
func updatePasswordByUserId(userId, oldPassword, newPassword string) error {
    oldPasswordHash, err := getUserPasswordHash(userId)
    if err != nil{
        return err
    }

    if !verifyPasswordHash(oldPasswordHash, oldPassword) {
        err := errors.New("previous password was not verfied against password in db")
        return err
    }

    hashedPassword, err := hashPassword(newPassword)
    if err != nil {
        return err
    }

    result, err := driver.Exec(
        "UPDATE users SET hashed_password = $1 WHERE user_id = $2",
        hashedPassword,
        userId,
    )
    if err != nil {
        return err
    } 

    if n, err:= result.RowsAffected(); err != nil {
        return err
    } else if n != 1 {
        return errors.New("incorrect number of rows affected when upoda")
    }

    return nil
}

/*
 * Validates verification token by searching for a corresponding value
 *    to the token as a key in redis
 */
func validateVerificationToken(token string) (string, error) {
    if email, err := redis.Get(redisContext, token).Result(); err != nil {
        return "", err
    } else {
        return email, nil
    }
}

// Verifies if a given email is unique by checking against the db
func verifyEmailIsUnique(email string) (bool, error) {
    if emailExists, err := driver.ValueExists("users", "email", email); err != nil {
        return false, err
    } else if emailExists {
        return false, nil
    }
    return true, nil
}

/*
*---------------------------------------------------------- 
* Only for testing
*----------------------------------------------------------
*/
type TestNewUser struct {
    Username    string `json:"username"`
    Password    string `json:"password"`
    Email       string `json:"email"`
}

func testCreateUser(c *gin.Context) {
    var newUser TestNewUser
    if err := c.BindJSON(&newUser); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Bad request",
        })
        return
    }
    
    if _, _, err := addNewUserIntoDb(
        newUser.Username,
        newUser.Email,
        newUser.Password,
    ); err != nil {
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "detail": err,
        })
        return
    }
}