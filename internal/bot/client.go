package bot

import (
	"net/http"
	"strings"
	"time"
)

// DefaultUserAgent is a realistic browser User-Agent string to avoid basic bot detection.
const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// NewHTTPClient creates an HTTP client with browser-like headers, retry logic,
// and Cloudflare challenge detection for personal-use scraping.
func NewHTTPClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
				"User-Agent":                DefaultUserAgent,
				"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Language":           "en-US,en;q=0.9",
				"Accept-Encoding":           "gzip, deflate, br",
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

// IsCloudflareChallenge checks if an HTML response contains a Cloudflare challenge page.
func IsCloudflareChallenge(body string) bool {
	lower := strings.ToLower(body)
	return strings.Contains(lower, "checking your browser") ||
		(strings.Contains(lower, "cloudflare") && strings.Contains(lower, "challenge"))
}

// SleepWithDelay pauses execution for a small delay between requests to avoid rate limiting.
func SleepWithDelay() {
	time.Sleep(1500 * time.Millisecond)
}
