package models

import (
	"fmt"
	"time"
)

type Entry struct {
	Path   string `json:"path"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Secure bool   `json:"secure"`
}

func (e *Entry) String() string {
	return fmt.Sprintf("Key: %s. Value: %s", e.Key, e.Value)
}

type Tracking struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Secure    bool      `json:"secure"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

func (e *Tracking) String() string {
	return fmt.Sprintf("Key: %s. Value: %s", e.Key, e.Value)
}
