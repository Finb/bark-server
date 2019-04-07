package main

import (
	"encoding/json"
	"os"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestMessages(t *testing.T) {
	const dbfile = "/tmp/bark.db"
	defer os.Remove(dbfile)

	conn, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	db := NewDB(conn)

	type Message struct {
		Topic string `json:"topic"`
	}
	t.Run("Append", func(t *testing.T) {
		message := &Message{}
		message.Topic = "a-value"
		key := "mock_key_append"
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Append(key, *message); err != nil {
			t.Fatal(err)
		}

		data := db.Get(bkt, key)
		if data == nil {
			t.Fatal("should not get nil for messages")
		}

		message.Topic = "b-value"
		if err := db.Append(key, *message); err != nil {
			t.Fatal(err)
		}
		data = db.Get(bkt, key)

		var ms []Message
		err = json.Unmarshal(data, &ms)
		if err != nil {
			t.Fatal(err)
		}
		if len(ms) != 2 {
			t.Fatal("messages number should be 2")
		}
		if ms[0].Topic != "a-value" || ms[1].Topic != "b-value" {
			t.Fatal("append not working")
		}
	})
}
