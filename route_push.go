package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/finb/bark-server/v2/apns"

	"github.com/gofiber/fiber/v2"

	"go.etcd.io/bbolt"
)

func init() {
	registerRoute("push", func(router *fiber.App) {
		router.Post("/push", func(c *fiber.Ctx) error { return routeDoPush(c, false) })
	})

	// compatible with old requests
	registerRouteWithWeight("push_compat", 1, func(router *fiber.App) {
		router.Get("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })

		router.Get("/:device_key/:category/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
		router.Post("/:device_key/:category/:title/:body", func(c *fiber.Ctx) error { return routeDoPush(c, true) })
	})
}

func routeDoPush(c *fiber.Ctx, compat bool) error {
	// default value
	msg := apns.PushMessage{
		Category:  "myNotificationCategory",
		Body:      "NoContent",
		Sound:     "1107",
		ExtParams: make(map[string]interface{}),
	}

	// always parse body(Lowest priority)
	if err := c.BodyParser(&msg); err != nil && err != fiber.ErrUnprocessableEntity {
		return c.Status(400).JSON(failed(400, "request bind failed: %v", err))
	}

	if compat {
		params := make(map[string]string)
		visitor := func(key, value []byte) {
			params[strings.ToLower(string(key))] = string(value)
		}
		// parse query args (medium priority)
		c.Request().URI().QueryArgs().VisitAll(visitor)
		// parse post args
		c.Request().PostArgs().VisitAll(visitor)

		// parse multipartForm values
		form, err := c.Request().MultipartForm()
		if err == nil {
			for key,val := range form.Value {
				if len(val) > 0{
					params[key] = val[0]
				}
			}
		}


		for key, val := range params {
			switch strings.ToLower(string(key)) {
			case "device_key":
				msg.DeviceKey = val
			case "category":
				msg.Category = val
			case "title":
				msg.Title = val
			case "body":
				msg.Body = val
			case "sound":
				msg.Sound = val + ".caf"
			default:
				msg.ExtParams[strings.ToLower(string(key))] = val
			}
		}

		// parse url path (highest priority)
		if pathDeviceKey := c.Params("device_key"); pathDeviceKey != "" {
			msg.DeviceKey = pathDeviceKey
		}
		if category := c.Params("category"); category != "" {
			str, err := url.QueryUnescape(category)
			if err != nil {
				return err
			}
			msg.Category = str
		}
		if title := c.Params("title"); title != "" {
			str, err := url.QueryUnescape(title)
			if err != nil {
				return err
			}
			msg.Title = str
		}
		if body := c.Params("body"); body != "" {
			str, err := url.QueryUnescape(body)
			if err != nil {
				return err
			}
			msg.Body = str
		}
	}

	if msg.DeviceKey == "" {
		return c.Status(400).JSON(failed(400, "device key is empty"))
	}

	err := db.View(func(tx *bbolt.Tx) error {
		if bs := tx.Bucket([]byte(bucketName)).Get([]byte(msg.DeviceKey)); bs == nil {
			return fmt.Errorf("failed to get [%s] device token from database", msg.DeviceKey)
		} else {
			msg.DeviceToken = string(bs)
			return nil
		}
	})
	if err != nil {
		return c.Status(400).JSON(failed(400, "failed to get device token: %v", err))
	}

	err = apns.Push(&msg)
	if err != nil {
		return c.Status(500).JSON(failed(500, "push failed: %v", err))
	}
	return c.JSON(success())
}
