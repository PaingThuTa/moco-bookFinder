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

type ZLibraryScraper struct {
	client *http.Client
	base   string
}

func NewZLibraryScraper(client *http.Client) *ZLibraryScraper {
	return &ZLibraryScraper{
		client: client,
		base:   "https://z-library.bz",
	}
}

func (s *ZLibraryScraper) Name() string {
	return "Z-Library"
}

func (s *ZLibraryScraper) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	query := title
	if author != "" {
		query = title + " " + author
	}

	searchURL := fmt.Sprintf("%s/s/%s", s.base, url.PathEscape(query))

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
		return nil, fmt.Errorf("Z-Library returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Z-Library response body: %w", err)
	}

	if isCloudflareChallenge(string(data)) {
		return nil, fmt.Errorf("Z-Library blocked by Cloudflare")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var results []BookResult

	// Z-Library results typically have book items with title links
	// Selector may need adjustment based on current site structure
	doc.Find("div.resItemBox, .book-item, [class*='resItem'], [class*='book']").Each(func(i int, sel *goquery.Selection) {
		titleSel := sel.Find("a[href*='/book/'], h3 a, .title a").First()
		if titleSel.Length() == 0 {
			return
		}

		bookTitle, _ := titleSel.Attr("title")
		if bookTitle == "" {
			bookTitle = strings.TrimSpace(titleSel.Text())
		}

		href, _ := titleSel.Attr("href")
		if href != "" {
			href = resolveURL(s.base, href)
		}

		auth := strings.TrimSpace(sel.Find(".author, a[href*='/author/']").First().Text())

		if bookTitle != "" {
			results = append(results, BookResult{
				Title:       bookTitle,
				Author:      auth,
				DetailURL:   href,
				DownloadURL: href,
				Source:      "Z-Library",
			})
		}
	})

	// Fallback: try any link-based search results
	if len(results) == 0 {
		doc.Find("a[href*='/book/']").Each(func(i int, sel *goquery.Selection) {
			if i >= 10 {
				return
			}
			title := strings.TrimSpace(sel.Text())
			href, _ := sel.Attr("href")
			if title != "" && href != "" {
				results = append(results, BookResult{
					Title:       title,
					Author:      "",
					DetailURL:   resolveURL(s.base, href),
					DownloadURL: resolveURL(s.base, href),
					Source:      "Z-Library",
				})
			}
		})
	}

	return results, nil
}

func resolveURL(base, relative string) string {
	if strings.HasPrefix(relative, "http") {
		return relative
	}
	return strings.TrimSuffix(base, "/") + relative
}
