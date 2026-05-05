package downloader

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LibGenDownloader resolves LibGen detail pages to downloadable EPUB/PDF files.
type LibGenDownloader struct {
	client *http.Client
}

// NewLibGenDownloader creates a LibGen file downloader.
func NewLibGenDownloader(client *http.Client) *LibGenDownloader {
	return &LibGenDownloader{client: client}
}

func (d *LibGenDownloader) SourceName() string {
	return "LibGen"
}

func (d *LibGenDownloader) DownloadFile(ctx context.Context, detailURL string) (*DownloadedFile, error) {
	data, err := FetchFile(ctx, d.client, detailURL)
	if err != nil {
		return nil, err
	}

	html := string(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// LibGen detail pages have download links in a table.
	// Look for EPUB links first, then PDF.
	epubLink := ""
	pdfLink := ""
	filename := ""
	fileType := ""

	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")
		text := strings.ToLower(strings.TrimSpace(sel.Text()))

		if href == "" {
			return
		}

		// Resolve relative URLs
		if !strings.HasPrefix(href, "http") {
			href = resolveLibGenURL(href)
		}

		if epubLink == "" && (strings.Contains(text, "epub") || strings.HasSuffix(strings.ToLower(href), ".epub")) {
			epubLink = href
			filename = extractFilename(href, "epub")
			fileType = "epub"
		} else if pdfLink == "" && (strings.Contains(text, "pdf") || strings.HasSuffix(strings.ToLower(href), ".pdf")) {
			pdfLink = href
			if filename == "" {
				filename = extractFilename(href, "pdf")
				fileType = "pdf"
			}
		}
	})

	// Prefer EPUB over PDF.
	downloadLink := epubLink
	if downloadLink == "" {
		downloadLink = pdfLink
	}
	if downloadLink == "" {
		return nil, ErrNoFileFound
	}

	if fileType == "" {
		fileType = detectFileType(downloadLink)
	}
	if filename == "" {
		filename = "book." + fileType
	}

	// Download the actual file.
	fileData, err := FetchFile(ctx, d.client, downloadLink)
	if err != nil {
		return nil, err
	}

	return &DownloadedFile{
		Data:     fileData,
		Filename: filename,
		FileType: fileType,
		Size:     int64(len(fileData)),
	}, nil
}

func resolveLibGenURL(href string) string {
	if strings.HasPrefix(href, "http") {
		return href
	}
	return "https://libgen.is" + href
}

func extractFilename(href string, defaultExt string) string {
	parts := strings.Split(href, "/")
	last := parts[len(parts)-1]
	last = strings.SplitN(last, "?", 2)[0]
	if last != "" && strings.Contains(last, ".") {
		return last
	}
	return "book." + defaultExt
}

func detectFileType(href string) string {
	lower := strings.ToLower(href)
	if strings.HasSuffix(lower, ".epub") {
		return "epub"
	}
	if strings.HasSuffix(lower, ".pdf") {
		return "pdf"
	}
	return "pdf"
}
