package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/*
 * Creates access token for user authorization
 * User authorization only requires acess token sent as a request header,
 *		with this covering CSRF attacks
 * Access token contains the user Id, username and whether the user is an admin
 * The expiration time of the the access token is given by ACCESS_TOKEN_EXPIRE_TIME,
 *		and is signed by ACCESS_TOKEN_SECRET_KEY (see config.go)
 */
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

/*
 * Creates refresh token used for requesting a new access token
 * Refresh tokens are sent via web cookies and only contain the user Id
 * The refresh tokens do eventually expire, prompting the user to log in again,
 *		but when a new access token is provisioned, a new refresh token is provided
 */
func createRefreshToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp" : time.Now().Add(REFRESH_TOKEN_EXPIRE_TIME).Unix(),
	}
	return jwtEncode(claims, REFRESH_TOKEN_SECRET_KEY)
}

// Deletes the authorization cookies (but not the accept cookie)
func deleteCookies(c *gin.Context) {
	c.SetCookie(REFRESH_TOKEN_NAME, "", -1, "", FRONTEND_DOMAIN, false, false)
}

// Deconstructs the authorization header to get the bearer token
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

// Returns whether the access token from a request is
// func hasAccessTokenExpired(c *gin.Context) bool {
// 	authToken, err := getAuthorizationToken(c)
// 	if err != nil {
// 		return true
// 	}
// 	print("Access token: ", authToken)
// 	claims, err := jwtDecodeClaims(authToken, ACCESS_TOKEN_SECRET_KEY)
// 	if err != nil || claims["exp"] == nil {
// 		return true
// 	}

// 	expireTime, ok := claims["exp"].(int64)
// 	if !ok {
// 		return true
// 	}

// 	if expireTime > time.Now().Unix() {
// 		return true
// 	}
// 	return false
// }

/*
 * Creates a refresh token and sets it as a cookie
 * The router method aborts if an error is thrown
 */
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

// Encodes and signs jwt claims into a token string
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

func jwtDecodeClaims(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface {}, error) {
			return secretKey, nil
		},
	)

	if err != nil {
		return claims, err
	} else if !token.Valid {
		return claims, errors.New("parsed JWT token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims, errors.New("error getting JWT token claims")
	}
	return claims, nil
}