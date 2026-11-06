package database

import "encoding/json"

// Supported device types
const (
	DeviceTypeIOS     = "ios"
	DeviceTypeWatchOS = "watchos"
)

// DeviceRecord describes a registered device's APNs token and device type.
type DeviceRecord struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	UpdatedAt int64  `json:"updated_at"`
}

// normalizeDeviceType returns the default device type when it is empty.
func normalizeDeviceType(deviceType string) string {
	if deviceType == "" {
		return DeviceTypeIOS
	}
	return deviceType
}

// parseDeviceRecord parses a stored device record.
// Compatible with legacy raw token values, which are treated as iOS devices.
func parseDeviceRecord(raw string) *DeviceRecord {
	record := &DeviceRecord{}
	if err := json.Unmarshal([]byte(raw), record); err != nil {
		return &DeviceRecord{Token: raw, Type: DeviceTypeIOS}
	}
	record.Type = normalizeDeviceType(record.Type)
	return record
}

// Database defines all of the db operation
type Database interface {
	CountAll() (int, error)                                             //Get db records count
	DeviceTokenByKey(key string) (*DeviceRecord, error)                 //Get specified device's token
	SaveDeviceTokenByKey(key, token, deviceType string) (string, error) //Create or update specified devices's token
	DeleteDeviceByKey(key string) error                                 //Delete specified device
	Close() error                                                       //Close the database
}
