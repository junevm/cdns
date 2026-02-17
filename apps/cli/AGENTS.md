# AGENTS.md: CLI Application

Governance for AI agents working in the Go CLI application.

## Authority

- **Inherits from**: `/AGENTS.md` (root)
- **Extends**: Adds CLI-specific rules
- **Scope**: All files within `apps/cli/`
- **Precedence**: These rules extend root rules; root rules cannot be contradicted

## Mandatory Planning Workflow

**ABSOLUTE REQUIREMENT**: Follow the planning workflow defined in `/.github/planning-with-files/` before any code modification.

This is NON-NEGOTIABLE for:

- Architecture changes
- New features
- Refactoring multiple files
- Complex bug fixes
- Dependency updates

See root `/AGENTS.md` for complete planning requirements.

## CLI Application Architecture

This application follows **Hexagonal Architecture (Ports & Adapters)** with **Dependency Injection** via Uber Fx.

### Core Stack

**Primary Technologies:**

- Go 1.25.5
- Cobra (CLI framework)
- Bubble Tea (TUI framework)
- Uber Fx (dependency injection)
- Koanf (configuration)
- slog (structured logging)

**Build System:**

- Task runner: `mise`
- Linter: `golangci-lint`
- Formatter: `gofmt`
- Hooks: `lefthook`

### Package Structure

```
apps/cli/
├── cmd/cdns/            # Application entry point (Fx wiring only)
├── internal/
│   ├── cli/             # Cobra adapters (primary port)
│   ├── config/          # Configuration infrastructure (Koanf)
│   ├── logger/          # Logging infrastructure (slog)
│   ├── ui/              # UI components (Bubble Tea)
│   └── features/        # Feature modules (business logic)
│       └── <feature>/
│           ├── module.go    # Fx module definition
│           ├── service.go   # Business logic
│           ├── command.go   # Cobra command adapter
│           └── *_test.go    # Tests
├── docs/                # Architecture docs and ADRs
└── config.yaml          # Default configuration
```

## Absolute Rules (Non-Negotiable)

Violating these rules invalidates the contribution.

### 1. No Global State

**Prohibited:**

- Global variables
- Package-level singletons
- `init()` functions for initialization

**Required:**

- All dependencies injected via Uber Fx
- Use Fx lifecycle hooks (`OnStart`, `OnStop`)

### 2. Strict Layer Boundaries

**Cobra Commands** (primary adapter):

- **CAN**: Parse flags, validate arguments, invoke services
- **CANNOT**: Contain business logic, access databases, make network calls, use `os.Exit`

**Services** (application layer):

- **CAN**: Execute business logic, call other services, use repositories
- **CANNOT**: Import `cobra`, `viper`, or any CLI-specific packages

**Infrastructure** (secondary adapters):

- **CAN**: Access databases, filesystems, networks, external APIs
- **MUST**: Implement interfaces defined in application layer

### 3. Mandatory Control Flow

Every operation must follow this exact flow:

```
User Input (flags/args)
    ↓
Cobra Command (adapter)
    ↓
Service (business logic)
    ↓
Domain logic / Infrastructure
    ↓
Output
```

Skipping layers or reversing the flow is forbidden.

### 4. No Silent Errors

**Prohibited:**

```go
result, _ := doSomething()  // Ignoring error
```

**Required:**

```go
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}
```

Exception: Explicitly justified with comment explaining why error is safe to ignore.

### 5. No Direct Output in Services

**Prohibited:**

```go
func (s *Service) DoWork() {
    fmt.Println("Working...")  // Direct output
    log.Println("Done")        // Wrong logger
}
```

**Required:**

```go
func (s *Service) DoWork(ctx context.Context) error {
    s.logger.InfoContext(ctx, "working", slog.String("task", "example"))
    return nil
}
```

Use `log/slog` for all logging. Cobra commands handle user-facing output.

## Feature Module Requirements

### Mandatory Structure

Every feature must be self-contained in `internal/features/<name>/`:

```
internal/features/myfeature/
├── module.go         # Fx module (REQUIRED)
├── service.go        # Business logic (REQUIRED)
├── command.go        # Cobra command (if needed)
├── service_test.go   # Unit tests (REQUIRED)
└── README.md         # Feature documentation (optional)
```

### Fx Module Pattern

**Required pattern:**

```go
// module.go
package myfeature

import "go.uber.org/fx"

var Module = fx.Module("myfeature",
    fx.Provide(NewService),      // Provide service
    fx.Provide(NewCommand),      // Provide command
    fx.Invoke(RegisterCommand),  // Register with root
)
```

All features must export a `Module` variable for inclusion in `main.go`.

