package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/mritd/logger"
	"go.etcd.io/bbolt"
)

var dbOnce sync.Once
var db *bbolt.DB

func databaseSetup(dataDir string) {
	dbOnce.Do(func() {
		logger.Infof("init database [%s]...", dataDir)
		if _, err := os.Stat(dataDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dataDir, 0755); err != nil {
				logger.Fatalf("failed to create database storage dir(%s): %v", dataDir, err)
			}
		} else if err != nil {
			logger.Fatalf("failed to open database storage dir(%s): %v", dataDir, err)
		}

		bboltDB, err := bbolt.Open(filepath.Join(dataDir, "bark.db"), 0600, nil)
		if err != nil {
			logger.Fatalf("failed to create database file(%s): %v", filepath.Join(dataDir, "bark.db"), err)
		}
		err = bboltDB.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			return err
		})
		if err != nil {
			logger.Fatalf("failed to create database bucket: %v", err)
		}
		db = bboltDB
	})
}
