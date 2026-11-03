package database

import (
	"fmt"
	"strings"
)

var (
	cacheKey          = "MemoryBaseKey"
	cacheDeviceRecord = &DeviceRecord{Type: DeviceTypeIOS}
)

type MemBase struct {
}

func NewMemBase() Database {
	return &MemBase{}
}

func (d *MemBase) CountAll() (int, error) {
	return 1, nil
}

func (d *MemBase) DeviceTokenByKey(key string) (*DeviceRecord, error) {
	if cacheKey == key && cacheDeviceRecord.Token != "" {
		return &DeviceRecord{Token: cacheDeviceRecord.Token, Type: cacheDeviceRecord.Type}, nil
	}
	return nil, fmt.Errorf("key not found")
}

func (d *MemBase) SaveDeviceTokenByKey(key, token, deviceType string) (string, error) {
	if key != "" && key != cacheKey {
		return "", fmt.Errorf("key not found")
	}
	// Deep copy prevents Fiber memory overwrite bugs.
	cacheDeviceRecord = &DeviceRecord{
		Token: strings.Clone(token),
		Type:  normalizeDeviceType(deviceType),
	}
	return key, nil
}

func (d *MemBase) DeleteDeviceByKey(key string) error {
	if key != "" && key != cacheKey {
		return fmt.Errorf("key not found")
	}
	cacheDeviceRecord = &DeviceRecord{Type: DeviceTypeIOS}
	return nil
}

func (d *MemBase) Close() error {
	return nil
}
