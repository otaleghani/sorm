package orm

import (
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

func Select(dest interface{}) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}
	sliceValue := rv.Elem()
	elemType := sliceValue.Type().Elem()

	tableName := elemType.Name()
	query := fmt.Sprintf("SELECT * FROM %s;", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
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
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, elem))
	}
	return rows.Err()
}
