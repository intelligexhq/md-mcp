# md-mcp - server for managing markdown files

lean mcp server for agent clients to work with markdown files within strictly limited directory / workspace on local machine. works over `stdio`.

## usage

point your agentic clients (like Open Code, Claude Code , 5ire etc) and their mcp configurations to the `mcp-md` binary. it will become available in the tools section.

## dev

```bash
make setup
make lint
make build
## run
./bin/mcp-md

make clean
```
