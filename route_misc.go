package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
)

func init() {
	registerRoute("misc", func(router *gin.Engine) {
		// ping func only returns a "pong" string, usually used to test server response
		router.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})

		// healthz func only returns a "ok" string, similar to ping func,
		// healthz func is usually used for health check
		router.GET("/healthz", func(c *gin.Context) {
			c.String(200, "ok")
		})

		// info func returns information about the server version
		router.GET("/info", func(c *gin.Context) {
			c.JSON(200, map[string]string{
				"version": version,
				"build":   buildDate,
				"arch":    runtime.GOOS + "/" + runtime.GOARCH,
				"commit":  commitID,
			})
		})
	})
}
