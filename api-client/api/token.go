package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func createAccessToken(userId string, username string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"username": username,
		"is_admin": isAdmin,
		"iat": time.Now().Unix(),
		"exp" : time.Now().Add(ACCESS_TOKEN_EXPIRE_TIME).Unix(),
	}
	return jwtEncode(claims, ACCESS_TOKEN_SECRET_KEY)
}

func createRefreshToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp" : time.Now().Add(REFRESH_TOKEN_EXPIRE_TIME).Unix(),
	}
	return jwtEncode(claims, REFRESH_TOKEN_SECRET_KEY)
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
	c.SetCookie(REFRESH_TOKEN_NAME, refreshToken, REFRESH_TOKEN_EXPIRE_TIME_SECONDS, "", FRONTEND_DOMAIN, false, false)
	return nil
}

func jwtEncode(claims jwt.MapClaims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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