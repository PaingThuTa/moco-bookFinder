# Environment Variables

<!-- AUTO-GENERATED -->

All configuration is loaded via environment variables at startup. See `internal/config/config.go`.

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `TELEGRAM_BOT_TOKEN` | Yes | Telegram bot API token from @BotFather | `123456789:ABCdefGHIjklMNOpqrsTUVwxyz` |
| `ALLOWED_USER_IDS` | Yes | Comma-separated list of authorized Telegram user IDs | `123456789,987654321` |

## Notes

- Both variables are mandatory. The bot will fail to start if either is missing.
- `ALLOWED_USER_IDS` values must be valid 64-bit integers.
- Unauthorized users receive a "You're not authorized" message and the bot ignores their requests.
- To add/remove users, update `ALLOWED_USER_IDS` and restart the process.
