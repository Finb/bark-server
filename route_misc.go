package main

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("misc", func(router fiber.Router) {
		// ping func only returns a "pong" string, usually used to test server response
		router.Get("/ping", func(c *fiber.Ctx) error {
			return c.JSON(CommonResp{
				Code:      200,
				Message:   "pong",
				Timestamp: time.Now().Unix(),
			})
		})

		// healthz func only returns a "ok" string, similar to ping func,
		// healthz func is usually used for health check
		router.Get("/healthz", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		// info func returns information about the server version
		router.Get("/info", func(c *fiber.Ctx) error {
			devices, _ := db.CountAll()
			return c.JSON(map[string]interface{}{
				"version": version,
				"build":   buildDate,
				"arch":    runtime.GOOS + "/" + runtime.GOARCH,
				"commit":  commitID,
				"devices": devices,
			})
		})
	})
}
