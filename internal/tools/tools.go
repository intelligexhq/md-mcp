package tools

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Config holds the configuration needed by the tools package to operate.
// It decouples the tools from global state, making them testable and reusable.
type Config struct {
	// WorkspaceDir is the absolute path to the directory where markdown files are stored.
	WorkspaceDir string
}

// RegisterAll registers every available tool with the given MCP server.
// This is the single entry point that main.go calls to wire up all tool
// definitions and their handlers.
func RegisterAll(s *server.MCPServer, cfg Config) {
	registerMarkdownTools(s, cfg)
	registerDateTimeTools(s)
}

// getSafePath resolves a filename relative to the configured workspace directory.
// It enforces three invariants:
//  1. The filename must end with ".md" (appends it if missing).
//  2. The resulting path is cleaned to eliminate traversal sequences (e.g. "../").
//  3. The final absolute path must remain within the permitted workspace directory.
//
// If any check fails an error is returned; otherwise the absolute safe path is returned.
func getSafePath(workspaceDir, filename string) (string, error) {
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	targetPath := filepath.Clean(filepath.Join(workspaceDir, filename))

	absWorkspaceDir, _ := filepath.Abs(workspaceDir)
	absTarget, _ := filepath.Abs(targetPath)

	if !strings.HasPrefix(absTarget, absWorkspaceDir) {
		return "", fmt.Errorf("security violation: access denied outside allowed directory")
	}

	return absTarget, nil
}

// toolDef is a small helper that pairs an MCP tool schema with its handler
// function, keeping the registration code in each file concise.
type toolDef struct {
	tool    mcp.Tool
	handler server.ToolHandlerFunc
}
