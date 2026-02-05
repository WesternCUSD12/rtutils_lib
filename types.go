package rtutils_lib

import (
	"encoding/json"
	"fmt"
	"strings"
)

// APIError represents an error returned by the Request Tracker API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("RT API Error: %d %s", e.StatusCode, e.Message)
}

// CustomFieldNotFoundError indicates a requested custom field does not exist in RT.
type CustomFieldNotFoundError struct {
	FieldName string
}

func (e *CustomFieldNotFoundError) Error() string {
	return fmt.Sprintf("custom field not found: %s", e.FieldName)
}

func requireNonEmpty(field string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	return nil
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
	Assets   []T    `json:"assets"`  // Support RT variant
	Tickets  []T    `json:"tickets"` // Support RT variant
	Queues   []T    `json:"queues"`  // Support RT variant
}

// Finalize ensures all variants are merged into Items.
func (r *SearchResult[T]) Finalize() {
	if len(r.Items) == 0 {
		if len(r.Assets) > 0 {
			r.Items = r.Assets
		} else if len(r.Tickets) > 0 {
			r.Items = r.Tickets
		} else if len(r.Queues) > 0 {
			r.Items = r.Queues
		}
	}
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
	ID          string      `json:"id"`
	Type        string      `json:"Type"`
	OldValue    string      `json:"OldValue,omitempty"`
	NewValue    string      `json:"NewValue,omitempty"`
	Field       string      `json:"Field,omitempty"`
	Data        string      `json:"Data,omitempty"`
	Description string      `json:"Description,omitempty"`
	Content     string      `json:"Content,omitempty"`
	Creator     string      `json:"Creator,omitempty"`
	Created     string      `json:"Created,omitempty"`
	Attachments []string    `json:"Attachments,omitempty"`
	Hyperlinks  []Hyperlink `json:"_hyperlinks,omitempty"`
}

// Hyperlink represents a link reference in API responses
type Hyperlink struct {
	URL  string      `json:"_url,omitempty"`
	ID   interface{} `json:"id,omitempty"`
	Ref  string      `json:"ref,omitempty"`
	Type string      `json:"type,omitempty"`
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

// Attachment represents an attachment to a transaction (typically contains comment content)
type Attachment struct {
	ID            interface{} `json:"id"`                // Can be number or string
	Content       string      `json:"Content,omitempty"` // Base64-encoded content
	ContentType   string      `json:"ContentType,omitempty"`
	Created       string      `json:"Created,omitempty"`
	Creator       interface{} `json:"Creator,omitempty"` // Can be string or object
	Subject       string      `json:"Subject,omitempty"`
	Headers       string      `json:"Headers,omitempty"`
	MessageId     string      `json:"MessageId,omitempty"`
	Parent        interface{} `json:"Parent,omitempty"`        // Can be number or object
	TransactionId interface{} `json:"TransactionId,omitempty"` // Can be number or object
}

// UnmarshalJSON handles custom unmarshaling for Attachment to handle numeric ID
func (a *Attachment) UnmarshalJSON(data []byte) error {
	type AttachmentAlias Attachment
	aux := struct {
		ID interface{} `json:"id"`
		*AttachmentAlias
	}{
		AttachmentAlias: (*AttachmentAlias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle ID (can be number or string, convert to string for consistency)
	switch v := aux.ID.(type) {
	case string:
		a.ID = v
	case float64:
		a.ID = fmt.Sprintf("%.0f", v)
	case int:
		a.ID = fmt.Sprintf("%d", v)
	}

	return nil
}
