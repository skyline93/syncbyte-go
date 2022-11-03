package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB []interface{}

func (j *JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j)
}
