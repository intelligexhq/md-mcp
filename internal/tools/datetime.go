package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerDateTimeTools defines and registers the date/time tool with the MCP server.
// This tool gives an AI agent awareness of the current moment, which is otherwise
// unavailable to it since LLMs have no built-in clock.
func registerDateTimeTools(s *server.MCPServer) {
	def := toolDef{
		tool: mcp.NewTool(
			"get_current_datetime",
			mcp.WithDescription("Returns the current date and time in ISO 8601 format (UTC) and the local timezone."),
		),
		handler: getDateTimeHandler,
	}
	s.AddTool(def.tool, def.handler)
}

// getDateTimeHandler returns the current time in both UTC and the server's
// local timezone, formatted as RFC 3339 (ISO 8601 compliant). The output
// includes the timezone abbreviation so the agent can interpret context.
func getDateTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	now := time.Now()
	localZone, _ := now.Zone()
	return mcp.NewToolResultText(
		fmt.Sprintf(
			"UTC:   %s\nLocal: %s (%s)",
			now.UTC().Format(time.RFC3339),
			now.Format(time.RFC3339),
			localZone,
		),
	), nil
}
