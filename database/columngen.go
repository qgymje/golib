package database

import (
	"bytes"
	"reflect"
)

func GenerateColumnWithPrefix(s interface{}, tag string, prefix string, exclude ...string) string {
	if prefix != "" {
		prefix = prefix + "."
	}
	cf := func(c string) string {
		return prefix + c
	}
	return generateColumn(cf, s, tag, false, exclude...)
}

func GenerateColumn(s interface{}, tag string, exclude ...string) string {
	cf := func(c string) string {
		return c
	}
	return generateColumn(cf, s, tag, true, exclude...)
}

type columnFunc func(string) string

func generateColumn(cf columnFunc, s interface{}, tag string, useBackQuote bool, exclude ...string) string {
	b := bytes.NewBufferString(" ")
	t := reflect.TypeOf(s)
	m := toMap(exclude)

	var sl int
	switch t.Kind() {
	case reflect.Struct:
		sl = t.NumField()
		for i := 0; i < sl; i++ {
			column := t.Field(i).Tag.Get(tag)
			if !isColumnValid(column) {
				continue
			}

			if _, ok := m[column]; ok {
				continue
			}
			if useBackQuote {
				b.WriteString("`" + cf(column) + "`,")
			} else {
				b.WriteString(cf(column) + ",")
			}
		}
	}
	return b.String()[:b.Len()-1] + " "
}

func isColumnValid(column string) bool {
	return column != "-" && column != ""
}

func toMap(ss []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	return m
}
