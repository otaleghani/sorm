package sorm

import (
	"os"
  "fmt"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//var db *sql.DB

type Database struct {
  Path string
  Connection *sql.DB
}

func CreateDatabase(dbPath string, foreignKeys bool) (*Database, error) {
  // logNotice(fmt.Sprintf("Opening database connection at %v", dbPath))
  db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
    logError(fmt.Sprintf("sql:Open, %v", err))
		return &Database{}, err
	}
  logInfo(fmt.Sprintf("sql.Open at path %v", dbPath))

  if foreignKeys {
    _, err = db.Exec("PRAGMA foreign_keys = ON;")
    if err != nil {
      logError(fmt.Sprintf("db.Exec(\"PRAGMA foreign_keys = ON;\"), %v", err))
      return &Database{}, err
    }
    logInfo("foreign keys activated")
  }

  // logSuccess("Database opened")
  return &Database{Connection: db, Path: dbPath}, db.Ping()
}

func DeleteDatabase(dbPath string) error {
	return os.Remove(dbPath)
}

