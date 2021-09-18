package database

import "strings"

func AddSpaceAround(s string) string {
	return " " + strings.Trim(s, " ") + " "
}
