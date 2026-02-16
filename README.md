<div align="center">
<img src="./assets/icon.svg" alt="CDNS" width="120" height="120" />

**change DNS servers effortlessly via terminal**

[**Usage**](#-usage) | [**Report Bugs**](https://github.com/junevm/cdns/issues) | [**Releases**](https://github.com/junevm/cdns/releases) | [**Contributing**](#-contributing)

</div>

**CDNS (change DNS)** is a dead-simple terminal tool that handles the messy details of `systemd-resolved` and `NetworkManager` for you, so you can swap DNS providers in seconds without the headache.

```
  ___ ___  _  _ ___
 / __|   \| \| / __|
| (__| |) | .  \__ \
 \___|___/|_|\_|___/
```

## Why CDNS?

- **🔐 Privacy in a click**: Easily switch to trusted providers like Quad9, Cloudflare, or AdGuard for a more secure browsing experience.
- **✨ Terminal-first**: A clean, reactive TUI that makes managing network settings actually enjoyable.
- **🧠 Zero-config discovery**: It just works. Whether you're on NetworkManager, systemd-resolved, or a plain old resolv.conf, CDNS finds it and handles the heavy lifting.
- **🚑 Fail-safe**: Messed something up? Roll back to your previous configuration instantly with zero stress.

## Compatibility

While CDNS is designed to work across any Linux distribution using standard networking stack, here is the current status:

| Distribution   | Status      |
| :------------- | :---------- |
| **Ubuntu**     | ✅ Verified |
| **Debian**     | ⚠️ Untested |
| **Fedora**     | ⚠️ Untested |
| **Arch Linux** | ⚠️ Untested |
| **Pop!\_OS**   | ⚠️ Untested |

_If it works for you on an untested distro, please [let us know](https://github.com/junevm/cdns/issues)!_

## Installation

### Option 1: Install Script (Recommended)

The easiest way to install the latest release is via our installer script:

```bash
curl -sfL https://raw.githubusercontent.com/junevm/cdns/main/install.sh | sh
```

### Option 2: Go Install

If you have Go installed:

```bash
go install github.com/junevm/cdns/apps/cli/cmd/app@latest
```

### Option 3: Homebrew (Linux)

```bash
brew tap junevm/homebrew-tap
brew install cdns
```

### Option 4: Manual Download

Download the latest binary for your architecture from the [Releases](https://github.com/junevm/cdns/releases) page.

## Usage

**Interactive Mode (Recommended)**

```bash
cdns
```

**Quick Commands**

```bash
# List available DNS presets
cdns list

# Set DNS to Cloudflare
sudo cdns set cloudflare

# Check current status
cdns status

# Restore previous settings
sudo cdns reset
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

See [LICENSE](./LICENSE) for details.
