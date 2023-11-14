package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"dev",       //os.Getenv("MYSQL_USER"),
		"dev",       //os.Getenv("MYSQL_PASSWORD"),
		"localhost", //os.Getenv("MYSQL_HOST"),
		"mysql",     //os.Getenv("MYSQL_PORT"),
		"dev",       //os.Getenv("MYSQL_DATABASE"),
	))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
