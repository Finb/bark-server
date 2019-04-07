package main

import (
	"encoding/json"
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestDB(t *testing.T) {
	const dbfile = "/tmp/bark.db"
	defer os.Remove(dbfile)

	conn, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		t.Fatal(err)
	}

	type Message struct {
		From string `json:"from"`
	}
	db := NewDB(conn)
	bkt := "bkt"
	key := "mock_key"

	t.Run("Set/Get", func(t *testing.T) {
		message := Message{
			From: "a@a.com",
		}
		msg, err := json.Marshal(message)
		if err != nil {
			t.Fatal(err)
		}

		if err := db.Set(bkt, key, msg); err != nil {
			t.Fatal(err)
		}

		data := db.Get(bkt, key)
		if data == nil {
			t.Fatal("should not get nil for messages")
		}

		var m Message
		err = json.Unmarshal(data, &m)
		if err != nil {
			t.Fatal(err)
		}
		if m.From != "a@a.com" {
			t.Fatal("incorrect data")
		}
	})
}
