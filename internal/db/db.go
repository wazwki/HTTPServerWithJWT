package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const connStr string = "user=user dbname=dbname sslmode=disable password=1234"

func Connect() error {
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXIST users(
		username TEXT,
		email TEXT
	)`)
	if err != nil {
		return err
	}

	return nil
}
