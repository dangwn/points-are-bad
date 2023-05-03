package services

import (
	"log"
	"errors"
	
	"points-areb-bad/api-client/schema"
)

func GetAllUsers() []schema.User{
	rows := driver.Query("SELECT * FROM users;")
	users := []schema.User{}
	for rows.Next() {
		var user schema.User
		if err := rows.Scan(
			&user.UserId,
			&user.Username,
			&user.Email,
			&user.HashedPassword,
			&user.IsAdmin,
		); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users
}

func GetUserById(userId int) (schema.User, error) {
	rows := driver.Query(
		"SELECT * FROM users WHERE users.user_id = $1 LIMIT 1",
		userId,
	)
	var user schema.User
	
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

func CreateNewUser(
	username string,
	email string,
	password string,
) error {
	// If the new user is the first in the db, they are an admin
	adminInDB, err := driver.ValueExists("users", "is_admin", true)
	if err != nil {
		log.Println(err)
		return err
	}
	isAdmin := !adminInDB

	emailExists, err := driver.ValueExists("users", "email", email)
	if err != nil {
		log.Println(err)
		return err
	}
	if emailExists {
		err = errors.New("Email already exists in db")
		log.Println(err)
		return err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Println(err)
		return err
	}

	result, err := driver.Exec(
		"INSERT INTO users (username, email, hashed_password, is_admin) VALUES ($1, $2, $3, $4);",
		username, 
		email, 
		hashedPassword, 
		isAdmin,
	)	
	if err != nil {
		log.Println(err)
		return err
	} 

	if _, err:= result.RowsAffected(); err != nil {
		log.Println(err)
		return err
	} 
	log.Println("New user created")
	return nil
}

func DeleteUserById(
	userId int,
) error {
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

func UpdateUsernameByUserId(
	userId int,
	username string,
) error {
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

func UpdatePasswordByUserId(
	userId int,
	oldPassword string,
	newPassword string,
) error {
	oldPasswordHash, err := getUserPasswordHash(userId)
	if err != nil{		
		log.Println(err)
		return err
	}

	if !verifyPasswordHash(oldPasswordHash, oldPassword) {
		err := errors.New("Old password was not correct")
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

func getUserPasswordHash(
	userId int,
) (string, error) {
	rows := driver.Query(
		"SELECT hashed_password FROM users WHERE user_id = $1 LIMIT 1",
		userId,
	)
	var hash string
	
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