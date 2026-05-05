# Book Finder

Telegram bot that searches multiple book sources (Z-Library, Ocean of PDF, LibGen) and returns download links as inline keyboard buttons.

## Setup

1. Configure environment variables:

```
TELEGRAM_BOT_TOKEN=your-bot-token
ALLOWED_USER_IDS=123456789,987654321
```

2. Build and run:

```bash
go run ./cmd/bot
```

## Commands

- `/start` — Show welcome message and usage instructions
- `/search <book name>` — Search for a book across all sources
- `/search <book name> --author <author>` — Search by title and author

## Sources

Searches sequentially with fallback:
1. Z-Library
2. Ocean of PDF
3. LibGen

## Project Structure

```
cmd/bot/            Entry point
internal/bot/       Telegram bot handlers
internal/config/    Environment variable configuration
internal/source/    Book source scrapers (Scraper interface + implementations)
```

## Testing

```bash
go test ./...          # Run all unit tests
go test -race ./...    # Run with race detection
```
