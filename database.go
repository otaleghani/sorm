package orm

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func New(dataSourceName string) (err error) {
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	return db.Ping()
}

// add delete
// add drop table?
