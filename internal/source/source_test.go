package source

import (
	"context"
	"errors"
	"testing"
)

type mockScraper struct {
	name    string
	results []BookResult
	err     error
}

func (m *mockScraper) Name() string {
	return m.name
}

func (m *mockScraper) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	return m.results, m.err
}

func TestSourceManager_ReturnsFirstSuccessfulResult(t *testing.T) {
	scrapers := []Scraper{
		&mockScraper{name: "fail", err: errors.New("down")},
		&mockScraper{
			name:    "success",
			results: []BookResult{{Title: "Test Book", Source: "success"}},
		},
	}
	mgr := NewSourceManager(scrapers)

	results, err := mgr.Search(context.Background(), "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Test Book" {
		t.Errorf("expected 'Test Book', got %q", results[0].Title)
	}
}

func TestSourceManager_AllFail(t *testing.T) {
	scrapers := []Scraper{
		&mockScraper{name: "a", err: errors.New("down")},
		&mockScraper{name: "b", err: errors.New("down")},
	}
	mgr := NewSourceManager(scrapers)

	_, err := mgr.Search(context.Background(), "test", "")
	if err == nil {
		t.Fatal("expected error when all sources fail")
	}
}

func TestSourceManager_EmptyResultsSkipped(t *testing.T) {
	scrapers := []Scraper{
		&mockScraper{name: "empty", results: []BookResult{}},
		&mockScraper{
			name:    "has-results",
			results: []BookResult{{Title: "Found", Source: "has-results"}},
		},
	}
	mgr := NewSourceManager(scrapers)

	results, err := mgr.Search(context.Background(), "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Source != "has-results" {
		t.Errorf("expected 'has-results' source, got %q", results[0].Source)
	}
}

func TestFormatBooks(t *testing.T) {
	results := []BookResult{
		{Title: "Book A", Author: "Auth X", Source: "Src"},
	}
	formatted := FormatBooks(results)
	if !contains(formatted, "Book A") {
		t.Error("expected formatted output to contain 'Book A'")
	}
	if !contains(formatted, "Auth X") {
		t.Error("expected formatted output to contain 'Auth X'")
	}
}

func TestFormatBooks_Empty(t *testing.T) {
	formatted := FormatBooks([]BookResult{})
	if formatted != "No results found." {
		t.Errorf("expected 'No results found.', got %q", formatted)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
