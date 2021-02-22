package main

import (
	"github.com/gofiber/fiber/v2"
	fiberbasicauth "github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/mritd/logger"
)

func routerAuth(user, passwd string, router *fiber.App) {
	if user != "" && passwd != "" {
		logger.Info("Bark Server Has Basic Auth Enabled.")
		basicAuth := fiberbasicauth.New(fiberbasicauth.Config{
			Users: map[string]string{user: passwd},
			Realm: "Coffee Time",
			Unauthorized: func(c *fiber.Ctx) error {
				return c.Status(418).SendString("I'm a teapot")
			},
		})
		router.Use("/push", basicAuth)
		router.Use("/:device_key/:body", basicAuth)
		router.Use("/:device_key/:title/:body", basicAuth)
		router.Use("/:device_key/:category/:title/:body", basicAuth)
	}
}
