package sorm

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func (db Database) InsertInto(models ...interface{}) error {
  // Starts database transaction
	tx, err := db.Connection.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback error: %v, original panic: %v", rbErr, p)
			} else {
				panic(p) // re-throw panic after Rollback
			}
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback error: %v, original error: %v", rbErr, err)
			}
		} else {
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("commit error: %v", err)
			}
		}
	}()

	// insert
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
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, strings.Join(fields, ", "), strings.Join(placeholders, ", ")) // #nosec G201
		stmt, err := tx.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return err
}
