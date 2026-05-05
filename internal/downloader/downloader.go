package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"book-finder/internal/bot"
)

const maxFileSize = 50 * 1024 * 1024 // 50MB Telegram limit

// DownloadedFile represents a book file ready to send via Telegram.
type DownloadedFile struct {
	Data     []byte
	Filename string
	FileType string // "epub" or "pdf"
	Size     int64
}

// FileDownloader resolves a detail page URL to an actual downloadable file.
type FileDownloader interface {
	SourceName() string
	DownloadFile(ctx context.Context, detailURL string) (*DownloadedFile, error)
}

// SourceManager holds per-source downloaders and orchestrates file downloads.
type SourceManager struct {
	downloaders map[string]FileDownloader
	client      *http.Client
}

// NewSourceManager creates a download manager from per-source implementations.
func NewSourceManager(downloaders []FileDownloader, client *http.Client) *SourceManager {
	dlMap := make(map[string]FileDownloader, len(downloaders))
	for _, d := range downloaders {
		dlMap[d.SourceName()] = d
	}
	return &SourceManager{
		downloaders: dlMap,
		client:      client,
	}
}

// DownloadFile looks up the appropriate downloader for the given source
// and fetches the book file, preferring EPUB over PDF.
func (m *SourceManager) DownloadFile(ctx context.Context, sourceName, detailURL string) (*DownloadedFile, error) {
	dl, ok := m.downloaders[sourceName]
	if !ok {
		return nil, fmt.Errorf("no downloader for source %q", sourceName)
	}

	bot.SleepWithDelay()
	return dl.DownloadFile(ctx, detailURL)
}

// FetchFile performs an authenticated HTTP GET with retry logic for rate-limited responses.
func FetchFile(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			bot.SleepWithDelay()
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limited (status %d)", resp.StatusCode)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP error %d from %s", resp.StatusCode, url)
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(io.LimitReader(resp.Body, maxFileSize+1))
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if int64(len(data)) > maxFileSize {
			return nil, ErrFileTooLarge
		}

		if bot.IsCloudflareChallenge(string(data)) {
			return nil, ErrCloudflareBlocked
		}

		return data, nil
	}

	return nil, fmt.Errorf("download failed after retries: %w", lastErr)
}
