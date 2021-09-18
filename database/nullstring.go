package database

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func NewNullString(s string) *NullString {
	ns := sql.NullString{
		String: s,
		Valid:  true,
	}

	if s == "" {
		ns.Valid = false
	}

	return &NullString{ns}
}

func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) ToString() string {
	return v.NullString.String
}
