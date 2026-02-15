# Agent Governance

## 🚨 Critical Instruction

**Strictly follow `./skills/planning-with-files/`. This is non-negotiable.**
Before writing code for complex tasks, you MUST start with a `task_plan.md`.

## 📜 Project Overview

**CDNS** is a monorepo containing a Go CLI tool and its documentation.

- **Root**: Workspace config and governance.
- **apps/cli**: Go-based CLI application (Cobra + Bubble Tea).
- **apps/documentation**: Docusaurus-based documentation site.

## ⚡ Task Execution Rules

1. **Scope**: Modify only what is requested. Isolate changes to the specific app (`cli` or `documentation`) unless it's a cross-cutting concern.
2. **Verification**: Always verify changes.
   - CLI: `mise run cli:test`
   - Docs: `mise run docs:build`
3. **Naming**: The project is named **CDNS**. Avoid "avtoolz" or "go-cli".

## 🛠️ Development Environment

- **Task Runner**: `mise` handles all dev tasks.
- **Linting**: `golangci-lint` (via `mise run cli:lint`).
- **Formatting**: Standard `gofmt`.

## 📝 Documentation Policy

- Code changes affecting user behavior MUST have corresponding documentation updates.
- Keep `AGENTS.md` in subdirectories (`apps/cli/AGENTS.md`) updated if specific rules change.

## 🚫 Constraints

- **No Secrets**: Never output or commit secrets.
- **No Force Push**: Respect git history.
- **No Unplanned Refactoring**: Stick to the task plan.
