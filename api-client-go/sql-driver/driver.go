package driver

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type SqlDriver struct {
	DB *sql.DB
}

func NewSqlDriver(user string, password string, db string) *SqlDriver {
	database, err := sql.Open(
		"postgres",
		"user=" + user + " password=" + password + " dbname=" + db + " sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	return &SqlDriver{
		DB: database,
	}
}

func (d *SqlDriver) Exec(statement string, args ...any) (sql.Result, error) {
	return d.DB.Exec(statement, args...)
}

func (d *SqlDriver) Query(queryStatement string, args ...any) (*sql.Rows, error) {
	return d.DB.Query(queryStatement, args...)
}

func (d *SqlDriver) QueryRow(queryStatement string, args ...any) *sql.Row {
	return d.DB.QueryRow(queryStatement, args...)
}

func (d *SqlDriver) ValueExists(table string, column string, value any) (bool, error) {
	var exists bool 
	
	if err := d.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM " + table + " WHERE " + column + " = $1)",
		value,
	).Scan(&exists); err != nil {
		return false, err
	}
	
	return exists, nil
}