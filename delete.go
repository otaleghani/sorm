package sorm

import (
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) Delete(model interface{}, condition string, args ...interface{}) error {
	t := reflect.TypeOf(model)
	tableName := t.Name()

  logNotice(fmt.Sprintf("Deleting from table %v with condition %v", tableName, condition))
  logInfo(fmt.Sprintf("Executing query on database %v", db.Path))
	
  query := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, condition)
	result, err := db.Connection.Exec(query, args...)
  if err != nil {
    logError(fmt.Sprintf("db.Connection.Exec: %v", err))
    return err
  }

  rowsAffected, err := result.RowsAffected()
	if err != nil {
    logError(fmt.Sprintf("result.RowsAffected: %v", err))
		return err
	}	
  if rowsAffected > 0 {
    logSuccess("Item deleted")
    return nil
	} else {
    logInfo("Item was not found, so it was not deleted")
    return nil
	}

	return nil
}
