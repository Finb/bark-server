package main

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mritd/logger"
)

// rotateRateLimit enforces a per-key cooldown to prevent spam rotation.
// A key can only be rotated at most once per rotateMinInterval.
const rotateMinInterval = 5 * time.Minute

var (
	rotateLastTime   = map[string]time.Time{}
	rotateLastTimeMu sync.Mutex
)

func init() {
	registerRoute("rotate", func(router fiber.Router) {
		router.Post("/register/rotate", doRotateKey)
	})
}

func doRotateKey(c *fiber.Ctx) error {
	var deviceInfo DeviceInfo
	if err := c.BodyParser(&deviceInfo); err != nil {
		return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
	}

	if deviceInfo.DeviceKey == "" && deviceInfo.OldDeviceKey != "" {
		deviceInfo.DeviceKey = deviceInfo.OldDeviceKey
	}

	if deviceInfo.DeviceKey == "" {
		return c.Status(400).JSON(failed(400, "device_key is required"))
	}

	// Rate-limit: reject if this key was rotated too recently.
	rotateLastTimeMu.Lock()
	lastTime, seen := rotateLastTime[deviceInfo.DeviceKey]
	if seen && time.Since(lastTime) < rotateMinInterval {
		rotateLastTimeMu.Unlock()
		remaining := rotateMinInterval - time.Since(lastTime)
		return c.Status(429).JSON(failed(429, "rotation rate limit exceeded, please wait %ds before rotating again", int(remaining.Seconds())+1))
	}
	// Record the attempt before releasing the lock so concurrent requests are blocked.
	rotateLastTime[deviceInfo.DeviceKey] = time.Now()
	rotateLastTimeMu.Unlock()

	newKey, err := db.RotateDeviceKey(deviceInfo.DeviceKey)
	if err != nil {
		// On failure, clear the rate-limit timestamp so the user can retry.
		rotateLastTimeMu.Lock()
		delete(rotateLastTime, deviceInfo.DeviceKey)
		rotateLastTimeMu.Unlock()

		logger.Errorf("device key rotation failed: %v", err)
		return c.Status(500).JSON(failed(500, "device key rotation failed: %v", err))
	}

	// After a successful rotation, track the new key's cooldown instead.
	rotateLastTimeMu.Lock()
	rotateLastTime[newKey] = time.Now()
	delete(rotateLastTime, deviceInfo.DeviceKey)
	rotateLastTimeMu.Unlock()

	logger.Infof("device key rotated: %s -> %s", deviceInfo.DeviceKey, newKey)

	return c.Status(200).JSON(data(map[string]string{
		"device_key": newKey,
		// compatible with old field name
		"key": newKey,
	}))
}
