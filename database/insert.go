package database

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
)

func buildQuestionMark(n int) string {
	s := strings.Repeat("?,", n)
	return s[:len(s)-1]
}

func BatchInsertionWithSliceMap(table string, sliceMap interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(sliceMap)
	var isArray bool
	var values []reflect.Value
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		isArray = true
		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Index(i))
		}
	}

	if !isArray {
		return "", nil, fmt.Errorf("%+v is not iterable", sliceMap)
	}

	if len(values) < 1 {
		return "", nil, fmt.Errorf("%+v len less than 1", sliceMap)
	}

	mapKeys := []reflect.Value{}
	columns := []string{}
	switch values[0].Kind() {
	case reflect.Map:
		for _, key := range values[0].MapKeys() {
			columns = append(columns, key.String())
			mapKeys = append(mapKeys, key)
		}
	}

	columnString := "`" + strings.Join(columns, "`,`") + "`"

	sqlStr := "insert into " + table + "(" + columnString + ") values"
	sqlBuf := bytes.NewBufferString(sqlStr)
	args := []interface{}{}
	for i, v := range values {
		switch v.Kind() {
		case reflect.Map:
			n := 0
			for _, key := range mapKeys {
				args = append(args, v.MapIndex(key).Interface())
				n++
			}

			sqlBuf.WriteString("(" + buildQuestionMark(n) + ")")
		}

		if i+1 != len(values) {
			sqlBuf.WriteString(",")
		}
	}
	return sqlBuf.String(), args, nil
}

// slices only suport []struct and [][]slice
func BatchInsertion(table string, columns string, slices interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(slices)
	var isArray bool
	var values []reflect.Value
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		isArray = true
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Kind() == reflect.Ptr {
				values = append(values, v.Index(i).Elem())
			} else {
				values = append(values, v.Index(i))
			}
		}
	}

	if !isArray {
		return "", nil, fmt.Errorf("%+v is not iterable", slices)
	}

	if len(values) < 1 {
		return "", nil, fmt.Errorf("%+v len less than 1", slices)
	}

	sqlStr := "insert into " + table + "(" + columns + ") values"
	sqlBuf := bytes.NewBufferString(sqlStr)
	args := []interface{}{}
	for i, v := range values {
		switch v.Kind() {
		case reflect.Struct:
			n := 0
			for i := 0; i < v.NumField(); i++ {
				args = append(args, v.Field(i).Interface())
				n++
			}
			sqlBuf.WriteString("(" + buildQuestionMark(n) + ")")

		case reflect.Slice, reflect.Array:
			n := 0
			for i := 0; i < v.Len(); i++ {
				args = append(args, v.Index(i).Interface())
				n++
			}
			sqlBuf.WriteString("(" + buildQuestionMark(n) + ")")

		default:
			log.Printf("%+v", v.Kind())
		}

		if i+1 != len(values) {
			sqlBuf.WriteString(",")
		}
	}
	return sqlBuf.String(), args, nil
}
