package sorm

import (
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func (db Database) Delete(model interface{}, condition string, args ...interface{}) error {
	t := reflect.TypeOf(model)
	tableName := t.Name()
	query := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, condition)
	fmt.Println(query)
	_, err := db.Connection.Exec(query, args...)
	return err
}
