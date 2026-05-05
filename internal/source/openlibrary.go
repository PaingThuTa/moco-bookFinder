package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OpenLibraryScraper struct {
	client *http.Client
}

func NewOpenLibraryScraper(client *http.Client) *OpenLibraryScraper {
	return &OpenLibraryScraper{client: client}
}

func (s *OpenLibraryScraper) Name() string {
	return "Open Library"
}

// olSearchResponse is the Open Library search API response format.
type olSearchResponse struct {
	NumFound int `json:"numFound"`
	Docs     []struct {
		Title         string   `json:"title"`
		AuthorName    []string `json:"author_name,omitempty"`
		ISBN          []string `json:"isbn,omitempty"`
		Key           string   `json:"key"`
		Publisher     []string `json:"publisher,omitempty"`
		PublishDate   string   `json:"publish_date,omitempty"`
		NumberOfPages int      `json:"number_of_pages,omitempty"`
		Format        []string `json:"format,omitempty"`
	} `json:"docs"`
}

func (s *OpenLibraryScraper) Search(ctx context.Context, title, author string) ([]BookResult, error) {
	query := title
	if author != "" {
		query = title + " " + author
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("fields", "title,author_name,isbn,key,publisher,publish_date,number_of_pages,format")
	params.Set("limit", "10")

	searchURL := fmt.Sprintf("https://openlibrary.org/search.json?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Open Library returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Open Library response body: %w", err)
	}

	if isCloudflareChallenge(string(data)) {
		return nil, fmt.Errorf("Open Library blocked by Cloudflare")
	}

	var searchResp olSearchResponse
	if err := json.Unmarshal(data, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse Open Library response: %w", err)
	}

	if searchResp.NumFound == 0 || len(searchResp.Docs) == 0 {
		return nil, nil
	}

	var results []BookResult
	for _, doc := range searchResp.Docs {
		detailURL := "https://openlibrary.org" + doc.Key

		var bookAuthor string
		if len(doc.AuthorName) > 0 {
			bookAuthor = strings.Join(doc.AuthorName, ", ")
		}

		var isbn string
		if len(doc.ISBN) > 0 {
			isbn = doc.ISBN[0]
		}

		desc := ""
		if doc.Publisher != nil && len(doc.Publisher) > 0 {
			desc = doc.Publisher[0]
			if doc.PublishDate != "" {
				desc += " (" + doc.PublishDate + ")"
			}
		}

		results = append(results, BookResult{
			Title:       doc.Title,
			Author:      bookAuthor,
			DetailURL:   detailURL,
			DownloadURL: detailURL,
			Source:      "Open Library",
		})

		_ = isbn
		_ = desc

		if len(results) >= 10 {
			break
		}
	}

	return results, nil
}
