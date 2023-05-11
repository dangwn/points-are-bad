package driver

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

type SqlDriver struct {
	DB *sql.DB
}

func NewSqlDriver(
	user string,
	password string,
	db string,
) (*SqlDriver) {
	database, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",user, password, db),
	)
	if err != nil{
		log.Fatal(err)
	}
	return &SqlDriver{
		DB: database,
	}
}

func (d *SqlDriver) Delete(table string, whereConditions string, args ...any) (bool, error) {
	whereClause := ""
	if whereConditions != "" {
		if strings.ToLower(whereConditions[:6]) == "where " {
			whereClause += whereConditions
		} else {
			whereClause += "WHERE " + whereConditions
		}
	}

	res, err := d.Exec("DELETE FROM "+table+" "+whereClause, args...)
	if err == nil {
		if _, err := res.RowsAffected(); err == nil {
			return true, err
		} else {
			return false, err
		}
	}
	return false, err
}

func (d *SqlDriver) Exec(statement string, args ...any) (sql.Result, error) {
	return d.DB.Exec(statement, args...)
}

func (d *SqlDriver) Insert(
	table string,
	columnsString string,
	valuesString string,
	args ...any,
) (sql.Result, error) {
	result, err := d.Exec(
		"INSERT INTO " + table + "(" + columnsString + ") VALUES (" + valuesString + ")",
		args...,
	)
	if err != nil {
		return result, err
	}
	if _, err := result.RowsAffected(); err != nil {
		return result, err
	}
	return result, nil
}

func (d *SqlDriver) InsertWithReturn(
	table string,
	columnsString string,
	valuesString string,
	returnString string,
	args ...any,
) *sql.Row {
	return d.DB.QueryRow(
		"INSERT INTO " + table + "(" + columnsString + ") VALUES (" + valuesString + ") RETURNING " + returnString,
		args...,
	)
}

func (d *SqlDriver) Query(queryStatement string, args ...any) (*sql.Rows, error) {
	return d.DB.Query(queryStatement, args...)
}

func (d *SqlDriver) QueryRow(queryStatement string, args ...any) *sql.Row {
	return d.DB.QueryRow(queryStatement, args...)
}

func (d *SqlDriver) ValueExists(table string, column string, value any) (bool, error) {
	var exists bool 
	
	err := d.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + column + " = $1)",
		value,
	).Scan(&exists)
	
	if err != nil {
		return false, err
	}
	
	return exists, nil
}