### Service Constructor Pattern

**Required pattern:**

```go
// service.go
type Service struct {
    logger *slog.Logger
    config *config.Config
}

func NewService(logger *slog.Logger, cfg *config.Config) *Service {
    return &Service{
        logger: logger,
        config: cfg,
    }
}

func (s *Service) Execute(ctx context.Context) error {
    s.logger.InfoContext(ctx, "executing")
    // Business logic here
    return nil
}
```

Constructors must:

- Accept dependencies as parameters (Fx injection)
- Return pointer to struct
- Be fast (no I/O, no heavy computation)

### Command Registration Pattern

**Required pattern:**

```go
// command.go
type CommandResult struct {
    fx.Out
    Cmd *cobra.Command `name:"myfeature"`
}

func NewCommand(service *Service) CommandResult {
    cmd := &cobra.Command{
        Use:   "myfeature",
        Short: "Description",
        RunE: func(cmd *cobra.Command, args []string) error {
            return service.Execute(cmd.Context())
        },
    }
    return CommandResult{Cmd: cmd}
}

func RegisterCommand(root *cobra.Command, cmd CommandResult) {
    root.AddCommand(cmd.Cmd)
}
```

## Configuration Management

### System

Use `koanf` exclusively via `internal/config`.

### Precedence Order

Configuration sources are merged in this order (later overrides earlier):

1. Defaults (in code)
2. Config file (`config.yaml`)
3. Environment variables
4. CLI flags

### Access Pattern

**Prohibited:**

```go
value := os.Getenv("MY_VAR")           // Direct access
value := viper.GetString("my.key")     // Framework coupling
```

**Required:**

```go
// Service receives typed config via Fx
func NewService(cfg *config.Config) *Service {
    value := cfg.MyValue
    // Use value
}
```

Services receive `*config.Config` struct via dependency injection.

## Logging Requirements

### Library

Use `log/slog` exclusively. No other logging libraries permitted.

### Structured Logging

**Prohibited:**

```go
fmt.Printf("User %s logged in\n", userID)
log.Println("Processing request")
```

**Required:**

```go
logger.InfoContext(ctx, "user logged in",
    slog.String("user_id", userID),
    slog.String("ip", ipAddr),
)
```

### Context Propagation

All service methods must accept `context.Context` as first parameter:

```go
func (s *Service) DoWork(ctx context.Context, param string) error {
    s.logger.InfoContext(ctx, "working", slog.String("param", param))
    return nil
}
```

## Error Handling

### Propagation

Errors must bubble up to the Cobra command layer:

```go
// Service returns error
func (s *Service) DoWork() error {
    if err := s.repository.Fetch(); err != nil {
        return fmt.Errorf("fetching data: %w", err)
    }
    return nil
}

// Command handles error
RunE: func(cmd *cobra.Command, args []string) error {
    if err := service.DoWork(); err != nil {
        return err  // Cobra displays it
    }
    return nil
}
```

### Wrapping

Always wrap errors with context using `fmt.Errorf` and `%w`:

```go
if err != nil {
    return fmt.Errorf("action description: %w", err)
}
```

### Panics

**Prohibited** in runtime code. Only allowed during bootstrap if critical dependency fails:

```go
// Acceptable (bootstrap only)
func main() {
    fx.New(
        // modules
    ).Run()
    // Fx will panic if wiring fails
}
```

## Testing Requirements

### Mandatory Tests

Every feature service must have unit tests in `*_test.go` files.

### Table-Driven Tests

**Required pattern:**

```go
func TestService_Execute(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "test", "result", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Mocking

Mock dependencies using interfaces:

```go
type Repository interface {
    Fetch() (string, error)
}

// In tests
type mockRepository struct{}

func (m *mockRepository) Fetch() (string, error) {
    return "mock data", nil
}
```

### Running Tests

**Required command:**

```bash
mise run cli:test
```

Tests must pass before committing.

## Build and Verification Commands

### Development

```bash
mise run cli:dev        # Run CLI in development mode
```

### Testing

```bash
mise run cli:test       # Run tests with coverage
go test ./...           # Direct test execution
```

### Linting

```bash
mise run cli:lint       # Run linter and formatter
golangci-lint run       # Direct linter execution
```

### Cleanup

```bash
mise run cli:clean      # Remove build artifacts
mise run cli:tidy       # Tidy Go modules
```

## CLI-Specific Conventions

### Command Structure

```go
&cobra.Command{
    Use:   "command [arguments]",
    Short: "One-line description",
    Long:  "Detailed description",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // Parse flags
        // Call service
        // Handle output
        return nil
    },
}
```

### Flag Definition

```go
cmd.Flags().StringP("output", "o", "", "Output format")
cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
```

### Output Conventions

- Use `cmd.OutOrStdout()` for normal output
- Use `cmd.ErrOrStderr()` for error messages
- Support `--quiet` flag for script usage
- Respect `NO_COLOR` environment variable

### Exit Codes

Never call `os.Exit` directly. Return errors:

```go
// Command returns error
RunE: func(cmd *cobra.Command, args []string) error {
    if invalid {
        return fmt.Errorf("invalid input: %s", reason)
    }
    return nil
}

