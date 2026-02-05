package rtutils_lib

import "fmt"

// APIError represents an error returned by the Request Tracker API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("RT API Error: %d %s", e.StatusCode, e.Message)
}

// SearchResult represents a paginated search response.
type SearchResult[T any] struct {
	Total    int    `json:"total"`
	Count    int    `json:"count"`
	Page     int    `json:"page"`
	Pages    int    `json:"pages"`
	PerPage  int    `json:"per_page"`
	NextPage string `json:"next_page"`
	Items    []T    `json:"items"`
}

// ActionResult represents the response from a create or update operation.
type ActionResult struct {
	ID      string `json:"id"`
	URL     string `json:"_url"`
	Type    string `json:"type"`
	Message string `json:"-"` // Message is often returned separately or needs custom handling
}

// Transaction represents a history entry for an object.
type Transaction struct {
	ID          string   `json:"id"`
	Type        string   `json:"Type"`
	OldValue    string   `json:"OldValue,omitempty"`
	NewValue    string   `json:"NewValue,omitempty"`
	Field       string   `json:"Field,omitempty"`
	Data        string   `json:"Data,omitempty"`
	Description string   `json:"Description,omitempty"`
	Content     string   `json:"Content,omitempty"`
	Creator     string   `json:"Creator,omitempty"`
	Created     string   `json:"Created,omitempty"`
	Attachments []string `json:"Attachments,omitempty"`
}
