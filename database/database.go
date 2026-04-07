package database

// Database defines all of the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	DeleteDeviceByKey(key string) error                     //Delete specified device
	RotateDeviceKey(oldKey, deviceToken string) (string, error) //Atomically rotate device key after verifying ownership via deviceToken
	Close() error                                           //Close the database
}
