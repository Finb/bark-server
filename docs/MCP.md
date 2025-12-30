## MCP

Bark supports the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) via HTTP Streamable, allowing AI agents (like Claude Desktop, Cherry Studio or n8n) to send notifications directly through Bark.

### Endpoints

| Endpoint           | Description                                                                                                   |
| ------------------ | ------------------------------------------------------------------------------------------------------------- |
| `/mcp`             | Generic MCP endpoint. Requires `device_key` to be provided in the tool arguments.                             |
| `/mcp/:device_key` | Device-specific MCP endpoint. The `device_key` is fixed by the URL, and the AI agent doesn't need to know it. |

### Examples

Cherry Studio:   
```json
{
  "mcpServers": {
    "bark": {
      "type": "streamableHttp",
      "url": "https://api.day.app/mcp/{key}"
    }
  }
}
```

VS Code:  
```js
{
  "servers": {
    "bark": {
      "type": "http",
      "url": "https://api.day.app/mcp/{key}"
    }
  }
}
```

Claude Code:   
```sh
claude mcp add bark --transport http https://api.day.app/mcp/{key}
```  
or  
```js
{
  "mcpServers": {
    "bark": {
      "type": "http",
      "url": "https://api.day.app/mcp/{key}"
    }
  }
}
```  
> Note: Replace {key} in the URL with your own key.

