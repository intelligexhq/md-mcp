package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"
	"local-markdown-mcp/internal/tools"
)

func main() {
	// Resolve and ensure the notes directory exists.
	home, _ := os.UserHomeDir()
	notesDir := filepath.Join(home, "tmp", "my_agent_notes")
	_ = os.MkdirAll(notesDir, 0755)

	// Set up structured JSON logging to stderr so it does not interfere with
	// the stdio-based MCP transport.
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// Create the MCP server instance.
	s := server.NewMCPServer("md-mcp", "0.1.0",
		server.WithLogger(logger),
	)

	// Register all tools via the tools package, passing in the config they need.
	tools.RegisterAll(s, tools.Config{
		NotesDir: notesDir,
	})

	slog.Info("starting MCP server", "name", "md-mcp", "version", "0.1.0")

	// Start serving over stdin/stdout — this blocks until the client disconnects.
	if err := server.ServeStdio(s); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
