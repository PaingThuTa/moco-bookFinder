# Deploy & Use Guide — Book Finder Telegram Bot

## Prerequisites

- Go 1.21+ installed locally
- A Telegram account
- A free Railway, Fly.io, or Render account for hosting

---

## Step 1: Create Your Telegram Bot

1. Open Telegram and search for **@BotFather**
2. Send `/newbot`
3. Follow the prompts:
   - Choose a display name (e.g., "Book Finder")
   - Choose a username (must end in `bot`, e.g., `mybookfinder_bot`)
4. BotFather gives you a **token** — save it. It looks like: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`

## Step 2: Find Your Telegram User ID

1. Open Telegram and search for **@userinfobot**
2. Send `/start`
3. It replies with your numeric user ID (e.g., `123456789`)
4. Save this ID — only whitelisted IDs can use the bot

## Step 3: Local Development

```bash
# Clone and enter the project
cd book-finder

# Set environment variables
export TELEGRAM_BOT_TOKEN="your-bot-token-from-botfather"
export ALLOWED_USER_IDS="your-user-id"

# Run the bot
go run ./cmd/bot/
```

You should see: `Authorized on account your_bot_username`

The bot is now running locally. Test it by sending messages on Telegram.

## Step 4: Test the Bot

In Telegram, message your bot:

```
/start          → Welcome message with usage instructions
/search clean code              → Search for "Clean Code"
/search clean code --author robert martin  → Search with author filter
```

When results appear, tap a "Download" button to get the direct link.

## Step 5: Deploy to Railway (Recommended — Easiest)

1. Push your code to GitHub:
   ```bash
   git remote add origin https://github.com/yourusername/book-finder.git
   git push -u origin main
   ```

2. Go to [railway.app](https://railway.app) and sign in with GitHub

3. Click **New Project** → **Deploy from GitHub repo**

4. Select the `book-finder` repository

5. Go to **Variables** tab and add:
   ```
   TELEGRAM_BOT_TOKEN=your-bot-token
   ALLOWED_USER_IDS=your-user-id
   ```

6. Railway auto-detects Go and builds. No `Procfile` or config needed — the entry point is `cmd/bot/main.go`.

7. Railway may need a `Procfile` to know which binary to run. If so, create one:
   ```
   worker: ./bot
   ```
   And a simple `Dockerfile`:
   ```dockerfile
   FROM golang:1.23-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN go build -o bot ./cmd/bot/

   FROM alpine:latest
   COPY --from=builder /app/bot .
   CMD ["./bot"]
   ```

8. Once deployed, Railway shows a deployment URL — the bot uses long polling, so **no URL configuration needed**. It just works.

## Step 6: Deploy to Fly.io (Alternative)

```bash
# Install fly CLI
brew install flyctl

# Login
fly auth login

# Initialize app
fly launch --name book-finder-bot

# Add secrets
fly secrets set TELEGRAM_BOT_TOKEN="your-token" ALLOWED_USER_IDS="your-id"

# Deploy
fly deploy
```

Fly.io auto-detects Go. If needed, it generates a `Dockerfile` automatically.

## Step 7: Deploy to Render (Alternative)

1. Go to [render.com](https://render.com) and sign in

2. Click **New** → **Web Service**

3. Connect your GitHub repo

4. Settings:
   - **Build Command**: `go build -o bot ./cmd/bot/`
   - **Start Command**: `./bot`
   - **Environment Variables**: Add `TELEGRAM_BOT_TOKEN` and `ALLOWED_USER_IDS`

5. Click **Create Web Service**

---

## Managing Allowed Users

To add or remove users, update the `ALLOWED_USER_IDS` environment variable on your hosting platform:

```bash
# Railway: Variables tab → edit ALLOWED_USER_IDS
# Fly.io:
fly secrets set ALLOWED_USER_IDS="123,456,789"
# Render: Environment tab → edit variable
```

Multiple IDs are comma-separated. No restart needed for most platforms (they reload env vars automatically).

---

## Troubleshooting

**Bot doesn't respond:**
- Check logs on your hosting platform for errors
- Verify `TELEGRAM_BOT_TOKEN` is correct
- Verify your user ID is in `ALLOWED_USER_IDS`

**"All sources busy, try again later":**
- One or more book sites may be temporarily down
- Try again in a few minutes
- The bot tries Z-Library → Ocean of PDF → LibGen in sequence

**Build fails locally:**
- Ensure Go 1.21+: `go version`
- Run `go mod tidy` to fix dependencies

**Rate limited by Telegram:**
- Telegram limits bots to ~30 messages/second
- The bot processes one message at a time, so this shouldn't happen with 2-5 users
