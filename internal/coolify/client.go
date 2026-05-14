package coolify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("%s %s failed with HTTP %d", e.Method, e.Path, e.StatusCode)
	}
	return fmt.Sprintf("%s %s failed with HTTP %d: %s", e.Method, e.Path, e.StatusCode, e.Body)
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, body any, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

func (c *Client) Patch(ctx context.Context, path string, body any, out any) error {
	return c.do(ctx, http.MethodPatch, path, body, out)
}

func (c *Client) Delete(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodDelete, path, nil, out)
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}

	fullURL, err := c.url(path)
	if err != nil {
		return err
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode %s %s request body: %w", method, path, err)
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reader)
	if err != nil {
		return fmt.Errorf("create %s %s request: %w", method, path, err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send %s %s request: %w", method, path, err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read %s %s response: %w", method, path, err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return &HTTPError{
			Method:     method,
			Path:       path,
			StatusCode: res.StatusCode,
			Body:       strings.TrimSpace(string(resBody)),
		}
	}

	if out == nil || len(bytes.TrimSpace(resBody)) == 0 {
		return nil
	}

	if err := json.Unmarshal(resBody, out); err != nil {
		return fmt.Errorf("decode %s %s response: %w", method, path, err)
	}
	return nil
}

func (c *Client) url(path string) (string, error) {
	if c.BaseURL == "" {
		return "", fmt.Errorf("Coolify endpoint is empty")
	}
	base, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse Coolify endpoint %q: %w", c.BaseURL, err)
	}
	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("Coolify endpoint must include scheme and host")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return strings.TrimRight(base.String(), "/") + path, nil
}
