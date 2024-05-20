package orm

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InsertInto(models ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			// panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	for _, model := range models {
		t := reflect.TypeOf(model)
		v := reflect.ValueOf(model)
		tableName := t.Name()
		fields := []string{}
		placeholders := []string{}
		values := []interface{}{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fields = append(fields, field.Name)
			placeholders = append(placeholders, "?")
			values = append(values, v.Field(i).Interface())
		}
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, strings.Join(fields, ", "), strings.Join(placeholders, ", "))
		_, err = tx.Exec(query, values...)
		if err != nil {
			return err
		}
	}

	return nil
}
