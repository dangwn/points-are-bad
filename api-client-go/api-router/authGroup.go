package apiRouter

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"points-areb-bad/api-client/schema"
	"points-areb-bad/api-client/services"
)

func (r Router) addAuthGroup(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/login", login)
	auth.DELETE("/login", logout)
	auth.POST("/refresh", refreshAccessToken)
	auth.POST("/verify", verifyNewUserEmail)
}

func login(c *gin.Context) {
	var loginUser schema.LoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}
	userId, err := services.ValidateLoginUser(
		loginUser.Email,
		loginUser.Password,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail":"Incorrect email or password",
		})
		return
	}

	accessToken, err := services.JwtEncode(
		strconv.Itoa(userId),
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
		strconv.Itoa(userId),
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
	c.JSON(http.StatusAccepted, gin.H{
		"access_token": accessToken,
		"token_type": "Bearer",
	})
}

func logout(c *gin.Context) {
	c.SetCookie(
		"X-REFRESH-TOKEN",
		"",
		-1, // Unlimited age
		"", // Default path
		FRONTEND_DOMAIN, // Frontend domain
		true, // Secure cookie
		true, // HTTP only
	)
	c.Status(http.StatusNoContent)
}

func refreshAccessToken(c *gin.Context) {
	cookie, err := c.Request.Cookie("X-Refresh-Token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": "Could not get refresh token",
		})
		return
	}

	userId, err := services.JwtDecode(
		cookie.Value,
		services.REFRESH_TOKEN_SECRET_KEY,
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": "Could not decode refresh token",
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
	var newUserEmail schema.Email
	if err := c.BindJSON(&newUserEmail); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Bad request",
		})
		return
	}

	if !services.VerifyEmailFormat(newUserEmail.Email) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail":"Email address not valid",
		})
		return
	}

	emailInDb, _ := services.IsEmailInDB(
		newUserEmail.Email,
	)

	if emailInDb {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail":"Email already taken",
		})
	}

	services.SendVerifyMessageToEmailClient(
		newUserEmail.Email,
	)
	c.Status(http.StatusAccepted)
}
