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

While CDNS is designed to work across various operating systems, here is the current status:

| Operating System | Distribution    | Status      |
| :--------------- | :-------------- | :---------- |
| **Linux**        | **Ubuntu**      | ✅ Verified |
| **Linux**        | **Debian**      | ⚠️ Untested |
| **Linux**        | **Fedora**      | ⚠️ Untested |
| **Linux**        | **Arch Linux**  | ⚠️ Untested |
| **Linux**        | **Manjaro**     | ⚠️ Untested |
| **Linux**        | **Pop!\_OS**    | ⚠️ Untested |
| **Linux**        | **Linux Mint**  | ⚠️ Untested |
| **Linux**        | **openSUSE**    | ⚠️ Untested |
| **Linux**        | **NixOS**       | ⚠️ Untested |
| **Linux**        | **CentOS**      | ⚠️ Untested |
| **Linux**        | **Kali Linux**  | ⚠️ Untested |
| **macOS**        | **Darwin**      | 🏗️ WIP      |
| **Windows**      | **Windows 10+** | 🏗️ WIP      |

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

CDNS works in two ways: through a friendly **Interactive TUI** (best for discovery) or via **Quick Commands** (best for automation and power users).

### 🖥️ Interactive Mode (Recommended)

Just run `cdns` without any arguments to open the main menu. From here, you can navigate through all features using your arrow keys.

```bash
cdns
```

### ⚡ Quick Commands

For those who prefer the speed of the command line, CDNS provides a set of intuitive subcommands.

#### 1. Set your DNS
The `set` command is the heart of CDNS. You can use it with a preset name, custom IPs, or even target specific interfaces.

> **Note:** Changing system DNS settings typically requires `sudo` privileges.

```bash
# Apply a preset (e.g., Cloudflare, Google, Quad9)
sudo cdns set cloudflare

# Use custom IP addresses
sudo cdns set 1.1.1.1 8.8.8.8

# Target a specific network interface
sudo cdns set google --interface eth0
```

**Helpful Flags for `set`:**
- `--dry-run`: See what would happen without making any actual changes.
- `--interface` or `-i`: Manually specify which interfaces to modify.
- `--yes`: Skip confirmation prompts (perfect for scripts).

#### 2. Explore Presets
Not sure which provider to use? List all available presets to see names and IP addresses.

```bash
cdns list
```

#### 3. Check Current Status
Verify your active DNS configuration and see which backend (NetworkManager, systemd-resolved, etc.) is being used.

```bash
cdns status

# Pro tip: Use --json for machine-readable output
cdns status --json
```

#### 4. Instant Reset
If you need to roll back to your previous configuration, the `reset` command has your back.

```bash
sudo cdns reset
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

See [LICENSE](./LICENSE) for details.
