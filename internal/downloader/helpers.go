package downloader

import (
	"net/http"
	"strings"
	"time"
)

// defaultUserAgent is a realistic browser User-Agent string.
const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// isCloudflareChallenge checks if an HTML response contains a Cloudflare challenge page.
func isCloudflareChallenge(body string) bool {
	lower := strings.ToLower(body)
	return strings.Contains(lower, "checking your browser") ||
		(strings.Contains(lower, "cloudflare") && strings.Contains(lower, "challenge"))
}

// sleepWithDelay pauses execution between requests to avoid rate limiting.
func sleepWithDelay() {
	time.Sleep(1500 * time.Millisecond)
}

// NewHTTPClient creates an HTTP client with browser-like headers.
func NewHTTPClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
				"User-Agent":                defaultUserAgent,
				"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Language":           "en-US,en;q=0.9",
				"Connection":                "keep-alive",
				"Upgrade-Insecure-Requests": "1",
			},
		},
	}
	return client
}

// headerTransport adds default headers to every request.
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}
	return t.base.RoundTrip(req)
}
