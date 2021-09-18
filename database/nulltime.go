package database

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullTime struct {
	Time   time.Time
	Valid  bool
	format string
}

func NewNullTime(t time.Time) NullTime {
	nt := NullTime{}
	if !t.IsZero() {
		nt.Valid = true
	}
	nt.Time = t
	return nt
}

func NewNullTimeWithFormat(t time.Time, format string) NullTime {
	nt := NewNullTime(t)
	nt.format = format
	return nt
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) SetFormat(format string) {
	nt.format = format
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if nt.format != "" && nt.Valid {
		t := nt.Time.Format(nt.format)
		buf := bytes.NewBufferString("\"")
		buf.WriteString(t)
		buf.WriteString("\"")
		return buf.Bytes(), nil
	}
	return json.Marshal(&nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &nt.Time)
}
