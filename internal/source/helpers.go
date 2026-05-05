package source

import (
	"net/http"
	"strings"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// isCloudflareChallenge checks if an HTML response contains a Cloudflare challenge page.
func isCloudflareChallenge(body string) bool {
	lower := strings.ToLower(body)
	return strings.Contains(lower, "checking your browser") ||
		(strings.Contains(lower, "cloudflare") && strings.Contains(lower, "challenge"))
}

// setBrowserHeaders ensures a request looks like a real browser.
func setBrowserHeaders(req *http.Request) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}
}
