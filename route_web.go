package main

import (
	"github.com/gofiber/fiber/v2"
)

func init() {
	registerRoute("web", func(router fiber.Router) {
		// bark markdown web
		router.Static("/web", "./web")
		router.Get("/web/:markdown_key", func(c *fiber.Ctx) error { return c.SendFile("./web/index.html") })
		router.Get("/web/md/:markdown_key", func(c *fiber.Ctx) error { return routeDoMarkdown(c) })
	})
}

func routeDoMarkdown(c *fiber.Ctx) error {
	markdownKey := c.Params("markdown_key")
	if markdownKey == "" {
		return c.Status(400).JSON(failed(400, "markdown_key is empty"))
	}
	content, err := db.GetMarkdownByKey(markdownKey)
	if err != nil {
		return c.Status(400).JSON(failed(400, "get markdown failed: %v", err))
	}
	return c.Status(200).JSON(data(map[string]string{
		"content": content,
	}))
}
