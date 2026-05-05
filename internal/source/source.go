package source

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Scraper interface {
	Name() string
	Search(ctx context.Context, title, author string) ([]BookResult, error)
}

type SourceManager struct {
	scrapers []Scraper
	client   *http.Client
}

func NewSourceManager(scrapers []Scraper) *SourceManager {
	return &SourceManager{
		scrapers: scrapers,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *SourceManager) HTTPClient() *http.Client {
	return m.client
}

func (m *SourceManager) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	for _, scraper := range m.scrapers {
		results, err := scraper.Search(ctx, title, author)
		if err == nil && len(results) > 0 {
			return results, nil
		}
		// Try next source
	}
	return nil, fmt.Errorf("all sources failed")
}
