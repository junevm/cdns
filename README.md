<div align="center">
<img src="./assets/icon.svg" alt="CDNS" width="120" height="120" />

# CDNS

**Linux DNS Management Made Simple**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

**CDNS** is a robust, interactive CLI tool for managing DNS settings on Linux. It abstracts away the complexity of modern Linux networking (NetworkManager, systemd-resolved, resolv.conf) to give you simple, safe control over your DNS providers.

## Visuals

```
  ___ ___  _  _ ___
 / __|   \| \| / __|
| (__| |) | .  \__ \
 \___|___/|_|\_|___/
```

## Features

- **🛡️ Privacy Focused**: Quickly switch to encrypted/private DNS providers like Quad9, Cloudflare, or AdGuard.
- **🖥️ Interactive TUI**: Beautiful terminal interface powered by Bubble Tea.
- **🔌 Multi-Backend**: Automatically detects and handles `NetworkManager`, `systemd-resolved`, and traditional `/etc/resolv.conf`.
- **↩️ Safe Reverts**: Automatic state backup allowing quick rollback if things go wrong.

## Installation

### From Source

```bash
go install github.com/junevm/cdns@latest
```

### Pre-built Binaries

Check the [Releases](https://github.com/junevm/cdns/releases) page for your architecture.

## Usage

**Interactive Mode (Recommended)**

```bash
sudo cdns
# or
sudo cdns menu
```

**Quick Commands**

```bash
# List available presets
cdns list

# Set DNS to Cloudflare
sudo cdns set --preset cloudflare

# Check current status
cdns status

# Restore previous settings
sudo cdns reset
```

## Development

This project uses `mise` for task management.

```bash
# Install dependencies
mise install

# Run tests
mise run cli:test

# Run linter
mise run cli:lint
```

## Documentation

Full documentation is available in the [`apps/documentation`](./apps/documentation) directory.

## License

MIT
