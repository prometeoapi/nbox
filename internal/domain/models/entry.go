package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Entry struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func (e *Entry) String() string {
	return fmt.Sprintf("Key: %s. Value: %s", e.Key, hex.EncodeToString(e.Value))
}
