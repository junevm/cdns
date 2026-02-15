---
sidebar_position: 1
---

# Introduction

CDNS is a trusted, Linux-first DNS management CLI tool.

## Features

- **Preset Management**: Quickly switch between trusted DNS providers like Cloudflare, Google, Quad9, and more.
- **Custom Presets**: Define your own DNS servers for work or personal use.
- **Interface Awareness**: Automatically detects active network interfaces.
- **Status Monitoring**: View current DNS settings and connection status in a formatted table.
- **Interactive Mode**: Easy-to-use TUI for selecting options.
- **Linux First**: Optimized for Linux systems using NetworkManager.

## Installation

```bash
# Clone the repository
git clone https://gitlab.com/junevm/cdns.git
cd cdns

# Setup environment
mise install

# Build
mise run cli:build
```

## Usage

```bash
# List available presets
cdns list

# Set DNS to Cloudflare
cdns set cloudflare

# Check current status
cdns status

# Reset to default
cdns reset
```
