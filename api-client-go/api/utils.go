package api

import (
	"errors"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func createDateRangeWhereClause(startDate *Date, endDate *Date) string {
	var whereClause string

	if startDate != nil {
		whereClause += " WHERE match_date >= " + startDate.String()
	}
	if endDate != nil {
		if whereClause != "" {
			return whereClause + " AND match_date < " + endDate.String()
		} else {
			return " WHERE match_date < " + endDate.String()
		}
	}
	return whereClause
}

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

func validateDateString(dateString string) bool {
	if _, err := time.Parse("2010-01-01", dateString); err != nil {
		return false
	}
	return true
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