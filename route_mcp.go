package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func init() {
	registerRoute("mcp", func(router fiber.Router) {
		mcpStreamable := setupMCPServer()
		router.All("/mcp", func(c *fiber.Ctx) error {
			return adaptor.HTTPHandlerFunc(mcpStreamable.ServeHTTP)(c)
		})
		router.All("/mcp/:device_key", func(c *fiber.Ctx) error {
			mcpWithDevice := setupMCPServer(c.Params("device_key"))
			return adaptor.HTTPHandlerFunc(mcpWithDevice.ServeHTTP)(c)
		})
	})
}

func setupMCPServer(defaultDeviceKey ...string) *server.StreamableHTTPServer {
	mcpServer := server.NewMCPServer(
		"Bark MCP Server",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// Register "notify" tool
	// device_key is required when defaultDeviceKey is not provided
	toolOpts := []mcp.ToolOption{
		mcp.WithDescription("Send a notification to a device via Bark"),
		mcp.WithString("title", mcp.Description("Notification title")),
		mcp.WithString("subtitle", mcp.Description("Notification subtitle")),
		mcp.WithString("body", mcp.Required(), mcp.Description("Notification body")),
		mcp.WithString("level",
			mcp.Description("Notification level, can be 'critical', 'active', 'timeSensitive', 'passive'"),
			mcp.Enum("critical", "active", "timeSensitive", "passive"),
		),
		mcp.WithNumber("badge", mcp.Description("Badge number")),
		mcp.WithString("sound", mcp.Description("Notification sound")),
		mcp.WithString("icon", mcp.Description("Notification icon URL")),
		mcp.WithString("group", mcp.Description("Notification group")),
		mcp.WithString("url", mcp.Description("Click action URL")),
		mcp.WithString("copy", mcp.Description("Text to copy on copy action")),
	}
	if len(defaultDeviceKey) == 0 {
		toolOpts = append(
			toolOpts,
			mcp.WithString("device_key", mcp.Required(), mcp.Description("Device Key")),
		)
	}

	notifyTool := mcp.NewTool("notify", toolOpts...)
	mcpServer.AddTool(notifyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argsJson, err := json.Marshal(request.Params.Arguments)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal arguments"), nil
		}

		var args map[string]any
		err = json.Unmarshal(argsJson, &args)
		if err != nil {
			return mcp.NewToolResultError("Failed to unmarshal arguments"), nil
		}

		// Handle device_key
		if _, ok := args["device_key"]; !ok && len(defaultDeviceKey) > 0 {
			args["device_key"] = defaultDeviceKey[0]
		}
		if _, ok := args["device_key"]; !ok {
			return mcp.NewToolResultError("device_key is required"), nil
		}

		// Use existing push function
		// We need to adapt the push function or its usage.
		// push(params map[string]interface{}) (int, error) is defined in route_push.go

		code, err := push(args)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to send notification: %v (code %d)", err, code)), nil
		}

		return mcp.NewToolResultText("Notification sent successfully"), nil
	})

	return server.NewStreamableHTTPServer(mcpServer)
}
