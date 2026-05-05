package source

import (
	"fmt"
	"strings"
)

type BookResult struct {
	Title       string
	Author      string
	DetailURL   string // URL to book detail page for file resolution
	DownloadURL string // Direct download link or search result link
	Source      string // "Z-Library", "Ocean of PDF", "LibGen"
}

func FormatBooks(results []BookResult) string {
	if len(results) == 0 {
		return "No results found."
	}
	var sb strings.Builder
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("%d. **%s**", i+1, r.Title))
		if r.Author != "" {
			sb.WriteString(fmt.Sprintf(" by %s", r.Author))
		}
		sb.WriteString(fmt.Sprintf("\n   Source: %s\n\n", r.Source))
	}
	return sb.String()
}
