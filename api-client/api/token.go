package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func createAccessToken(userId string) (string, error) {
	return jwtEncode(userId, ACCESS_TOKEN_SECRET_KEY, ACCESS_TOKEN_EXPIRE_TIME)
}

func createRefreshToken(userId string) (string, error) {
	return jwtEncode(userId, REFRESH_TOKEN_SECRET_KEY, REFRESH_TOKEN_EXPIRE_TIME)
}

func deleteCookies(c *gin.Context) {
	c.SetCookie(REFRESH_TOKEN_NAME, "", -1, "", FRONTEND_DOMAIN, false, false)
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

func setRefreshTokenCookie(c *gin.Context, userId string) error {
	refreshToken, err := createRefreshToken(userId)
	if err != nil {
		abortRouterMethod(c, http.StatusUnauthorized, err, err)
		deleteCookies(c)
		return err
	}
	c.SetCookie(REFRESH_TOKEN_NAME, refreshToken, REFRESH_TOKEN_EXPIRE_TIME, "", FRONTEND_DOMAIN, false, false)
	return nil
}

func jwtEncode(subject string, secretKey []byte, expireDuration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"iat": time.Now().Unix(),
		"exp" : time.Now().Add(expireDuration).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(secretKey)
}

func jwtDecode(tokenString string, secretKey []byte) (string, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)

	if err != nil{
		return "", err
	}

	if !token.Valid {
		return "", errors.New("parsed JWT token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error getting JWT token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("error getting JWT token sub claim")
	}

	return sub, nil
}