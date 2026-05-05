package source

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type OceanPDFScraper struct {
	client *http.Client
	base   string
}

func NewOceanPDFScraper(client *http.Client) *OceanPDFScraper {
	return &OceanPDFScraper{
		client: client,
		base:   "https://oceanofpdf.com",
	}
}

func (s *OceanPDFScraper) Name() string {
	return "Ocean of PDF"
}

func (s *OceanPDFScraper) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	query := title
	if author != "" {
		query = title + " " + author
	}

	searchURL := fmt.Sprintf("%s/?s=%s", s.base, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ocean of PDF returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Ocean of PDF response body: %w", err)
	}

	if isCloudflareChallenge(string(data)) {
		return nil, fmt.Errorf("Ocean of PDF blocked by Cloudflare")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var results []BookResult

	doc.Find("article, .post, [class*='post']").Each(func(i int, sel *goquery.Selection) {
		titleSel := sel.Find("h2 a, h1 a, .entry-title a").First()
		if titleSel.Length() == 0 {
			return
		}

		bookTitle := strings.TrimSpace(titleSel.Text())
		href, _ := titleSel.Attr("href")

		if bookTitle != "" && href != "" {
			results = append(results, BookResult{
				Title:       bookTitle,
				Author:      "",
				DetailURL:   href,
				DownloadURL: href,
				Source:      "Ocean of PDF",
			})
		}
	})

	return results, nil
}
