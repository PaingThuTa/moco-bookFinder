package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"book-finder/internal/bot"
	"book-finder/internal/config"
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

	// Shared HTTP client for scrapers
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Set up scrapers with fallback
	scrapers := []source.Scraper{
		source.NewZLibraryScraper(httpClient),
		source.NewOceanPDFScraper(httpClient),
		source.NewLibGenScraper(httpClient),
	}
	mgr := source.NewSourceManager(scrapers)

	// Set up handler
	handler := bot.NewHandler(cfg, mgr)

	// Configure updates polling
	ctx := context.Background()
	updates, err := botClient.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}

	// Process updates
	for update := range updates {
		if update.Message != nil {
			msg := update.Message
			// Check if message is a command
			if strings.HasPrefix(msg.Text, "/") {
				cmd, _, _ := tu.ParseCommand(msg.Text)
				switch cmd {
				case "start":
					handler.HandleStart(botClient, update)
				case "search":
					handler.HandleSearch(botClient, update)
				default:
					botClient.SendMessage(ctx, &telego.SendMessageParams{
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
}
