# Agent Governance

## 🚨 Critical Instruction

**Strictly follow `./skills/planning-with-files/`. This is non-negotiable.**
Before writing code for complex tasks, you MUST start with a `task_plan.md`.

## 📜 Project Overview

**CDNS** is a Go-based CLI tool (Cobra + Bubble Tea) for managing DNS settings on Linux.

- **Root**: Application source code, configuration, and governance.
- **internal/**: Core logic (DNS backends, models, features).
- **main.go**: Application entry point.

## ⚡ Task Execution Rules

1. **Scope**: Modify only what is requested.
2. **Verification**: Always verify changes.
   - Run: `mise test`
3. **Naming**: The project is named **CDNS**. Avoid "avtoolz" or "go-cli".

## 🛠️ Development Environment

- **Task Runner**: `mise` handles all dev tasks.
- **Linting**: `golangci-lint` (via `mise lint`).
- **Formatting**: Standard `gofmt`.

## 📝 Documentation Policy

- Code changes affecting user behavior MUST have corresponding documentation updates.
- Keep `AGENTS.md` in subdirectories (`apps/cli/AGENTS.md`) updated if specific rules change.

## 🚫 Constraints

- **No Secrets**: Never output or commit secrets.
- **No Force Push**: Respect git history.
- **No Unplanned Refactoring**: Stick to the task plan.
