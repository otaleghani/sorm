package sorm

import (
	"fmt"
	"reflect"
	"strings"
)

func (db *Database) Update(model interface{}, conditions string, args ...interface{}) error {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)

	tableName := t.Name()
	fields := []string{}
	values := []interface{}{}
  logNotice(fmt.Sprintf("Updating in table %v", tableName))

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if value != 0 {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value != 0 {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		case reflect.Float32, reflect.Float64:
			if value != 0.0 {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		case reflect.Bool:
			if value != 0 {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		case reflect.String:
			if value != "" {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		default:
			if value != "" {
				fields = append(fields, fmt.Sprintf("%s = ?", field.Name))
				values = append(values, value)
			}
		}
	}

	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(fields, ", "))
	if conditions != "" {
		query = fmt.Sprintf("%s WHERE %s", query, conditions)
	}

	// Append condition arguments to the values slice
	values = append(values, args...)

	result, err := db.Connection.Exec(query, values...)
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
    logSuccess(fmt.Sprintf("Update successful. Rows affected: %v", rowsAffected))
    return nil
	} else {
    logInfo("No rows were affected")
    return nil
	}

	return nil
}
