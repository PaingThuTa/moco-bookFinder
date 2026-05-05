# Deploy & Use Guide — Book Finder Telegram Bot

## Prerequisites

- A Telegram account
- A GitHub account
- A free Render account (no credit card needed)

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
cd book-finder

# Set environment variables
export TELEGRAM_BOT_TOKEN="your-bot-token-from-botfather"
export ALLOWED_USER_IDS="your-user-id"

# Run in polling mode (no WEBHOOK_URL = polling)
go run ./cmd/bot/
```

You should see: `Running in polling mode` and `Authorized on account your_bot_username`

## Step 4: Test Locally

```
/start          → Welcome message
/search clean code  → Search for books
```

Tap **Send File** to get the file directly, or **Get Link** for the download URL.

## Step 5: Push to GitHub

```bash
git remote add origin https://github.com/yourusername/book-finder.git
git push -u origin main
```

## Step 6: Deploy to Render (Free, No Credit Card)

1. Go to [render.com](https://render.com) and sign up (GitHub login)

2. Click **New +** → **Web Service**

3. Connect your GitHub account and select the `book-finder` repo

4. Configure the service:
   - **Name**: `book-finder-bot`
   - **Branch**: `main`
   - **Environment**: Docker
   - **Plan**: Free

5. Add **Environment Variables**:
   ```
   TELEGRAM_BOT_TOKEN=your-bot-token-from-botfather
   ALLOWED_USER_IDS=your-user-id
   WEBHOOK_URL=https://your-service-name.onrender.com
   PORT=10000
   ```
   Replace `your-service-name` with the Render subdomain you'll get after creation (you can set this after the first deploy).

6. Click **Create Web Service**

### Step 6a: Set the Webhook URL

After the first deploy, Render gives you a URL like `https://book-finder-bot-abc123.onrender.com`.

Go to the **Environment** tab on Render and set `WEBHOOK_URL` to that URL. The service will redeploy.

> **Important**: The `WEBHOOK_URL` must exactly match the Render URL (with trailing slash removed). The bot will automatically register the Telegram webhook endpoint at `/<bot-username>`.

## How It Works

| Mode | Trigger | When to Use |
|------|---------|-------------|
| **Polling** | Bot polls Telegram for updates | Local development |
| **Webhook** | Telegram sends HTTP POST to your server | Deployed on Render/any web host |

The bot auto-detects the mode: if `WEBHOOK_URL` is set, it runs in webhook mode. Otherwise, polling.

## Managing Allowed Users

To add or remove users, update `ALLOWED_USER_IDS` on Render's Environment tab. Comma-separated: `123456789,987654321`.

## Render Free Tier Notes

- Free instances spin down after 15 minutes of inactivity
- When a user messages the bot, Telegram sends a webhook POST that wakes up the instance
- First message after idle may take 10-30 seconds while Render boots the instance
- 750 free hours/month (enough for one always-on service)
- 512MB RAM limit — sufficient for this bot (downloaded files are capped at 50MB)

## Troubleshooting

**First message takes 30+ seconds after idle:**
- This is Render's cold start. Normal for free tier. Subsequent messages are instant.

**Bot doesn't respond:**
- Check Render logs for errors
- Verify `TELEGRAM_BOT_TOKEN` is correct
- Verify your user ID is in `ALLOWED_USER_IDS`
- Check `WEBHOOK_URL` matches your Render URL

**Webhook not working:**
- Check `/health` endpoint works: visit `https://your-url.onrender.com/health`
- Check Render logs for "Webhook set to..." message
- You can check current webhook info via Telegram API: `https://api.telegram.org/bot<token>/getWebhookInfo`

**"Source temporarily unavailable":**
- The book site may have Cloudflare protection
- Try another source or use the **Get Link** button
