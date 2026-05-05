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

type LibGenScraper struct {
	client *http.Client
	base   string
}

func NewLibGenScraper(client *http.Client) *LibGenScraper {
	return &LibGenScraper{
		client: client,
		base:   "https://libgen.is",
	}
}

func (s *LibGenScraper) Name() string {
	return "LibGen"
}

func (s *LibGenScraper) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	params := url.Values{}
	params.Set("req", title)
	if author != "" {
		params.Set("column", "author")
	} else {
		params.Set("column", "title")
	}

	searchURL := fmt.Sprintf("%s/search.php?%s", s.base, params.Encode())

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
		return nil, fmt.Errorf("LibGen returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read LibGen response body: %w", err)
	}

	if isCloudflareChallenge(string(data)) {
		return nil, fmt.Errorf("LibGen blocked by Cloudflare")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var results []BookResult

	doc.Find("table#results tr, table[width='100%'] tr").Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			return
		}

		linkSel := sel.Find("td > a").First()
		if linkSel.Length() == 0 {
			return
		}

		bookTitle := strings.TrimSpace(linkSel.Text())
		href, _ := linkSel.Attr("href")

		if bookTitle == "" {
			return
		}

		author := ""
		sel.Find("td").Each(func(j int, cell *goquery.Selection) {
			if j == 1 {
				author = strings.TrimSpace(cell.Text())
			}
		})

		if href != "" {
			href = resolveURL(s.base, href)
		}

		results = append(results, BookResult{
			Title:       bookTitle,
			Author:      author,
			DetailURL:   href,
			DownloadURL: href,
			Source:      "LibGen",
		})
	})

	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}
