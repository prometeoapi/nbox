package models

import (
	"fmt"
)

type Entry struct {
	Path string `json:"-"`
	Key  string `json:"key"`
	//Value json.RawMessage `json:"value"`
	Value string `json:"value"`
}

func (e *Entry) String() string {
	return fmt.Sprintf("Key: %s. Value: %s", e.Key, e.Value)
}
