package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
 * Structs
 */
type Username struct {
	Username string `json:"username"`
}

type NewPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword 	string `json:"new_password"`
}

type SessionUser struct {
	Username string `json:"username"`
	IsAdmin  bool	`json:"is_admin"`
}

type User struct {
	UserId 			string 	`json:"user_id"`
	Username 		string  `json:"username"`
	Email 			string  `json:"email"`
	HashedPassword 	string  `json:"hashed_password"`
	IsAdmin 		bool 	`json:"is_admin"`
}

type NewUser struct {
	Token 	 string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
 * Router Methods
 */

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

func createNewUser(c *gin.Context) {
    var newUser NewUser
    if err := c.BindJSON(&newUser); err != nil {
		abortRouterMethod(c, http.StatusBadRequest, "Could not bind json", err)
        return
    }

    email, err := validateVerificationToken(newUser.Token)
    if err != nil {
        abortRouterMethod(c, http.StatusBadRequest, "Could not decode token", err)
        return
    }

    userId, err := addNewUserIntoDb(
        newUser.Username, email, newUser.Password,
    )
    if err != nil {
        abortRouterMethod(c, http.StatusBadRequest, "Could not add user into db", err)
        return
    }

	accessToken, err := createAccessToken(userId)
	if err != nil {
		abortRouterMethod(c, http.StatusUnauthorized, "Could not create access token", err)
		return
	}

	if err := setRefreshTokenCookie(c, userId); err != nil{
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"access_token": accessToken,
		"token_type": "Bearer",
	})
}

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

func deleteCurrentUser(c *gin.Context) {
    currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err := deleteUserByUserId(currentUserId); err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "detail":err,
        })
        return
    }

	deleteCookies(c)
    c.Status(http.StatusNoContent)
}

func editUsername(c *gin.Context) {
    var username Username
    if err := c.BindJSON(&username); err != nil {
        log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
        return
    }

    currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err := updateUsernameByUserId(
        currentUserId,
        username.Username,
    ); err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "detail":"Could not update username",
        })
        return
    }

    c.Status(http.StatusAccepted)
}

func editPassword(c *gin.Context) {
    var newPassword NewPassword
    if err := c.BindJSON(&newPassword); err != nil {
        log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
        return
    }

    currentUserId, err := getCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err = updatePasswordByUserId(
        currentUserId,
        newPassword.CurrentPassword,
        newPassword.NewPassword,
    ); err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "detail":"Could not update password",
        })
        return
    }

    c.Status(http.StatusAccepted)
}

/*
 * Services
 */

func addNewUserIntoDb(username string, email string, password string) (string, error) {
	// If the new user is the first in the db, they are an admin
	adminInDB, err := driver.ValueExists("users", "is_admin", true)
	if err != nil {
		log.Println(err)
		return "", err
	}
	isAdmin := !adminInDB

	if _, err = verifyEmailIsUnique(email); err != nil {
		log.Println(err)
		return "", err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	userId := createUserId()

	if insertUserIntoDb(
		userId, username, email, hashedPassword, isAdmin,
	); err != nil {
		log.Println(err)
		return "", err
	}

	if err := populatePredictionsByUserId(userId); err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("New user created")
	return userId, nil
}

func createUserId() string {
	return uuid.New().String()
}

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
	} else if otherAdminExists {
		return errors.New("no other admins exist")
	}

	if result, err := driver.Exec(
		"DELETE FROM users WHERE user_id = $1",
		userId,
	); err != nil {
		return err
	} else {
		if _, err:= result.RowsAffected(); err != nil {
			return err
		} 
	}

	log.Println("User deleted")
	return nil
}

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

func getUserPasswordHash(userId string) (string, error) {
	var hash string
	err := driver.QueryRow("SELECT hashed_password FROM users WHERE user_id = $1", userId).Scan(&hash)
	return hash, err
}

func insertUserIntoDb(userId string, username string, email string, hashedPassword string, isAdmin bool) error {
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

func updateUsernameByUserId(userId string, username string) error {
	result, err := driver.Exec(
		"UPDATE users SET username = $1 WHERE user_id = $2",
		username,
		userId,
	)
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 

	log.Println("Username updated")
	return nil
}

func updatePasswordByUserId(userId string, oldPassword string,	newPassword string) error {
	oldPasswordHash, err := getUserPasswordHash(userId)
	if err != nil{		
		log.Println(err)
		return err
	}

	if !verifyPasswordHash(oldPasswordHash, oldPassword) {
		err := errors.New("old password was not correct")
		log.Println(err)
		return err
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := driver.Exec(
		"UPDATE users SET hashed_password = $1 WHERE user_id = $2",
		hashedPassword,
		userId,
	)
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 

	log.Println("Password changed")
	return nil
}

func validateVerificationToken(token string) (string, error) {
	// Retrieves email from redis using the token
	if email, err := redis.Get(redisContext, token).Result(); err != nil {
		return "", err
	} else {
		return email, nil
	}
}

func verifyEmailIsUnique(email string) (bool, error) {
	emailExists, err := driver.ValueExists("users", "email", email)
	if err != nil {
		return false, err
	}
	if emailExists {
		err = errors.New("email already exists in db")
		return false, err
	}
	return true, nil
}

/*
*---------------------------------------------------------- 
* Only for testing
*----------------------------------------------------------
*/
type TestNewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email 	 string `json:"email"`
}

func testCreateUser(c *gin.Context) {
	var newUser TestNewUser
	if err := c.BindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}

	if _, err := addNewUserIntoDb(
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