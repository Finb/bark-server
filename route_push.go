package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

func init() {
	registerRoute("push", func(router *gin.Engine) {
		router.POST("/push", func(c *gin.Context) {
			var pushInfo PushInfo
			err := c.Bind(&pushInfo)
			if err != nil {
				c.JSON(400, failed(400, "request bind failed: %v", err))
				return
			}

			if pushInfo.DeviceKey == "" {
				if pushInfo.OldDeviceKey != "" {
					pushInfo.DeviceKey = pushInfo.OldDeviceKey
				} else {
					c.JSON(400, failed(400, "device token is empty"))
					return
				}
			}

			err = db.View(func(tx *bolt.Tx) error {
				if bs := tx.Bucket([]byte(bucketName)).Get([]byte(pushInfo.DeviceKey)); bs == nil {
					return fmt.Errorf("failed to get [%s] device token from database", pushInfo.DeviceKey)
				} else {
					pushInfo.DeviceToken = string(bs)
					return nil
				}
			})
			if err != nil {
				c.JSON(400, failed(400, "failed to get device token: %v", err))
				return
			}
			err = apnsPush(&pushInfo)
			if err != nil {
				c.JSON(500, failed(500, "push failed: %v", err))
				return
			}
		})
	})
}
