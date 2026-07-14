package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

// Define a safe local directory for the agent's notes
var notesDir string

func main() {
	// Set up the notes folder on your Desktop
	home, _ := os.UserHomeDir()
	notesDir = filepath.Join(home, "tmp", "my_agent_notes")
	_ = os.MkdirAll(notesDir, 0755)

	// setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// setup hooks for request logging
	hooks := &server.Hooks{}
	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		slog.Info("client connected", "sessionID", session.SessionID())
	})
	hooks.AddOnUnregisterSession(func(ctx context.Context, session server.ClientSession) {
		slog.Info("client disconnected", "sessionID", session.SessionID())
	})
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		slog.Info("tool call", "tool", message.Params.Name)
	})
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result any) {
		slog.Info("tool result", "tool", message.Params.Name)
	})

	// Create a new fast MCP server, with logger enabled
	s := server.NewMCPServer("Markdown Manager", "1.0.0",
		server.WithLogger(logger),
		server.WithHooks(hooks),
	)

	// 1. Tool to list all files
	listTool := mcp.NewTool("list_markdown_files",
		mcp.WithDescription("Lists all the markdown files available in the workspace."),
	)
	s.AddTool(listTool, listFilesHandler)

	// 2. Tool to read a file
	readTool := mcp.NewTool("read_markdown_file",
		mcp.WithDescription("Reads and returns the complete text content of a specific markdown file."),
		mcp.WithString("filename", mcp.Required(), mcp.Description("The name of the file to read (e.g., 'notes.md')")),
	)
	s.AddTool(readTool, readFileHandler)

	// 3. Tool to write a file
	writeTool := mcp.NewTool("write_markdown_file",
		mcp.WithDescription("Creates a new markdown file or overwrites an existing one with new content."),
		mcp.WithString("filename", mcp.Required(), mcp.Description("The name of the file to create/update")),
		mcp.WithString("content", mcp.Required(), mcp.Description("The full markdown content to write")),
	)
	s.AddTool(writeTool, writeFileHandler)

	// 4. Tool to get the current date and time
	timeTool := mcp.NewTool("get_current_datetime",
		mcp.WithDescription("Returns the current date and time in ISO 8601 format (UTC) and the local timezone."),
	)
	s.AddTool(timeTool, getDateTimeHandler)

	slog.Info("starting MCP server", "name", "mcp-md", "version", "0.1.0")

	// Start the server using standard input/output (stdio)
	if err := server.ServeStdio(s); err != nil {
		slog.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

// --- Tool Handlers ---

func listFilesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	files, err := os.ReadDir(notesDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read directory: %v", err)), nil
	}

	var filenames []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			filenames = append(filenames, f.Name())
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("Files found: %v", filenames)), nil
}

func readFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filename, err := request.RequireString("filename")
	if err != nil {
		return mcp.NewToolResultError("missing required argument: filename"), nil
	}
	
	path, err := getSafePath(filename)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("File not found or unreadable: %v", err)), nil
	}

	return mcp.NewToolResultText(string(content)), nil
}

func writeFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filename, err := request.RequireString("filename")
	if err != nil {
		return mcp.NewToolResultError("missing required argument: filename"), nil
	}
	content, err := request.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("missing required argument: content"), nil
	}

	path, err := getSafePath(filename)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully wrote to %s", filename)), nil
}

func getDateTimeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	now := time.Now()
	localZone, _ := now.Zone()
	return mcp.NewToolResultText(fmt.Sprintf("UTC:   %s\nLocal: %s (%s)", now.UTC().Format(time.RFC3339), now.Format(time.RFC3339), localZone)), nil
}

// Helper to keep the agent locked in the notes directory
func getSafePath(filename string) (string, error) {
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}

	// Join and clean the path to prevent directory traversal tricks (like "../../etc/passwd")
	targetPath := filepath.Clean(filepath.Join(notesDir, filename))
	
	// Ensure the absolute path starts with the notes folder path
	absNotesDir, _ := filepath.Abs(notesDir)
	absTarget, _ := filepath.Abs(targetPath)

	if !strings.HasPrefix(absTarget, absNotesDir) {
		return "", fmt.Errorf("security violation: access denied outside allowed directory")
	}

	return absTarget, nil
}