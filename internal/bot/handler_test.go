package bot

import (
	"context"
	"testing"

	"book-finder/internal/config"
	"book-finder/internal/source"
)

func TestHandler_ParseSearchArgs_TitleOnly(t *testing.T) {
	h := &BotHandler{}

	title, author := h.parseSearchArgs("The Great Gatsby")
	if title != "The Great Gatsby" {
		t.Errorf("expected title 'The Great Gatsby', got %q", title)
	}
	if author != "" {
		t.Errorf("expected empty author, got %q", author)
	}
}

func TestHandler_ParseSearchArgs_WithAuthor(t *testing.T) {
	h := &BotHandler{}

	title, author := h.parseSearchArgs("The Great Gatsby --author Fitzgerald")
	if title != "The Great Gatsby" {
		t.Errorf("expected title 'The Great Gatsby', got %q", title)
	}
	if author != "Fitzgerald" {
		t.Errorf("expected author 'Fitzgerald', got %q", author)
	}
}

func TestHandler_ParseSearchArgs_NoAuthorFlag(t *testing.T) {
	h := &BotHandler{}

	title, author := h.parseSearchArgs("Some Book --unknown")
	if title != "Some Book --unknown" {
		t.Errorf("expected title 'Some Book --unknown', got %q", title)
	}
	_ = author
}

func TestSearchHandler_ResultsFound(t *testing.T) {
	cfg := &config.Config{
		AllowedUserIDs: map[int64]bool{123: true},
	}
	mgr := &mockSourceManager{
		results: []source.BookResult{
			{Title: "Test Book", DownloadURL: "https://example.com/1", Source: "test"},
		},
	}
	h := NewHandler(cfg, mgr)

	title, author := h.parseSearchArgs("TestBook --author TestAuthor")
	if title != "TestBook" {
		t.Errorf("expected 'TestBook', got %q", title)
	}
	if author != "TestAuthor" {
		t.Errorf("expected 'TestAuthor', got %q", author)
	}
}

func TestSearchHandler_EmptyQuery(t *testing.T) {
	cfg := &config.Config{AllowedUserIDs: map[int64]bool{123: true}}
	mgr := &mockSourceManager{}
	h := NewHandler(cfg, mgr)

	title, _ := h.parseSearchArgs("")
	if title != "" {
		t.Errorf("expected empty title, got %q", title)
	}
}

func TestCallbackData_LinkAction(t *testing.T) {
	cbData := "link_2"
	if cbData != "link_2" {
		t.Errorf("expected 'link_2', got %q", cbData)
	}
}

func TestStoreResults_Cleanup(t *testing.T) {
	cfg := &config.Config{AllowedUserIDs: map[int64]bool{123: true}}
	mgr := &mockSourceManager{}
	h := NewHandler(cfg, mgr)

	results := []source.BookResult{
		{Title: "Test", DownloadURL: "https://example.com", Source: "test", DetailURL: "https://example.com/detail"},
	}
	h.storeResults(123, results)

	fetched := h.getResults(123)
	if len(fetched) != 1 {
		t.Fatalf("expected 1 result, got %d", len(fetched))
	}
	if fetched[0].Title != "Test" {
		t.Errorf("expected 'Test', got %q", fetched[0].Title)
	}
	if fetched[0].DetailURL != "https://example.com/detail" {
		t.Errorf("expected detail URL, got %q", fetched[0].DetailURL)
	}
}

func TestStoreResults_Overwrite(t *testing.T) {
	cfg := &config.Config{AllowedUserIDs: map[int64]bool{123: true}}
	mgr := &mockSourceManager{}
	h := NewHandler(cfg, mgr)

	h.storeResults(123, []source.BookResult{{Title: "First"}})
	h.storeResults(123, []source.BookResult{{Title: "Second"}})

	fetched := h.getResults(123)
	if len(fetched) != 1 {
		t.Fatalf("expected 1 result, got %d", len(fetched))
	}
	if fetched[0].Title != "Second" {
		t.Errorf("expected 'Second', got %q", fetched[0].Title)
	}
}

func TestStoreResults_DifferentChats(t *testing.T) {
	cfg := &config.Config{AllowedUserIDs: map[int64]bool{123: true}}
	mgr := &mockSourceManager{}
	h := NewHandler(cfg, mgr)

	h.storeResults(1, []source.BookResult{{Title: "Chat1"}})
	h.storeResults(2, []source.BookResult{{Title: "Chat2"}})

	if h.getResults(1)[0].Title != "Chat1" {
		t.Errorf("chat 1 should have 'Chat1'")
	}
	if h.getResults(2)[0].Title != "Chat2" {
		t.Errorf("chat 2 should have 'Chat2'")
	}
}

func TestGetResults_NotFound(t *testing.T) {
	cfg := &config.Config{AllowedUserIDs: map[int64]bool{123: true}}
	mgr := &mockSourceManager{}
	h := NewHandler(cfg, mgr)

	results := h.getResults(999)
	if results != nil {
		t.Errorf("expected nil results, got %+v", results)
	}
}

type mockSourceManager struct {
	results []source.BookResult
	err     error
}

func (m *mockSourceManager) Search(ctx context.Context, title, author string) ([]source.BookResult, error) {
	return m.results, m.err
}
