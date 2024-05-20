package orm

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InsertInto(model interface{}) error {
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
	_, err := db.Exec(query, values...)
	return err
}
