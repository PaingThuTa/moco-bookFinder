# Contributing Guide

## Prerequisites

- Go 1.25.7 or later
- A Telegram bot token (from @BotFather)
- Your Telegram user ID (from @userinfobot)

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/book-finder.git
   cd book-finder
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set environment variables:
   ```bash
   export TELEGRAM_BOT_TOKEN="your-bot-token"
   export ALLOWED_USER_IDS="your-user-id"
   ```

4. Run the bot:
   ```bash
   go run ./cmd/bot
   ```

## Available Commands

| Command | Description |
|---------|-------------|
| `go run ./cmd/bot` | Start the bot with long polling |
| `go build -o bot ./cmd/bot/` | Build the production binary |
| `go test ./...` | Run all unit tests |
| `go test -race ./...` | Run tests with race detection |
| `go test -cover ./...` | Run tests with coverage report |
| `go mod tidy` | Clean up module dependencies |
| `go vet ./...` | Run Go static analysis |

## Testing

Tests use the standard `go test` framework with table-driven patterns.

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...

# Run a single package's tests
go test ./internal/bot/

# Run a specific test
go test ./internal/source/ -run TestSourceManager_ReturnsFirstSuccessfulResult
```

Test files:
- `internal/config/config_test.go` — config loading and authorization
- `internal/bot/handler_test.go` — search argument parsing
- `internal/source/source_test.go` — source manager fallback logic

## Code Style

- Format with `gofmt` and `goimports`
- Follow idiomatic Go: accept interfaces, return structs
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Keep interfaces small (1-3 methods)
- No emojis in code, comments, or commit messages

## Git Workflow

Use conventional commits:

```
<type>: <description>
```

Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `build`, `ci`

## PR Submission Checklist

- [ ] Tests pass: `go test -race ./...`
- [ ] Code formatted: `gofmt -l .` (no output expected)
- [ ] No hardcoded secrets
- [ ] Commit message follows conventional commit format
- [ ] Changes are minimal and focused
