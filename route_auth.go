package main

import (
	"github.com/gofiber/fiber/v2"
	fiberbasicauth "github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/mritd/logger"
	"strings"
)

func routerAuth(user, passwd string, router fiber.Router) {
	if user != "" && passwd != "" {
		logger.Info("Bark Server Has Basic Auth Enabled.")
		basicAuth := fiberbasicauth.New(fiberbasicauth.Config{
			Users: map[string]string{user: passwd},
			Realm: "Coffee Time",
			Unauthorized: func(c *fiber.Ctx) error {
				authFreeRouters := []string{"/ping", "/register", "/healthz"}
				for _, item := range authFreeRouters {
					if strings.HasPrefix(c.Path(), item) {
						return c.Next()
					}
				}
				return c.Status(418).SendString("I'm a teapot")
			},
		})

		router.Use("/+", basicAuth)
	}
}
