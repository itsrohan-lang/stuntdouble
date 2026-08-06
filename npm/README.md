# stuntdouble-sandbox-cli

Secure AI Agent Sandbox with Docker isolation and git-based rollback.

## Installation

```bash
npm install -g stuntdouble-sandbox-cli
```

The package automatically downloads the appropriate native binary for your platform (Linux, macOS, Windows) during installation.

## Usage

After installation, the `sd` and `stuntdouble` commands are available:

```bash
# Run an agent in a sandboxed container
sd run <agent-name>

# Create a workspace checkpoint
sd checkpoint save <name>

# Restore a checkpoint
sd checkpoint restore <name>

# Show help
sd --help
```

## What it does

StuntDouble runs autonomous coding agents inside a restricted Docker container with:
- Dropped capabilities and resource limits
- Scoped workspace mount
- Git-based snapshot rollback

**Network egress filtering is NOT implemented.** Sandboxed agents can still reach the network. See [ENFORCEMENT.md](https://github.com/itsrohan-lang/stuntdouble/blob/main/docs/ENFORCEMENT.md) for what is and is not enforced.

## Platform support

Pre-built binaries are available for:
- Linux (x64, ARM64)
- macOS (x64, ARM64)
- Windows (x64)

If your platform isn't supported, [build from source](https://github.com/itsrohan-lang/stuntdouble).

## Documentation

Full documentation, deployment guides, and enforcement details:

**https://github.com/itsrohan-lang/stuntdouble**

## License

MIT
