package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringList []string

func (l StringList) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}

	encoded, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}

func (l *StringList) Scan(src interface{}) error {
	if src == nil {
		*l = nil
		return nil
	}

	var raw []byte
	switch value := src.(type) {
	case []byte:
		raw = value
	case string:
		raw = []byte(value)
	default:
		return fmt.Errorf("unsupported type for StringList: %T", src)
	}

	if len(raw) == 0 {
		*l = nil
		return nil
	}

	return json.Unmarshal(raw, l)
}
