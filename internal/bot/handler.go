package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"book-finder/internal/config"
	"book-finder/internal/source"
)

type SourceManagerInterface interface {
	Search(ctx context.Context, title, author string) ([]source.BookResult, error)
}

type storedResults struct {
	results []source.BookResult
	timer   *time.Timer
}

type BotHandler struct {
	cfg     *config.Config
	mgr     SourceManagerInterface
	results map[int64]*storedResults
	mu      sync.RWMutex
}

func NewHandler(cfg *config.Config, mgr SourceManagerInterface) *BotHandler {
	return &BotHandler{
		cfg:     cfg,
		mgr:     mgr,
		results: make(map[int64]*storedResults),
	}
}

func (h *BotHandler) HandleStart(bot *telego.Bot, update telego.Update) {
	if !h.checkAuthorized(bot, update) {
		return
	}

	msg := telego.SendMessageParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Text: "Welcome to Book Finder Bot!\n\n" +
			"Usage:\n" +
			"/search <book name> — search for a book\n" +
			"/search <book name> --author <author> — search by author\n",
	}
	bot.SendMessage(context.Background(), &msg)
}

func (h *BotHandler) HandleSearch(bot *telego.Bot, update telego.Update) {
	if !h.checkAuthorized(bot, update) {
		return
	}

	args := h.extractCommandArgs(update.Message.Text)
	if strings.TrimSpace(args) == "" {
		msg := telego.SendMessageParams{
			ChatID: tu.ID(update.Message.Chat.ID),
			Text:   "Please provide a book name. Usage: /search <book name>",
		}
		bot.SendMessage(context.Background(), &msg)
		return
	}

	title, author := h.parseSearchArgs(args)

	// Send "searching" message
	searchingMsg := telego.SendMessageParams{
		ChatID: tu.ID(update.Message.Chat.ID),
		Text:   "Searching...",
	}
	sent, err := bot.SendMessage(context.Background(), &searchingMsg)
	if err != nil {
		return
	}

	results, err := h.mgr.Search(context.Background(), title, author)
	if err != nil || len(results) == 0 {
		editMsg := telego.EditMessageTextParams{
			ChatID:    tu.ID(update.Message.Chat.ID),
			MessageID: sent.MessageID,
			Text:      "Book not found",
		}
		bot.EditMessageText(context.Background(), &editMsg)
		return
	}

	// Format results with inline keyboard buttons (one per result)
	var textBuilder strings.Builder
	textBuilder.WriteString(fmt.Sprintf("Found %d result(s):\n\n", len(results)))

	var buttons [][]telego.InlineKeyboardButton
	for i, r := range results {
		textBuilder.WriteString(fmt.Sprintf("%d. **%s**", i+1, r.Title))
		if r.Author != "" {
			textBuilder.WriteString(fmt.Sprintf(" by %s", r.Author))
		}
		textBuilder.WriteString(fmt.Sprintf("\n   Source: %s\n\n", r.Source))

		btn := telego.InlineKeyboardButton{
			Text:         fmt.Sprintf("Download %s", r.Source),
			CallbackData: fmt.Sprintf("link_%d", i),
		}
		buttons = append(buttons, []telego.InlineKeyboardButton{btn})
	}

	// Store results in memory for callback lookup
	h.storeResults(update.Message.Chat.ID, results)

	textMsg := telego.SendMessageParams{
		ChatID:      tu.ID(update.Message.Chat.ID),
		Text:        textBuilder.String(),
		ParseMode:   "Markdown",
		ReplyMarkup: tu.InlineKeyboard(buttons...),
	}

	bot.SendMessage(context.Background(), &textMsg)
}

func (h *BotHandler) HandleCallback(bot *telego.Bot, update telego.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.GetChat().ID

	bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	if !strings.HasPrefix(data, "link_") {
		return
	}

	indexStr := strings.TrimPrefix(data, "link_")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return
	}

	results := h.getResults(chatID)
	if results == nil || index < 0 || index >= len(results) {
		bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: tu.ID(chatID),
			Text:   "Download link not found",
		})
		return
	}

	r := results[index]
	bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID:    tu.ID(chatID),
		Text:      fmt.Sprintf("Download link for **%s**:\n%s", r.Title, r.DownloadURL),
		ParseMode: "Markdown",
	})
}

func (h *BotHandler) checkAuthorized(bot *telego.Bot, update telego.Update) bool {
	var userID int64
	if update.Message != nil {
		userID = update.Message.From.ID
	} else if update.CallbackQuery != nil {
		userID = update.CallbackQuery.From.ID
	}

	if !h.cfg.IsAllowed(userID) {
		if update.Message != nil {
			bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID: tu.ID(update.Message.Chat.ID),
				Text:   "You're not authorized",
			})
		}
		return false
	}
	return true
}

func (h *BotHandler) parseSearchArgs(args string) (title string, author string) {
	authorFlag := "--author "
	idx := strings.Index(args, authorFlag)
	if idx == -1 {
		return strings.TrimSpace(args), ""
	}

	title = strings.TrimSpace(args[:idx])
	author = strings.TrimSpace(args[idx+len(authorFlag):])
	return title, author
}

// extractCommandArgs extracts the arguments portion from a command message.
// e.g. "/search Clean Code --author Martin" -> "Clean Code --author Martin"
func (h *BotHandler) extractCommandArgs(text string) string {
	_, _, args := tu.ParseCommand(text)
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, " ")
}

func (h *BotHandler) storeResults(chatID int64, results []source.BookResult) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, exists := h.results[chatID]; exists {
		old.timer.Stop()
	}

	timer := time.AfterFunc(10*time.Minute, func() {
		h.mu.Lock()
		delete(h.results, chatID)
		h.mu.Unlock()
	})

	h.results[chatID] = &storedResults{
		results: results,
		timer:   timer,
	}
}

func (h *BotHandler) getResults(chatID int64) []source.BookResult {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stored, exists := h.results[chatID]
	if !exists {
		return nil
	}
	return stored.results
}
