package rtutils_lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client represents a Request Tracker API client.
type Client struct {
	baseURL string
	token   string
	client  *http.Client

	Tickets *TicketService
	Users   *UserService
	Assets  *AssetService
}

// NewClient creates a new Request Tracker API client.
func NewClient(baseURL, token string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/REST/2.0") {
		baseURL += "/REST/2.0"
	}

	c := &Client{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
	c.Tickets = &TicketService{client: c}
	c.Users = &UserService{client: c}
	c.Assets = &AssetService{client: c}
	return c
}

func (c *Client) request(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	url := c.baseURL + path
	// If path starts with http, assume it's a full URL (e.g. from _url field)
	if strings.HasPrefix(path, "http") {
		url = path
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		// Try to read message from body if possible, though RT might not always send JSON error
		bodyBytes, _ := io.ReadAll(resp.Body)
		apiErr.Message = string(bodyBytes)
		return &apiErr
	}

	if out != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if err := json.Unmarshal(bodyBytes, out); err != nil {
			// Include a sample of the body in the error message for debugging
			sample := string(bodyBytes)
			if len(sample) > 500 {
				sample = sample[:500] + "..."
			}
			return fmt.Errorf("failed to decode response: %w. Body: %s", err, sample)
		}
	}

	return nil
}
