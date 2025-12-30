package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type contextKey string

const deviceKeyCtxKey contextKey = "device_key"

func init() {
	registerRoute("mcp", func(router fiber.Router) {
		mcpGenericStreamable := setupGenericMCPServer()
		mcpSpecificStreamable := setupSpecificMCPServer()

		// Basic endpoint - requires device_key in tool arguments
		router.All("/mcp", func(c *fiber.Ctx) error {
			return adaptor.HTTPHandlerFunc(mcpGenericStreamable.ServeHTTP)(c)
		})

		// Device-specific endpoint - device_key is pre-filled from URL path
		router.All("/mcp/:device_key", func(c *fiber.Ctx) error {
			deviceKey := c.Params("device_key")
			return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), deviceKeyCtxKey, deviceKey)
				mcpSpecificStreamable.ServeHTTP(w, r.WithContext(ctx))
			})(c)
		})
	})
}

func setupGenericMCPServer() *server.StreamableHTTPServer {
	s := server.NewMCPServer("Bark MCP Server", version,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	opts := getCommonToolOpts()
	opts = append(opts,
		mcp.WithString("device_key",
			mcp.Required(),
			mcp.Description("Device Key"),
		),
	)

	s.AddTool(mcp.NewTool("notify", opts...), notifyHandler)
	return server.NewStreamableHTTPServer(s)
}

func setupSpecificMCPServer() *server.StreamableHTTPServer {
	s := server.NewMCPServer("Bark MCP Server (Specific)", version,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	s.AddTool(mcp.NewTool("notify", getCommonToolOpts()...), notifyHandler)
	return server.NewStreamableHTTPServer(s)
}

func notifyHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := request.Params.Arguments.(map[string]any)
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format"), nil
	}

	// Resolve device_key: tool args > context (from URL)
	var deviceKey string
	if val, ok := args["device_key"]; ok {
		if tmpDeviceKey, ok := val.(string); ok {
			deviceKey = tmpDeviceKey
		}
	}
	if val := ctx.Value(deviceKeyCtxKey); val != nil {
		if tmpDeviceKey, ok := val.(string); ok {
			deviceKey = tmpDeviceKey
		}
	}
	if len(deviceKey) == 0 {
		return mcp.NewToolResultError("device_key is required"), nil
	}

	args["device_key"] = deviceKey
	code, err := push(args)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send notification: %v (code %d)", err, code)), nil
	}
	return mcp.NewToolResultText("Notification sent successfully"), nil
}

func getCommonToolOpts() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("Send a notification to a device via Bark"),
		mcp.WithString("title", mcp.Description("Notification title")),
		mcp.WithString("subtitle", mcp.Description("Notification subtitle")),
		mcp.WithString("body", mcp.Description("Notification content")),
		mcp.WithString("markdown", mcp.Description("Basic Markdown notification content. Overrides body.")),
		mcp.WithString("level",
			mcp.Description("Notification level"),
			mcp.Enum("critical", "active", "timeSensitive", "passive"),
		),
		mcp.WithNumber("volume",
			mcp.Description("Alert volume for important notification"),
			mcp.DefaultNumber(5),
			mcp.Max(10),
			mcp.Min(0),
		),
		mcp.WithNumber("badge", mcp.Description("Badge number")),
		mcp.WithString("call", mcp.Description("Set to '1' to repeat the notification ringtone")),
		mcp.WithString("sound", mcp.Description("Notification sound")),
		mcp.WithString("icon", mcp.Description("Notification icon URL")),
		mcp.WithString("image", mcp.Description("Notification image URL")),
		mcp.WithString("group", mcp.Description("Notification group")),
		mcp.WithString("isArchive", mcp.Description("Set to '1' to save the notification or any other value to skip saving")),
		mcp.WithString("url", mcp.Description("Click action URL")),
		mcp.WithString("copy", mcp.Description("Text to copy on copy action")),
	}
}
