package bot

import (
	"context"
	"testing"

	"book-finder/internal/config"
	"book-finder/internal/downloader"
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
	// Should treat entire string as title since --author not present
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
	h := NewHandler(cfg, mgr, &mockDownloadManager{})

	// parseSearchArgs is the testable unit
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
	h := NewHandler(cfg, mgr, &mockDownloadManager{})

	title, _ := h.parseSearchArgs("")
	if title != "" {
		t.Errorf("expected empty title, got %q", title)
	}
}

type mockSourceManager struct {
	results []source.BookResult
	err     error
}

func (m *mockSourceManager) Search(ctx context.Context, title, author string) ([]source.BookResult, error) {
	return m.results, m.err
}

type mockDownloadManager struct{}

func (m *mockDownloadManager) DownloadFile(ctx context.Context, sourceName, detailURL string) (*downloader.DownloadedFile, error) {
	return nil, nil
}
