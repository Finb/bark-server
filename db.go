package main

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

// DB database
type DB struct {
	client *bolt.DB
}

// NewDB create a db instance
func NewDB(conn *bolt.DB) *DB {
	return &DB{
		client: conn,
	}
}

// Set update val of key
func (db *DB) Set(bkt string, key string, val []byte) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}

	return db.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bkt))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), []byte(val))
	})
}

// Get get
func (db *DB) Get(bkt string, key string) []byte {
	var val []byte
	if err := db.client.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bkt))
		if err != nil {
			return err
		}
		val = bucket.Get([]byte(key))
		return nil
	}); err != nil {
		return nil
	}
	return val
}
