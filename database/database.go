package database

// Database defines all of the db operation
type Database interface {
	CountAll() (int, error)                                 //Get db records count
	DeviceTokenByKey(key string) (string, error)            //Get specified device's token
	SaveDeviceTokenByKey(key, token string) (string, error) //Create or update specified devices's token
	GetMarkdownByKey(key string) (string, error)            // Get Markdown Content by key
	SaveMarkdown(content string) (string, error)            // Save Markdown Content
	Close() error                                           //Close the database
}
