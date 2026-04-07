package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mritd/logger"
)

func init() {
	registerRoute("rotate", func(router fiber.Router) {
		router.Post("/register/rotate", doRotateKey)
	})
}

// doRotateKey handles POST /register/rotate.
//
// Security model: device_key is public (it appears in push URLs). Anyone who
// knows the key could call this endpoint and rotate it, permanently invalidating
// the victim's push URL. To prevent this, the caller must also supply
// device_token — the APNs token that is only available on the physical device
// and is never embedded in push URLs. The server verifies the token matches
// what is stored before performing the rotation.
func doRotateKey(c *fiber.Ctx) error {
	var deviceInfo DeviceInfo
	if err := c.BodyParser(&deviceInfo); err != nil {
		return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
	}

	if deviceInfo.DeviceKey == "" && deviceInfo.OldDeviceKey != "" {
		deviceInfo.DeviceKey = deviceInfo.OldDeviceKey
	}
	if deviceInfo.DeviceToken == "" && deviceInfo.OldDeviceToken != "" {
		deviceInfo.DeviceToken = deviceInfo.OldDeviceToken
	}

	if deviceInfo.DeviceKey == "" {
		return c.Status(400).JSON(failed(400, "device_key is required"))
	}
	if deviceInfo.DeviceToken == "" {
		return c.Status(400).JSON(failed(400, "device_token is required for ownership verification"))
	}

	newKey, err := db.RotateDeviceKey(deviceInfo.DeviceKey, deviceInfo.DeviceToken)
	if err != nil {
		// Distinguish ownership failures from other errors so the client can
		// show a meaningful message without leaking internals.
		logger.Errorf("device key rotation failed for key %s: %v", deviceInfo.DeviceKey, err)
		if err.Error() == "device token mismatch: ownership verification failed" {
			return c.Status(401).JSON(failed(401, "ownership verification failed: device_token does not match"))
		}
		return c.Status(500).JSON(failed(500, "device key rotation failed: %v", err))
	}

	logger.Infof("device key rotated: %s -> %s", deviceInfo.DeviceKey, newKey)

	return c.Status(200).JSON(data(map[string]string{
		"device_key": newKey,
		// compatible with old field name
		"key": newKey,
	}))
}
