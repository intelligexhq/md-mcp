# md-mcp - markdown file mcp server

lean, localy run mcp server for agentic clients to work with markdown files on local machines. within strictly limited directories / workspaces. server works over `stdio`.

## usage

if you fancy little bit of control and security while using agentic clients and applications on your machine, the good practice is usually to first disable the inbuilt file or system management tools in these agentic clients.

so you can decide what do you allow them to do on your devices.

next, you provide them with multiple mcp tools to work with your filesystem, os etc. this way, you can control, log and manage how agentic applications are behaving and what are they allowed to do.

this project, `md-mcp` is one of such very fast local mcp tools - dedicated to provide access to markdown file management, within specifically allowed workspace.

point your agentic clients (5ire, etc) within their mcp configurations to the `mcp-md` binary you build. agents will run it and use mcp `stdio` to interact with it. it will become available in the tools section.

## dev

```bash
make setup
make lint
make build
## run
./bin/mcp-md

make clean
```
