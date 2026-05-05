package downloader

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ZLibraryDownloader struct {
	client *http.Client
}

func NewZLibraryDownloader(client *http.Client) *ZLibraryDownloader {
	return &ZLibraryDownloader{client: client}
}

func (d *ZLibraryDownloader) SourceName() string {
	return "Z-Library"
}

func (d *ZLibraryDownloader) DownloadFile(ctx context.Context, detailURL string) (*DownloadedFile, error) {
	// Fetch the detail page
	data, err := FetchFile(ctx, d.client, detailURL)
	if err != nil {
		return nil, err
	}

	html := string(data)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Z-Library detail pages have EPUB/PDF download buttons
	// Look for buttons with EPUB or PDF in text or href
	epubLink := ""
	pdfLink := ""
	filename := ""
	fileType := ""

	doc.Find("a, button, [href]").Each(func(i int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")
		text := strings.ToLower(strings.TrimSpace(sel.Text()))

		if href == "" {
			return
		}

		// Resolve relative URLs
		if !strings.HasPrefix(href, "http") {
			href = strings.TrimSuffix("https://z-library.bz", "/") + href
		}

		if epubLink == "" && (strings.Contains(text, "epub") || strings.Contains(strings.ToLower(href), "epub")) {
			epubLink = href
			filename = extractFilename(href, "epub")
			fileType = "epub"
		} else if pdfLink == "" && (strings.Contains(text, "pdf") || strings.Contains(strings.ToLower(href), "pdf")) {
			pdfLink = href
			if filename == "" {
				filename = extractFilename(href, "pdf")
				fileType = "pdf"
			}
		}
	})

	// Prefer EPUB over PDF
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

	// Download the actual file
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
