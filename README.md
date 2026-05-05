# Book Finder

Telegram bot that searches multiple book sources (Z-Library, Ocean of PDF, LibGen) and returns download links as inline keyboard buttons.

## Quick Start

1. Set environment variables:

```bash
export TELEGRAM_BOT_TOKEN="your-bot-token"
export ALLOWED_USER_IDS="123456789,987654321"
```

2. Install dependencies and run:

```bash
go mod download
go run ./cmd/bot
```

## Commands

- `/start` — Show welcome message and usage instructions
- `/search <book name>` — Search for a book across all sources
- `/search <book name> --author <author>` — Search by title and author

## Sources

Searches sequentially with fallback (first source with results wins):
1. Z-Library (`z-library.bz`)
2. Ocean of PDF (`oceanofpdf.com`)
3. LibGen (`libgen.is`)

## Project Structure

```
cmd/bot/            Entry point and long-polling loop
internal/bot/       Telegram bot handler (commands, callbacks, auth)
internal/config/    Environment variable loading and user whitelist
internal/source/    Scraper interface + Z-Library, OceanPDF, LibGen implementations
```

## Documentation

- [Contributing Guide](docs/CONTRIBUTING.md) — setup, testing, code style
- [Environment Variables](docs/ENV.md) — configuration reference
- [Architecture Codemap](docs/CODEMAPS/architecture.md) — system design and data flow
- [Module Codemap](docs/CODEMAPS/modules.md) — package-level API reference
- [Deploy Guide](DEPLOY.md) — deployment to Railway, Fly.io, Render

## Testing

```bash
go test ./...          # Run all unit tests
go test -race ./...    # Run with race detection
go test -cover ./...   # Run with coverage report
```
