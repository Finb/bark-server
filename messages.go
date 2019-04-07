package main

import (
	"encoding/json"
)

const bkt = "messages"

// Append data
func (db *DB) Append(key string, v interface{}) error {
	exists := db.Get(bkt, key)
	var messages []interface{}
	if exists != nil {
		if err := json.Unmarshal(exists, &messages); err != nil {
			return err
		}
	}
	list := append(messages, v)
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return db.Set(bkt, key, data)
}
