package services

import (
	"time"
	"errors"
	"regexp"
	
	"github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
)

func JwtEncode(
	subject string,
	secretKey []byte,
	expireDuration time.Duration,
) (string, error) {
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

func JwtDecode(
	tokenString string,
	secretKey []byte,
) (string, error) {
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
		return "", errors.New("Parsed JWT token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Error getting JWT token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("Error getting JWT token sub claim")
	}

	return sub, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func verifyPasswordHash(hashedPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	); err != nil{
		return false
	}
	return true
}

func VerifyEmailFormat(email string) bool {
	emailReg := regexp.MustCompile(`[^@]+@[^@]+\.[^@]+`)
	return emailReg.MatchString(email)
}