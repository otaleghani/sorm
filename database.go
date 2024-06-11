package sorm

import (
	"os"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//var db *sql.DB

type Database struct {
  Path string
  Connection *sql.DB
}

func CreateDatabase(dbPath string, foreignKeys bool) (*Database, error) {
  db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &Database{}, err
	}
  if foreignKeys {
    _, err = db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {
      return &Database{}, err
    }
  }
  return &Database{Connection: db, Path: dbPath}, db.Ping()
}

func DeleteDatabase(dbPath string) error {
	return os.Remove(dbPath)
}
