package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
	"go.etcd.io/bbolt"
)

// BboltDB implement Database interface with ETCD's bbolt
type BboltDB struct {
}

var dbOnce sync.Once
var db *bbolt.DB

const (
	bucketName = "device"
)

func NewBboltdb(dataDir string) Database {
	bboltSetup(dataDir)

	return &BboltDB{}
}

// CountAll Fetch records count
func (d *BboltDB) CountAll() (int, error) {
	var keypairCount int
	err := db.View(func(tx *bbolt.Tx) error {
		keypairCount = tx.Bucket([]byte(bucketName)).Stats().KeyN
		return nil
	})

	if err != nil {
		return 0, err
	}

	return keypairCount, nil
}

// Close close the db file
func (d *BboltDB) Close() error {
	return db.Close()
}

// DeviceTokenByKey get device token of specified key
func (d *BboltDB) DeviceTokenByKey(key string) (string, error) {
	var token string
	err := db.View(func(tx *bbolt.Tx) error {
		if bs := tx.Bucket([]byte(bucketName)).Get([]byte(key)); bs == nil {
			return fmt.Errorf("failed to get [%s] device token from database", key)
		} else {
			token = string(bs)
			return nil
		}
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// SaveDeviceToken create or update device token of specified key
func (d *BboltDB) SaveDeviceTokenByKey(key, deviceToken string) (string, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		// If the deviceKey is empty or the corresponding deviceToken cannot be obtained from the database,
		// it is considered as a new device registration
		if key == "" || bucket.Get([]byte(key)) == nil {
			// Generate a new UUID as the deviceKey when a new device register
			key = shortuuid.New()
		}

		// update the deviceToken
		return bucket.Put([]byte(key), []byte(deviceToken))
	})

	if err != nil {
		return "", err
	}

	return key, nil
}

// bboltSetup setup the bbolt database
func bboltSetup(dataDir string) {
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
