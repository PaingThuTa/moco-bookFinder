package downloader

import "errors"

var (
	// ErrNoFileFound is returned when no EPUB or PDF download link is found on a detail page.
	ErrNoFileFound = errors.New("no downloadable file found on page")

	// ErrFileTooLarge is returned when the file exceeds Telegram's 50MB limit.
	ErrFileTooLarge = errors.New("file too large to send via Telegram (>50MB)")

	// ErrCloudflareBlocked is returned when a Cloudflare challenge page is detected.
	ErrCloudflareBlocked = errors.New("source blocked by Cloudflare challenge")
)
