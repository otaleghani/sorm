package sorm

import (
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) Delete(model interface{}, condition string, args ...interface{}) error {
	t := reflect.TypeOf(model)
	tableName := t.Name()
	
  query := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, condition)
	result, err := db.Connection.Exec(query, args...)
  if err != nil {
    logError(fmt.Sprintf("%v 400, db.Connection.Exec(), %v", query, err))
    return err
  }

  rowsAffected, err := result.RowsAffected()
	if err != nil {
    logError(fmt.Sprintf("%v 400, result.RowsAffected, %v", query, err))
		return err
	}	
  if rowsAffected > 0 {
    logInfo(fmt.Sprintf("%v item deleted", query))
    return nil
	} else {
    logInfo(fmt.Sprintf("%v, item not found", query))
    return nil
	}

	return nil
}
