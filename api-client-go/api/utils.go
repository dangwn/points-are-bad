package api

import (
	"log"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func createDateRangeWhereClause(startDate *Date, endDate *Date) string {
	var whereClause string

	if startDate != nil {
		whereClause += " WHERE match_date >= '" + startDate.String() + "'"
	}
	if endDate != nil {
		if whereClause != "" {
			return whereClause + " AND match_date < '" + endDate.String() + "'"
		} else {
			return " WHERE match_date < '" + endDate.String() + "'"
		}
	}
	
	return whereClause
}

func abortRouterMethod(c *gin.Context, statusCode int, msg interface{}, logs ...any) {
	for _, l := range logs {
		log.Println(l)
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"detail": msg})
}


func redisKeyExists(key string) (bool, error) {
	exists, err := redis.Exists(redisContext, key).Result()
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
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

func verifyEmailFormat(email string) bool {
	emailReg := regexp.MustCompile(`[^@]+@[^@]+\.[^@]+`)
	return emailReg.MatchString(email)
}