package sorm

import (
	"fmt"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func (db *Database) Select(dest interface{}, conditions string, args ...interface{}) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	sliceValue := rv.Elem()
	elemType := sliceValue.Type().Elem()
	tableName := elemType.Name()
	fields := []string{}
  logNotice(fmt.Sprintf("Selecting from %v with condition %v", tableName, conditions))

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		fields = append(fields, field.Name)
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), tableName)
	if conditions != "" {
		query = fmt.Sprintf("%s WHERE %s", query, conditions)
	}

	stmt, err := db.Connection.Prepare(query)
	if err != nil {
    logError(fmt.Sprintf("db.Connection.Prepare: %v", err))
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
    logError(fmt.Sprintf("stmt.Query: %v", err))
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
    logError(fmt.Sprintf("rows.Columns: %v", err))
		return err
	}

	for rows.Next() {
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()
		fieldPtrs := make([]interface{}, len(columns))
		for i := 0; i < len(columns); i++ {
			fieldPtrs[i] = elem.Field(i).Addr().Interface()
		}
		if err := rows.Scan(fieldPtrs...); err != nil {
    logError(fmt.Sprintf("rows.Scan: %v", err))
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, elem))
	}

  logSuccess("Items selected and returned")
	return rows.Err()
}
