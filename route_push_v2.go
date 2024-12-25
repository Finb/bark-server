package main

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("push", func(router fiber.Router) {
		router.Post("/push", func(c *fiber.Ctx) error { return routeDoPushV2(c) })
	})
}

type PushV2Request struct {
	Title             string `json:"title,omitempty"`
	SubTitle          string `json:"subtitle,omitempty"`
	Body              string `json:"body,omitempty"`
	DeviceKey         string `json:"device_key"`
	Level             int    `json:"level,omitempty"`
	Badge             int    `json:"badge,omitempty"`
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty"` // deprecated
	AutoCopy          string `json:"autoCopy,omitempty"`
	Copy              string `json:"copy,omitempty"`
	Sound             string `json:"sound,omitempty"`
	Icon              string `json:"icon,omitempty"`
	Group             string `json:"group,omitempty"`
	IsArchive         int    `json:"isArchive,omitempty"`
	Url               string `json:"url,omitempty"`
	Call              string `json:"call,omitempty"`
	Ciphertext        string `json:"ciphertext,omitempty"`
}

func routeDoPushV2(c *fiber.Ctx) error {
	params := &PushV2Request{}

	if err := c.BodyParser(params); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
	}

	var deviceKeys []string
	if strings.Contains(params.DeviceKey, ",") {
		deviceKeys = strings.Split(params.DeviceKey, ",")
	} else {
		deviceKeys = []string{params.DeviceKey}
	}
	params.DeviceKey = ""
	count := len(deviceKeys)
	if count == 0 {
		return c.Status(400).JSON(failed(400, "unknown device key"))
	}
	if count > 1 {
		if count > maxBatchPushCount && maxBatchPushCount != -1 {
			return c.Status(400).JSON(failed(400, "batch push count exceeds the maximum limit: %d", maxBatchPushCount))
		}
		if batchErrors := batchPush(deviceKeys, params); batchErrors != nil {
			var result []map[string]interface{}
			for deviceKey, err := range batchErrors {
				result = append(result, map[string]interface{}{
					"device_key": deviceKey,
					"message":    err.Error(),
				})
			}
			return c.JSON(data(result))
		}
	} else {
		if err := push(deviceKeys[0], params); err != nil {
			return c.Status(500).JSON(failed(500, "push failed: %v", err))
		}
	}
	return c.JSON(success())
}

// Maximum number of batch pushes allowed, -1 means no limit
var maxBatchPushCount = -1

// Set the maximum number of batch pushes allowed
func SetMaxBatchPushCount(count int) {
	maxBatchPushCount = count
}
