package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const Timeout = 10 * time.Second

type Response struct {
	Body        []byte
	ETag        string
	NotModified bool
}

type Client struct {
	HTTP *http.Client
}

func New() *Client {
	return &Client{
		HTTP: &http.Client{Timeout: Timeout},
	}
}

func (c *Client) Get(ctx context.Context, url string, etag string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting %s: %w", url, err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotModified:
		return &Response{NotModified: true}, nil
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("requesting %s: unexpected status %s", url, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", url, err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("reading response from %s: empty body", url)
	}

	return &Response{
		Body: body,
		ETag: res.Header.Get("ETag"),
	}, nil
}
