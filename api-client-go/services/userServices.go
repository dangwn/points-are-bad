package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"

	"points-are-bad/api-client/schema"
)

func createUserId() string {
	return uuid.New().String()
}

func CreateNewUser(username string, email string, password string) (string, error) {
	// If the new user is the first in the db, they are an admin
	adminInDB, err := driver.ValueExists("users", "is_admin", true)
	if err != nil {
		log.Println(err)
		return "", err
	}
	isAdmin := !adminInDB

	if _, err = verifyEmailIsUnique(email); err != nil {
		log.Println(err)
		return "", err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	userId := createUserId()

	if insertUserIntoDb(
		userId, username, email, hashedPassword, isAdmin,
	); err != nil {
		log.Println(err)
		return "", err
	}

	if err = insertPointsIntoDb(userId); err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("New user created")
	return userId, nil
}

func decodeVerificationToken(token string) (string, string, error) {
	// Returns email, token, error
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Println("Could not decode token")
		return "", "", nil
	}

	var data map[string]string
	err = json.Unmarshal(decodedToken, &data)
	if err != nil {
		log.Println("Could not unmarshal JSON")
		return "", "", nil
	}

	for key := range data {
		return key, data[key], nil
	}
	return "", "", errors.New("data was empty")
}

func DeleteUserById(userId string) error {
	result, err := driver.Exec(
		"DELETE FROM users WHERE user_id = $1",
		userId,
	)	
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 
	log.Println("User deleted")
	return nil
}

func GetUserById(userId string) (schema.User, error) {
	var user schema.User
	rows, err := driver.Query(
		"SELECT * FROM users WHERE users.user_id = $1 LIMIT 1",
		userId,
	)
	if err != nil {
		return user, err
	}
	
	for rows.Next() {
		if err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.Email,
			&user.HashedPassword,
			&user.IsAdmin,
		); err != nil {
			return user, err
		}
		break
	}
	return user, nil
}

func getUserPasswordHash(userId string) (string, error) {
	var hash string
	rows, err := driver.Query(
		"SELECT hashed_password FROM users WHERE user_id = $1 LIMIT 1",
		userId,
	)
	if err != nil {
		return hash, err
	}
	
	for rows.Next() {
		if err := rows.Scan(
			&hash,
		); err != nil {
			return hash, err
		}
		break
	}
	return hash, nil
}

func insertUserIntoDb(
	userId string,
	username string,
	email string, 
	hashedPassword string, 
	isAdmin bool,
) error {
	_, err := driver.Insert(
		"users",
		"user_id, username, email, hashed_password, is_admin",
		"$1, $2, $3, $4, $5",
		userId,
		username, 
		email, 
		hashedPassword, 
		isAdmin,
	)
	return err
}

func UpdateUsernameByUserId(userId string, username string) error {
	result, err := driver.Exec(
		"UPDATE users SET username = $1 WHERE user_id = $2",
		username,
		userId,
	)
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 

	log.Println("Username updated")
	return nil
}

func UpdatePasswordByUserId(userId string, oldPassword string,	newPassword string) error {
	oldPasswordHash, err := getUserPasswordHash(userId)
	if err != nil{		
		log.Println(err)
		return err
	}

	if !verifyPasswordHash(oldPasswordHash, oldPassword) {
		err := errors.New("old password was not correct")
		log.Println(err)
		return err
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := driver.Exec(
		"UPDATE users SET hashed_password = $1 WHERE user_id = $2",
		hashedPassword,
		userId,
	)
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 

	log.Println("Password changed")
	return nil
}

func ValidateVerificationToken(token string) (string, error) {
	// Email comes back from decode function in the form verify--<email>
	// Only return the email for creating new user
	email, decodedToken, err := decodeVerificationToken(token)
	if err != nil {
		return "", err
	}

	tokenInRedis, err := redis.Get(redisContext, email).Result()
	if err != nil {
		return "", err
	}

	if decodedToken != tokenInRedis {
		return "", errors.New("token in redis did not match verified token")
	}

	return email[len("verify--"):], nil
}

func verifyEmailIsUnique(email string) (bool, error) {
	emailExists, err := driver.ValueExists("users", "email", email)
	if err != nil {
		return false, err
	}
	if emailExists {
		err = errors.New("email already exists in db")
		return false, err
	}
	return true, nil
}