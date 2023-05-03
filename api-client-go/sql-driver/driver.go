package driver

import (
	"fmt"
	"log"
	"database/sql"

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

func (d *SqlDriver) Query(queryStatement string, args ...any) *sql.Rows {
	rows, err := d.DB.Query(queryStatement, args...)
	if err != nil {
		log.Fatal(err)
	}	
	return rows
}

func (d *SqlDriver) Exec(statement string, args ...any) (sql.Result, error) {
	return d.DB.Exec(statement, args...)
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