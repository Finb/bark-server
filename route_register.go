package main

import (
	"github.com/gin-gonic/gin"

	"github.com/lithammer/shortuuid"

	"go.etcd.io/bbolt"

	"github.com/mritd/logger"
)

const (
	bucketName = "device"
)

func init() {
	registerRoute("register", func(router *gin.Engine) {
		router.POST("/register", func(c *gin.Context) {
			deviceInfo := struct {
				DeviceKey   string `form:"device_key,omitempty" json:"device_key,omitempty"`
				DeviceToken string `form:"device_token,omitempty" json:"device_token,omitempty"`

				// compatible with old req
				OldDeviceKey   string `form:"key,omitempty" json:"key,omitempty"`
				OldDeviceToken string `form:"devicetoken,omitempty" json:"devicetoken,omitempty"`
			}{}
			err := c.Bind(&deviceInfo)
			if err != nil {
				c.JSON(400, failed(400, "request bind failed: %v", err))
				return
			}

			if deviceInfo.DeviceKey == "" && deviceInfo.OldDeviceKey != "" {
				deviceInfo.DeviceKey = deviceInfo.OldDeviceKey
			}

			if deviceInfo.DeviceToken == "" {
				if deviceInfo.OldDeviceToken != "" {
					deviceInfo.DeviceToken = deviceInfo.OldDeviceToken
				} else {
					c.JSON(400, failed(400, "device token is empty"))
					return
				}
			}

			err = db.Update(func(tx *bbolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
				if err != nil {
					return err
				}

				// If the deviceKey is empty or the corresponding deviceToken cannot be obtained from the database,
				// it is considered as a new device registration
				if deviceInfo.DeviceKey == "" || bucket.Get([]byte(deviceInfo.DeviceKey)) == nil {
					// Generate a new UUID as the deviceKey when a new device register
					deviceInfo.DeviceKey = shortuuid.New()
				}

				// update the deviceToken
				return bucket.Put([]byte(deviceInfo.DeviceKey), []byte(deviceInfo.DeviceToken))
			})

			if err != nil {
				logger.Errorf("route [/register] error: device registration failed: %v", err)
				c.JSON(500, failed(500, "device registration failed: %v", err))
				return
			}

			logger.Infof("route [/register]: new device registered successfully, deviceKey %s, deviceToken %s", deviceInfo.DeviceKey, deviceInfo.DeviceToken)
			c.JSON(200, data(map[string]string{
				// compatible with old resp
				"key":          deviceInfo.DeviceKey,
				"device_key":   deviceInfo.DeviceKey,
				"device_token": deviceInfo.DeviceToken,
			}))
		})
	})
}
