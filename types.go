package rtutils_lib

import (
	"encoding/json"
	"fmt"
)

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

// UnmarshalJSON handles custom unmarshaling for Transaction to support case-insensitive fields
func (t *Transaction) UnmarshalJSON(data []byte) error {
	type TransactionAlias Transaction
	aux := struct {
		ID       interface{} `json:"id"`
		Creator  interface{} `json:"Creator"`
		OldValue interface{} `json:"OldValue"`
		NewValue interface{} `json:"NewValue"`
		*TransactionAlias
	}{
		TransactionAlias: (*TransactionAlias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle ID (can be number or string)
	switch v := aux.ID.(type) {
	case string:
		t.ID = v
	case float64:
		t.ID = fmt.Sprintf("%.0f", v)
	case int:
		t.ID = fmt.Sprintf("%d", v)
	}

	// Handle Creator (can be string or object)
	t.Creator = parseStringOrObject(aux.Creator)

	// Handle OldValue (can be string or object)
	t.OldValue = parseStringOrObject(aux.OldValue)

	// Handle NewValue (can be string or object)
	t.NewValue = parseStringOrObject(aux.NewValue)

	return nil
}

// parseStringOrObject extracts a string value from either a string or an object with "id" field
func parseStringOrObject(field interface{}) string {
	if field == nil {
		return ""
	}

	switch v := field.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Extract ID from object like {"id": "value", "_url": "...", "type": "..."}
		if id, ok := v["id"].(string); ok {
			return id
		}
	}

	return ""
}
