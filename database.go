package sorm

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func CreateDatabase(dbPath string) (err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	return db.Ping()
}

func DeleteDatabase(dbPath string) error {
	return os.Remove(dbPath)
}
