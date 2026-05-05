package downloader

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// OceanPDFDownloader resolves Ocean of PDF detail pages to downloadable EPUB or PDF files.
type OceanPDFDownloader struct {
	client *http.Client
}

// NewOceanPDFDownloader creates a new downloader for Ocean of PDF.
func NewOceanPDFDownloader(client *http.Client) *OceanPDFDownloader {
	return &OceanPDFDownloader{client: client}
}

func (d *OceanPDFDownloader) SourceName() string {
	return "Ocean of PDF"
}

func (d *OceanPDFDownloader) DownloadFile(ctx context.Context, detailURL string) (*DownloadedFile, error) {
	data, err := FetchFile(ctx, d.client, detailURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch detail page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	epubLink := ""
	pdfLink := ""
	filename := ""
	fileType := ""

	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		href, _ := sel.Attr("href")
		text := strings.ToLower(strings.TrimSpace(sel.Text()))
		class, _ := sel.Attr("class")

		if href == "" {
			return
		}

		if !strings.HasPrefix(href, "http") {
			href = strings.TrimSuffix("https://oceanofpdf.com", "/") + href
		}

		isDownload := strings.Contains(text, "download") ||
			strings.Contains(text, "epub") ||
			strings.Contains(text, "pdf") ||
			strings.Contains(strings.ToLower(class), "download") ||
			strings.Contains(strings.ToLower(href), ".epub") ||
			strings.Contains(strings.ToLower(href), ".pdf")

		if !isDownload {
			return
		}

		if epubLink == "" && (strings.Contains(text, "epub") || strings.Contains(strings.ToLower(href), ".epub")) {
			epubLink = href
			filename = extractFilename(href, "epub")
			fileType = "epub"
		} else if pdfLink == "" && (strings.Contains(text, "pdf") || strings.Contains(strings.ToLower(href), ".pdf")) {
			pdfLink = href
			if filename == "" {
				filename = extractFilename(href, "pdf")
				fileType = "pdf"
			}
		}
	})

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

	fileData, err := FetchFile(ctx, d.client, downloadLink)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return &DownloadedFile{
		Data:     fileData,
		Filename: filename,
		FileType: fileType,
		Size:     int64(len(fileData)),
	}, nil
}

