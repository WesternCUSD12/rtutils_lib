package rtutils_lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	Queues  *QueueService
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
	c.Queues = &QueueService{client: c}
	return c
}

// QueueService handles communication with the queue related methods of the RT API.
type QueueService struct {
	client *Client
}

type Queue struct {
	ID          string `json:"id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// UnmarshalJSON handles inconsistencies in RT's field naming and types
func (q *Queue) UnmarshalJSON(data []byte) error {
	type Alias Queue
	aux := &struct {
		ID          interface{} `json:"id"`
		Name        interface{} `json:"Name"`
		Description interface{} `json:"Description"`
		*Alias
	}{
		Alias: (*Alias)(q),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle ID which can be numeric
	q.ID = parseIDField(aux.ID)

	// Try to get Name from various possible sources
	q.Name = parseStringOrObject(aux.Name)
	if q.Name == "" {
		// Try lowercase 'name'
		var raw map[string]interface{}
		json.Unmarshal(data, &raw)
		if n, ok := raw["name"].(string); ok {
			q.Name = n
		}
	}

	// Try to get Description from various sources
	q.Description = parseStringOrObject(aux.Description)
	if q.Description == "" {
		var raw map[string]interface{}
		json.Unmarshal(data, &raw)
		if d, ok := raw["description"].(string); ok {
			q.Description = d
		}
	}

	return nil
}

func (s *QueueService) Search(ctx context.Context, query string) (*SearchResult[Queue], error) {
	path := "/queues"
	if query == "" {
		path = "/queues/all"
	} else {
		path += "?query=" + url.QueryEscape(query)
	}
	var result SearchResult[Queue]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	result.Finalize()

	// Automatically expand if names are missing
	needsExpansion := false
	for _, q := range result.Items {
		if q.Name == "" {
			needsExpansion = true
			break
		}
	}
	if needsExpansion && len(result.Items) > 0 {
		s.Expand(ctx, &result)
	}

	return &result, nil
}

// Expand fetches full details for each queue in the search result.
func (s *QueueService) Expand(ctx context.Context, result *SearchResult[Queue]) error {
	log.Printf("DEBUG: Queues: Expanding %d search results...", len(result.Items))
	for i := range result.Items {
		var fullQueue Queue
		err := s.client.request(ctx, "GET", "/queue/"+result.Items[i].ID, nil, &fullQueue)
		if err != nil {
			log.Printf("DEBUG: Queues: Failed to expand queue %s: %v", result.Items[i].ID, err)
			continue
		}
		result.Items[i] = fullQueue
	}
	return nil
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

	bodyBytes, _ := io.ReadAll(resp.Body)

	// Log raw response for debugging via standard logger
	log.Printf("DEBUG: API Response [%d] from %s: %s", resp.StatusCode, url, string(bodyBytes))

	if resp.StatusCode >= 400 {
		var apiErr APIError
		apiErr.StatusCode = resp.StatusCode
		apiErr.Message = string(bodyBytes)
		return &apiErr
	}

	if out != nil {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			sample := string(bodyBytes)
			if len(sample) > 500 {
				sample = sample[:500] + "..."
			}
			return fmt.Errorf("failed to decode response: %w. Body: %s", err, sample)
		}
	}

	return nil
}