// Cobra handles exit code
```

## CLI Stability Requirements

### Backward Compatibility

Command interfaces are contracts:

- **Flag names**: Cannot change once released
- **Argument positions**: Cannot reorder
- **Output format**: Must remain parseable
- **Exit codes**: Must remain consistent

### Breaking Changes

Require major version bump:

- Removing commands
- Removing flags
- Changing flag types
- Changing default behavior

### Deprecation Process

1. Mark as deprecated in help text
2. Log warning when used
3. Maintain for at least one minor version
4. Remove in next major version

## Documentation Requirements

### Code Documentation

**Required:**

- Package documentation for each package
- Exported functions must have doc comments
- Complex logic must have inline comments

**Format:**

```go
// Package myfeature provides functionality for X.
package myfeature

// Service handles business logic for myfeature.
type Service struct {
    logger *slog.Logger
}

// Execute performs the main operation.
// It returns an error if the operation fails.
func (s *Service) Execute(ctx context.Context) error {
    // Implementation
}
```

### Architecture Decision Records

Significant decisions must be documented in `docs/adr/`:

**When required:**

- Changing architectural patterns
- Adding major dependencies
- Modifying control flow
- Changing configuration strategy

**Format:** Follow standard ADR format (Context, Decision, Consequences)

### Runbook Updates

Update `docs/runbooks/` when adding:

- New features requiring documentation
- New CLI commands
- Configuration options
- Error handling procedures

## Anti-Patterns (Forbidden)

### The "Fat Command"

**Prohibited:**

```go
RunE: func(cmd *cobra.Command, args []string) error {
    // Business logic in command
    db := connectDatabase()
    data := db.Query()
    processData(data)
    return nil
}
```

**Required:** Delegate to service.

### The "Global Singleton"

**Prohibited:**

```go
var DB *sql.DB

func init() {
    DB = connectDatabase()
}
```

**Required:** Use Fx injection.

### The "Framework-Aware Service"

**Prohibited:**

```go
func (s *Service) DoWork() error {
    flag := viper.GetString("flag")  // Coupling
    // ...
}
```

**Required:** Receive config via constructor.

### The "Silent Failure"

**Prohibited:**

```go
result, _ := operation()
```

**Required:** Handle or propagate all errors.

## Knowledge Verification Requirements

Before implementing non-trivial changes:

### Research Declaration

**Required format:**

```markdown
## Sources Consulted

- Go standard library: context package (pkg.go.dev/context)
- Cobra official docs: Command flags
- Uber Fx docs: Lifecycle hooks
- DeepWiki: Go error handling patterns
```

### Required Sources (in priority order)

1. **Official Documentation**
   - Go standard library (pkg.go.dev)
   - Cobra documentation
   - Uber Fx documentation
   - Koanf documentation
   - Bubble Tea documentation

2. **Project Documentation**
   - `apps/cli/docs/`
   - ADRs in `apps/cli/docs/adr/`
   - Runbooks in `apps/cli/docs/runbooks/`

3. **MCP Servers** (if available)
   - DeepWiki for Go best practices
   - DeepWiki for framework-specific patterns

### Prohibited Sources

- Blog posts as primary source
- Stack Overflow as primary source
- Outdated tutorials
- Deprecated documentation

### Verification

Before implementation, verify:

- Concepts fully understood
- APIs are current (not deprecated)
- Edge cases identified
- Testing strategy defined

If uncertain, continue research or request clarification. **Silent guessing is forbidden.**

## Git and Version Control

### Commit Messages

Follow Conventional Commits with CLI scope:

```
feat(cli): add user authentication command
fix(cli): correct flag parsing in list command
refactor(cli): extract validation logic to service
docs(cli): update architecture decision records
```

### Pre-commit Verification

Before committing:

1. Run `mise run cli:lint`
2. Run `mise run cli:test`
3. Verify lefthook hooks pass

### Code Review Checklist

- [ ] Follows hexagonal architecture
- [ ] No global state
- [ ] Services framework-agnostic
- [ ] All errors handled
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Conventional commit format
