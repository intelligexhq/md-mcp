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
	// NotesDir is the absolute path to the directory where markdown files are stored.
	NotesDir string
}

// RegisterAll registers every available tool with the given MCP server.
// This is the single entry point that main.go calls to wire up all tool
// definitions and their handlers.
func RegisterAll(s *server.MCPServer, cfg Config) {
	registerMarkdownTools(s, cfg)
	registerDateTimeTools(s)
}

// getSafePath resolves a filename relative to the configured notes directory.
// It enforces three invariants:
//  1. The filename must end with ".md" (appends it if missing).
//  2. The resulting path is cleaned to eliminate traversal sequences (e.g. "../").
//  3. The final absolute path must remain within the permitted notes directory.
//
// If any check fails an error is returned; otherwise the absolute safe path is returned.
func getSafePath(notesDir, filename string) (string, error) {
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	targetPath := filepath.Clean(filepath.Join(notesDir, filename))

	absNotesDir, _ := filepath.Abs(notesDir)
	absTarget, _ := filepath.Abs(targetPath)

	if !strings.HasPrefix(absTarget, absNotesDir) {
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
