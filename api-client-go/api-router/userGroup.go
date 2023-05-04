package apiRouter

import (
    "net/http"
    "log"
    "github.com/gin-gonic/gin"
    "points-are-bad/api-client/schema"
    "points-are-bad/api-client/services"
)

func (r Router) addUserGroup(rg *gin.RouterGroup) {
    user := rg.Group("/user")
    user.GET("/", displayCurrentUser)
    user.PUT("/username", editUsername)
    user.PUT("/password", editPassword)
    user.DELETE("/", deleteCurrentUser)

    // Only for testing
    user.POST("/testCreateUser", testCreateUser)
}

func createNewUser(c *gin.Context) {
    var newUser schema.NewUser
    if err := c.BindJSON(&newUser); err != nil {
        log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
        return
    }

    email, err := services.ValidateVerificationToken(newUser.Token)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not decode token",
		})
        return
    }

    userId, err := services.CreateNewUser(
        email, newUser.Username, newUser.Password,
    )
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Could not create new user",
		})
        return
    }

    accessToken, err := services.JwtEncode(
		userId,
		services.ACCESS_TOKEN_SECRET_KEY,
		services.ACCESS_TOKEN_EXPIRE_TIME,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":err,
		})
		return
	}
	refreshToken, err := services.JwtEncode(
		userId,
		services.REFRESH_TOKEN_SECRET_KEY,
		services.REFRESH_TOKEN_EXPIRE_TIME,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":err,
		})
		return
	}

	c.SetCookie(
		"X-REFRESH-TOKEN",
		refreshToken,
		0, // Unlimited age
		"", // Default path
		FRONTEND_DOMAIN, // Frontend domain
		false, // Secure cookie
		true, // HTTP only
	)
	c.SetCookie(
		"X-CSRF-TOKEN",
		refreshToken,
		0, // Unlimited age
		"", // Default path
		FRONTEND_DOMAIN, // Frontend domain
		false, // Secure cookie
		true, // HTTP only
	)
	c.JSON(http.StatusCreated, gin.H{
		"access_token": accessToken,
		"token_type": "Bearer",
	})
}

func displayCurrentUser(c *gin.Context) {
    currentUserId, err := services.GetCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    user, err := services.GetUserById(currentUserId)
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
    currentUserId, err := services.GetCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err := services.DeleteUserById(currentUserId); err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not delete user",
        })
        return
    }

    c.Status(http.StatusNoContent)
}

func editUsername(c *gin.Context) {
    var username schema.Username
    if err := c.BindJSON(&username); err != nil {
        log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
        return
    }

    currentUserId, err := services.GetCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err := services.UpdateUsernameByUserId(
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
    var newPassword schema.NewPassword
    if err := c.BindJSON(&newPassword); err != nil {
        log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
        return
    }

    currentUserId, err := services.GetCurrentUser(c)
    if err != nil {
        log.Println(err)
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "detail":"Could not retreieve current user",
        })
        return
    }

    if err = services.UpdatePasswordByUserId(
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
*---------------------------------------------------------- 
* Only for testing
*----------------------------------------------------------
*/
func testCreateUser(c *gin.Context) {
    var newUser schema.TestNewUser
	if err := c.BindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}

    if _, err := services.CreateNewUser(
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