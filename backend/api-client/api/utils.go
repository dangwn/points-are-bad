package api

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

/*
 * Used by router handler functions to abort with a given status code and message
 * Optional error log messages can be sent to the console
 */
func abortRouterMethod(c *gin.Context, statusCode int, msg interface{}, logs ...any) {
	for _, l := range logs {
		Logger.Error(l)
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"detail": msg})
}

/*
 * Creates a where clause for sql queries between zero, one or two dates
 * Start date inclusive, end date exclusive
 */
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

// Creates a hash for a given password
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

// Verifies a given password a password hash
func verifyPasswordHash(hashedPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	); err != nil{
		return false
	}
	return true
}

// Verifies a string is of email format
func verifyEmailFormat(email string) bool {
	emailReg := regexp.MustCompile(`[^@]+@[^@]+\.[^@]+`)
	return emailReg.MatchString(email)
}