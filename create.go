package sorm

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func sqlType(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.String:
		return "TEXT"
	default:
		return "TEXT" // Default case, may need more specific handling
	}
}

func (db Database) CreateTable(model interface{}) error {
	t := reflect.TypeOf(model)
	tableName := t.Name()
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sqlType := sqlType(field.Type)
		if i == 0 {
			fields = append(fields, fmt.Sprintf("%s %s PRIMARY KEY", field.Name, sqlType))
		} else {
			fields = append(fields, fmt.Sprintf("%s %s", field.Name, sqlType))
		}
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(fields, ", "))
	_, err := db.Connection.Exec(query)
	return err
}
