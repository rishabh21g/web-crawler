package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rishabh21g/web-crawler/internal/models"
)

const (
	StatusUnsupportedMediaType = 415
	StatusPayloadTooLarge      = 413
	StatusInternalServerError  = 500
	StatusInvalidURL           = 1001
)

const MaxBodySize = 5 * 1024 * 1024 // 5MB

type HTTPError struct {
	Message string
	Code    int
}

func FetchHTML(ctx context.Context, task models.URLTask) (string, *HTTPError) {

	// 1. Parse and validate URL
	parsedURL, err := url.Parse(task.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", &HTTPError{
			Message: "Invalid URL",
			Code:    StatusInvalidURL,
		}
	}

	// 2. Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, task.URL, nil)
	if err != nil {
		return "", &HTTPError{
			Message: "Failed to create request: " + err.Error(),
			Code:    StatusInternalServerError,
		}
	}

	// 3. Headers
	req.Header.Set("User-Agent", "Crawler/1.0")
	req.Header.Set("Accept", "text/html")

	// 4. HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 5. Execute request
	res, err := client.Do(req)
	if err != nil {
		return "", &HTTPError{
			Message: "Request failed: " + err.Error(),
			Code:    StatusInternalServerError,
		}
	}
	defer res.Body.Close()

	// 6. Status check
	if res.StatusCode != http.StatusOK {
		return "", &HTTPError{
			Message: "Bad status: " + res.Status,
			Code:    res.StatusCode,
		}
	}

	// 7. Content-Type check
	contentType := res.Header.Get("Content-Type")
	if contentType == "" || !strings.Contains(contentType, "text/html") {
		return "", &HTTPError{
			Message: "Unsupported content type: " + contentType,
			Code:    StatusUnsupportedMediaType,
		}
	}

	// 8. Limit body size
	limitedReader := io.LimitReader(res.Body, MaxBodySize+1)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", &HTTPError{
			Message: "Failed to read body: " + err.Error(),
			Code:    StatusInternalServerError,
		}
	}

	// 9. Check if limit exceeded
	if int64(len(body)) > MaxBodySize {
		return "", &HTTPError{
			Message: "Response body too large",
			Code:    StatusPayloadTooLarge,
		}
	}

	return string(body), nil
}
