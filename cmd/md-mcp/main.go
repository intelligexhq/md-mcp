package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"local-markdown-mcp/internal/tools"

	"github.com/mark3labs/mcp-go/server"
)

// version, commit, and date are injected at build time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var showVersion bool
	var workspaceDir string

	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "tmp", "my_agent_notes")

	flag.BoolVar(&showVersion, "version", false, "Print version information and exit")
	flag.StringVar(&workspaceDir, "workspace", defaultDir, "Path to the workspace directory")
	flag.Parse()

	if showVersion {
		fmt.Printf("md-mcp %s (commit=%s, built=%s)\n", version, commit, date)
		return
	}

	_ = os.MkdirAll(workspaceDir, 0o755)

	// Set up structured JSON logging to stderr so it does not interfere with
	// the stdio-based MCP transport.
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// Create the MCP server instance.
	s := server.NewMCPServer(
		"md-mcp", version,
		server.WithLogger(logger),
	)

	// Register all tools via the tools package, passing in the config they need.
	tools.RegisterAll(s, tools.Config{
		WorkspaceDir: workspaceDir,
	})

	slog.Info("starting MCP server", "name", "md-mcp", "version", version, "workspace", workspaceDir)

	// Start serving over stdin/stdout — this blocks until the client disconnects.
	if err := server.ServeStdio(s); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
