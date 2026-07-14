package tools

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerMarkdownTools defines and registers all markdown-file-related tools
// with the MCP server. These tools let an AI agent list, read, and write
// markdown files within the permitted notes directory.
func registerMarkdownTools(s *server.MCPServer, cfg Config) {
	defs := []toolDef{
		{
			tool: mcp.NewTool("list_markdown_files",
				mcp.WithDescription("Lists all the markdown files available in the workspace."),
			),
			handler: makeListFilesHandler(cfg),
		},
		{
			tool: mcp.NewTool("read_markdown_file",
				mcp.WithDescription("Reads and returns the complete text content of a specific markdown file."),
				mcp.WithString("filename", mcp.Required(), mcp.Description("The name of the file to read (e.g., 'notes.md')")),
			),
			handler: makeReadFileHandler(cfg),
		},
		{
			tool: mcp.NewTool("write_markdown_file",
				mcp.WithDescription("Creates a new markdown file or overwrites an existing one with new content."),
				mcp.WithString("filename", mcp.Required(), mcp.Description("The name of the file to create/update")),
				mcp.WithString("content", mcp.Required(), mcp.Description("The full markdown content to write")),
			),
			handler: makeWriteFileHandler(cfg),
		},
	}

	for _, d := range defs {
		s.AddTool(d.tool, d.handler)
	}
}

// makeListFilesHandler returns a handler that reads the notes directory and
// returns the names of all files ending in ".md". Non-markdown files and
// subdirectories are ignored.
func makeListFilesHandler(cfg Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		files, err := os.ReadDir(cfg.WorkspaceDir)
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
}

// makeReadFileHandler returns a handler that reads and returns the full text
// content of the requested markdown file. The filename is validated and
// sandboxed to the configured notes directory via getSafePath.
func makeReadFileHandler(cfg Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError("missing required argument: filename"), nil
		}

		path, err := getSafePath(cfg.WorkspaceDir, filename)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("File not found or unreadable: %v", err)), nil
		}

		return mcp.NewToolResultText(string(content)), nil
	}
}

// makeWriteFileHandler returns a handler that creates or overwrites a markdown
// file with the provided content. The filename is validated and sandboxed to
// the configured notes directory via getSafePath.
func makeWriteFileHandler(cfg Config) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filename, err := request.RequireString("filename")
		if err != nil {
			return mcp.NewToolResultError("missing required argument: filename"), nil
		}
		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError("missing required argument: content"), nil
		}

		path, err := getSafePath(cfg.WorkspaceDir, filename)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully wrote to %s", filename)), nil
	}
}
