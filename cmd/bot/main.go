package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"book-finder/internal/bot"
	"book-finder/internal/config"
	"book-finder/internal/downloader"
	"book-finder/internal/source"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	botClient, err := telego.NewBot(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	me, err := botClient.GetMe(context.Background())
	if err != nil {
		log.Fatalf("Failed to get bot info: %v", err)
	}
	log.Printf("Authorized on account %s", me.Username)

	// Shared HTTP client with browser headers for scraping
	httpClient := downloader.NewHTTPClient(30 * time.Second)

	// Set up scrapers with fallback
	scrapers := []source.Scraper{
		source.NewOpenLibraryScraper(httpClient),
		source.NewZLibraryScraper(httpClient),
		source.NewOceanPDFScraper(httpClient),
		source.NewLibGenScraper(httpClient),
	}
	mgr := source.NewSourceManager(scrapers)

	// Set up handler
	handler := bot.NewHandler(cfg, mgr)

	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL != "" {
		runWebhook(botClient, handler, webhookURL, me.Username)
	} else {
		runPolling(botClient, handler)
	}
}

func runWebhook(botClient *telego.Bot, handler *bot.BotHandler, webhookURL, botName string) {
	ctx := context.Background()
	webhookPath := "/" + botName
	addr := ":" + os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		addr = ":10000"
	}

	// Register webhook on Telegram
	err := botClient.SetWebhook(ctx, &telego.SetWebhookParams{
		URL: webhookURL + webhookPath,
	})
	if err != nil {
		log.Fatalf("Failed to set webhook: %v", err)
	}
	log.Printf("Webhook set to %s%s", webhookURL, webhookPath)

	// Set up update channel from webhook
	updates, err := botClient.UpdatesViaWebhook(
		ctx,
		func(h telego.WebhookHandler) error {
			return telego.WebhookHTTPServeMux(
				http.DefaultServeMux,
				webhookPath,
			)(h)
		},
	)
	if err != nil {
		log.Fatalf("Failed to set up webhook: %v", err)
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Process updates
	go func() {
		for update := range updates {
			processUpdate(botClient, handler, update)
		}
	}()

	log.Printf("Starting webhook server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runPolling(botClient *telego.Bot, handler *bot.BotHandler) {
	ctx := context.Background()

	// Delete any existing webhook so polling can work
	if err := botClient.DeleteWebhook(ctx, &telego.DeleteWebhookParams{}); err != nil {
		log.Printf("Failed to delete webhook: %v", err)
	}

	updates, err := botClient.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}
	log.Println("Running in polling mode")

	for update := range updates {
		processUpdate(botClient, handler, update)
	}
}

func processUpdate(botClient *telego.Bot, handler *bot.BotHandler, update telego.Update) {
	if update.Message != nil {
		msg := update.Message
		if strings.HasPrefix(msg.Text, "/") {
			cmd, _, _ := tu.ParseCommand(msg.Text)
			switch cmd {
			case "start":
				handler.HandleStart(botClient, update)
			case "search":
				handler.HandleSearch(botClient, update)
			default:
				botClient.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID: tu.ID(msg.Chat.ID),
					Text:   "Unknown command. Use /start to see available commands.",
				})
			}
		}
	}

	if update.CallbackQuery != nil {
		handler.HandleCallback(botClient, update)
	}
}
