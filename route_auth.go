package main

import (
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberbasicauth "github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/mritd/logger"
)

func routerAuth(user, passwd string, router fiber.Router, urlPrefix string) {
	if user == "" && passwd == "" {
		logger.Info("Bark Server Has No Basic Auth.")
		return
	}

	logger.Info("Bark Server Has Basic Auth Enabled.")
	authFreeRouters := []string{"/ping", "/register", "/healthz"}
	basicAuth := fiberbasicauth.New(fiberbasicauth.Config{
		Users: map[string]string{user: passwd},
		Realm: "Coffee Time",
		Unauthorized: func(c *fiber.Ctx) error {
			for _, item := range authFreeRouters {
				if strings.HasPrefix(c.Path(), path.Join(urlPrefix, item)) {
					return c.Next()
				}
			}
			return c.Status(418).SendString("I'm a teapot")
		},
	})

	router.Use("/+", basicAuth)
}
