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
	fields := []string{}
  constraints := []string{}
	t := reflect.TypeOf(model)
	tableName := t.Name()
  logNotice(fmt.Sprintf("Creating table: %v", tableName))

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
    logInfo(fmt.Sprintf("\t%v", fields[i]))
	}

  query := ""
  if len(constraints) == 0 {
	  query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(fields, ", "))
  } else {
	  query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, %s);", tableName, strings.Join(fields, ", "), strings.Join(constraints, ", "))
  }

  logInfo(fmt.Sprintf("Executing query on database %v", db.Path))
	_, err := db.Connection.Exec(query)
  if err != nil {
    logError(fmt.Sprintf("db.Exec, %v", err))
    return err
  }
  logSuccess(fmt.Sprintf("Table %v created", tableName))

  logInfo(fmt.Sprintf("Inserting default nil value into %v", tableName))
  err = db.InsertNil(model)
  if err != nil {
    logError(fmt.Sprintf("db.InsertNil, %v", err))
    return err
  }
  logSuccess(fmt.Sprintf("Nil value inserted into %v", tableName))

	return nil
}

func (db *Database) InsertNil(model interface{}) error {
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

  t := reflect.TypeOf(model)
  // v := reflect.ValueOf(model)
  tableName := t.Name()
  fields := []string{}
  placeholders := []string{}
  values := []interface{}{}
  for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fields = append(fields, field.Name)
    placeholders = append(placeholders, "?")
		sqlType := sqlType(field.Type)

    switch sqlType {
    case "INTEGER":
      values = append(values, 0)
    case "REAL":
      values = append(values, 0)
    case "BOOLEAN":
      values = append(values, false)
    case "TEXT":
      values = append(values, "nil")
    }
  }

  query := fmt.Sprintf(
    "INSERT INTO %s (%s) VALUES (%s);", 
    tableName,
    strings.Join(fields, ", "),
    strings.Join(placeholders, ", "),
  )

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

  return nil
}
