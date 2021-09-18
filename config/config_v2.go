package config

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func Create(c interface{}, tag string) {
	CreateWithReplacer(c, tag, NoneReplacer())
}

func CreateWithMark(c interface{}, tag string, secretMap map[string]string) {
	CreateWithReplacer(c, tag, SecretMapReplacer(secretMap))
}

func CreateWithReplacer(c interface{}, tag string, replacer Replacer) {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		panic("need pointer")
	}
	v := reflect.ValueOf(c).Elem()
	setStructValue(v, tag, replacer)
}

type Replacer func(content string) string

func NoneReplacer() Replacer {
	return func(src string) string {
		return src
	}
}

func SecretMapReplacer(secret map[string]string, surrend ...string) Replacer {
	surrendLeft := "{{"
	surrendRight := "}}"
	if len(surrend) >= 2 {
		surrendLeft = surrend[0]
		surrendRight = surrend[1]
	}
	var pattern = regexp.MustCompile(surrendLeft + "(.*?)" + surrendRight)

	return func(src string) string {
		dst := pattern.ReplaceAllFunc([]byte(src), func(matched []byte) []byte {
			val, ok := secret[string(matched[2:len(matched)-2])]
			if ok {
				return []byte(val)
			}
			return matched
		})
		return string(dst)
	}
}

func setStructValue(v reflect.Value, tag string, replacer Replacer) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tagVal := t.Field(i).Tag.Get(tag)
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			setStructValue(v.Field(i), tag, replacer)
		case reflect.String:
			replacedVal := replacer(tagVal)
			v.Field(i).SetString(replacedVal)
		case reflect.Int, reflect.Int64:
			intVal, _ := strconv.ParseInt(tagVal, 10, 64)
			v.Field(i).SetInt(intVal)
		case reflect.Bool:
			boolVal, _ := strconv.ParseBool(tagVal)
			v.Field(i).SetBool(boolVal)
		case reflect.Slice:
			sliceVal := strings.Split(tagVal, ",")
			switch v.Field(i).Interface().(type) {
			case []string:
				v.Field(i).Set(reflect.ValueOf(sliceVal))
			case []int:
				intVal := make([]int, 0, len(sliceVal))
				for _, sval := range sliceVal {
					iv, _ := strconv.Atoi(sval)
					intVal = append(intVal, iv)
				}
				v.Field(i).Set(reflect.ValueOf(intVal))
			}
		default:
			panic("can't parse config: field name:" + t.Field(i).Name + ", type:" + t.Field(i).Type.Kind().String())
		}
	}
}
