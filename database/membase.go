package database

import (
	"fmt"
)

var (
	cacheKey         = "MemoryBaseKey"
	cacheDeviceToken = ""
)

type MemBase struct {
}

func NewMemBase() Database {
	return &MemBase{}
}

func (d *MemBase) CountAll() (int, error) {
	return 1, nil
}

func (d *MemBase) DeviceTokenByKey(key string) (string, error) {
	if cacheKey == key && cacheDeviceToken != "" {
		return cacheDeviceToken, nil
	}
	return "nil", fmt.Errorf("key not found")
}

func (d *MemBase) SaveDeviceTokenByKey(key, token string) (string, error) {
	if key != "" && key != cacheKey {
		return "", fmt.Errorf("key not found")
	}
	cacheDeviceToken = token
	return key, nil
}

func (d *MemBase) Close() error {
	return nil
}
