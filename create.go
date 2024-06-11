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

func (db *Database) CreateTable(model interface{}) error {
	t := reflect.TypeOf(model)
	tableName := t.Name()
	fields := []string{}
  constraints := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		sqlType := sqlType(field.Type)

    if field.Name == "Id" {
			fields = append(fields, fmt.Sprintf("%s %s PRIMARY KEY", field.Name, sqlType))
      continue
    }

    parts := strings.Split(field.Name, "_")
    if len(parts) > 1 {
      suffix := parts[len(parts)-1]
      switch suffix {
      case "u":
        fields = append(fields, fmt.Sprintf("%s %s UNIQUE", field.Name, sqlType))
      case "n":
        fields = append(fields, fmt.Sprintf("%s %s NOT NULL", field.Name, sqlType))
      case "nu":
        fields = append(fields, fmt.Sprintf("%s %s NOT NULL UNIQUE", field.Name, sqlType))
      case "id":
			  fields = append(fields, fmt.Sprintf("%s %s", field.Name, sqlType))
        constraints = append(constraints, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(Id) ON UPDATE CASCADE", field.Name, parts[0]))
      default: 
			  fields = append(fields, fmt.Sprintf("%s %s", field.Name, sqlType))
      }
    } else {
			fields = append(fields, fmt.Sprintf("%s %s", field.Name, sqlType))
    }
	}
  query := ""
  if len(constraints) == 0 {
	  query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(fields, ", "))
  } else {
	  query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, %s);", tableName, strings.Join(fields, ", "), strings.Join(constraints, ", "))
  }
  fmt.Println(query)
	_, err := db.Connection.Exec(query)
	return err
}
