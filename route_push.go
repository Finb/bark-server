package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/finb/bark-server/v2/pusher"

	"github.com/gofiber/fiber/v2"
)

func init() {
	// compatible with old requests
	registerRouteWithWeight("push_compat", 1, func(router fiber.Router) {
		// deprecated: this rules will match /favicon.ico.
		// Use v2 api instead.
		router.Get("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c) })
		router.Post("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c) })

		router.Get("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })
		router.Post("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })

		router.Get("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })
		router.Post("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })

		router.Get("/:device_key/:title/:subtitle/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })
		router.Post("/:device_key/:title/:subtitle/:body", func(c *fiber.Ctx) error { return routeDoPush(c) })
	})
}

func routeDoPush(c *fiber.Ctx) error {
	deviceKey := c.Params("device_key")
	params := getAllParams(c)
	delete(params, "device_key")
	if err := push(deviceKey, params); err != nil {
		return c.Status(500).JSON(failed(500, err.Error()))
	} else {
		return c.JSON(success())
	}
}
func batchPush(deviceKeys []string, params interface{}) map[string]error {
	var wg sync.WaitGroup
	result := make(map[string]error)
	var mu sync.Mutex

	for _, deviceKey := range deviceKeys {
		wg.Add(1)
		go func(deviceKey string) {
			defer wg.Done()
			if err := push(deviceKey, params); err != nil {
				result[deviceKey] = err
			}
			mu.Lock()
			mu.Unlock()
		}(deviceKey)
	}
	wg.Wait()
	return result
}

func push(deviceKey string, params interface{}) (err error) {
	var deviceToken string
	if deviceToken, err = db.DeviceTokenByKey(strings.Trim(deviceKey, " ")); err != nil {
		return fmt.Errorf("failed to get device token: %v", err)
	}
	msg := pusher.NewApnsMessage()
	msg.SetDeviceToken(deviceToken)
	if err := msg.LoadMessages(params); err != nil {
		return fmt.Errorf("failed to load message: %v", err)
	}
	if err := msg.Check(); err != nil {
		return fmt.Errorf("invalid message: %v", err)
	}
	pusherProvider := pusher.PushProviderApns
	err = pusher.GetPusher(pusherProvider, maxApnClientCount).Push(msg)
	if err != nil {
		return fmt.Errorf("push failed: %v", err)
	}
	return
}

// Maximum number of APN clients allowed, -1 means no limit
var maxApnClientCount = 1

// Set the maximum number of APN clients allowed
func SetMaxApnClientCount(count int) {
	maxApnClientCount = count
}
