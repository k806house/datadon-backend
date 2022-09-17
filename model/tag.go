package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type TagList []Tag

func (a TagList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *TagList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
