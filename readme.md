# mcp server for markdown 

localy run `mcp server` which allows agentic clients to work with markdown files on the local machines. focused on providing agents with access to strictly defined directories / workspaces. its a single binary server which works over `mcp stdio`.

## usage

if you are interested in control, security and restricting what agentic clients do on your machine, the good practice is usually to disable the inbuilt file and system management tools in these agentic clients.

so you can decide and configure what actions you allow for agentic clients running on your machine.

next, you provide them with multiple mcp tools to work with your filesystem, internet access, database access etc. this way, you control, log and manage how agentic applications are using your machine and OS.

this project, `md-mcp` is one of such fast local mcp tools - dedicated to provide access for agents to markdown file management, only within specifically allowed workspace directory.

point your agentic clients (opencode, cline, 5ire, etc) within their mcp configurations to the `mcp-md` binary. agents will use mcp `stdio` to start and interact with it. `mcp-md` will become available in the tools section within agentic client.

## dev

```bash
make setup
make lint
make build
## run
./bin/mcp-md

make clean
```
