package sorm

import (
	"fmt"
	"reflect"
)

func (db *Database) Join(dest interface{}, model1 interface{}, model2 interface{}, joinCondition string) error {
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}
	sliceValue := rv.Elem()
	elemType := sliceValue.Type().Elem()

	tableName1 := reflect.TypeOf(model1).Name()
	tableName2 := reflect.TypeOf(model2).Name()
	query := fmt.Sprintf("SELECT * FROM %s JOIN %s ON %s;", tableName1, tableName2, joinCondition)
	rows, err := db.Connection.Query(query)
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
