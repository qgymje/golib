package database

import (
	"database/sql/driver"
	"encoding/json"
)

type JsonColumn struct {
	value interface{}
	valid bool
}

func NewJsonColumn(v interface{}) *JsonColumn {
	return &JsonColumn{
		value: v,
		valid: true,
	}
}

func NewEmptyJsonColumn() *JsonColumn {
	return &JsonColumn{valid: false}
}

func NewJsonColumnWithMapString(v string) (*JsonColumn, error) {
	value := make(map[string]interface{})
	if err := json.Unmarshal([]byte(v), &value); err != nil {
		return nil, err
	}
	return &JsonColumn{
		valid: true,
		value: value,
	}, nil
}

func NewJsonColumnWithSliceString(v string) (*JsonColumn, error) {
	value := []map[string]interface{}{}
	if err := json.Unmarshal([]byte(v), &value); err != nil {
		return nil, err
	}
	return &JsonColumn{
		valid: true,
		value: value,
	}, nil
}

func (j *JsonColumn) MarshalJSON() ([]byte, error) {
	if j.valid {
		return json.Marshal(&j.value)
	}
	return []byte(`""`), nil
}

func (j *JsonColumn) UnmarshalJSON(data []byte) error {
	if j.valid {
		return json.Unmarshal(data, &j.value)
	}
	return nil
}

func (m *JsonColumn) Scan(value interface{}) error {
	if value == nil {
		m.value = nil
		m.valid = false
		return nil
	}
	m.value = value
	m.valid = true
	return json.Unmarshal([]byte(value.([]uint8)), m)
}

func (m JsonColumn) Value() (driver.Value, error) {
	if m.valid {
		return json.Marshal(&m)
	}
	return "", nil
}
