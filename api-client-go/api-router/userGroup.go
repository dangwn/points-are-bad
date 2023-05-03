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

    // Only for testing
    user.POST("/testCreateUser", testCreateUser)
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

    if err := services.CreateNewUser(
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