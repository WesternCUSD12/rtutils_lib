package rtutils_lib

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// TicketService handles communication with the ticket related methods of the
// Request Tracker API.
type TicketService struct {
	client *Client
}

// Create creates a new ticket.
func (s *TicketService) Create(ctx context.Context, ticket *Ticket) (string, error) {
	var result ActionResult
	err := s.client.request(ctx, "POST", "/ticket", ticket, &result)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

// Get fetches a ticket by ID.
func (s *TicketService) Get(ctx context.Context, id string) (*Ticket, error) {
	var ticket Ticket
	err := s.client.request(ctx, "GET", "/ticket/"+id, nil, &ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// GetByURL fetches a ticket using a full RT URL (e.g. from _url in search results).
func (s *TicketService) GetByURL(ctx context.Context, url string) (*Ticket, error) {
	var ticket Ticket
	err := s.client.request(ctx, "GET", url, nil, &ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// Search searches for tickets using TicketSQL.
func (s *TicketService) Search(ctx context.Context, query string, page int, perPage int) (*SearchResult[Ticket], error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	path := fmt.Sprintf("/tickets?query=%s&page=%d&per_page=%d", url.QueryEscape(query), page, perPage)

	var result SearchResult[Ticket]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// SearchBySubject searches for tickets whose subject contains the query string.
func (s *TicketService) SearchBySubject(ctx context.Context, query string) (*SearchResult[Ticket], error) {
	if err := requireNonEmpty("query", query); err != nil {
		return nil, err
	}
	escaped := strings.ReplaceAll(query, "'", "''")
	searchQuery := fmt.Sprintf("Subject LIKE '%%%s%%'", escaped)
	return s.Search(ctx, searchQuery, 1, 20)
}

// Update updates a ticket.
func (s *TicketService) Update(ctx context.Context, id string, ticket *Ticket) error {
	path := "/ticket/" + id
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, ticket, &result)
	return err
}

// Delete deletes a ticket.
func (s *TicketService) Delete(ctx context.Context, id string) error {
	path := "/ticket/" + id
	var result ActionResult
	err := s.client.request(ctx, "DELETE", path, nil, &result)
	return err
}

// GetHistory fetches the history of a ticket.
func (s *TicketService) GetHistory(ctx context.Context, id string) ([]Transaction, error) {
	path := "/ticket/" + id + "/history"
	var result SearchResult[Transaction]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// GetTransaction fetches a single transaction by ID with full details.
func (s *TicketService) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	// Use expand parameter to get all transaction details including attachments
	path := "/transaction/" + id + "?expand=true"
	var transaction Transaction
	err := s.client.request(ctx, "GET", path, nil, &transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetAttachment fetches an attachment by ID.
func (s *TicketService) GetAttachment(ctx context.Context, id string) (*Attachment, error) {
	path := "/attachment/" + id
	var result Attachment
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCorrespondence fetches correspondence/comments for a ticket.
func (s *TicketService) GetCorrespondence(ctx context.Context, ticketID string) ([]map[string]interface{}, error) {
	path := "/ticket/" + ticketID + "/correspond"
	var result SearchResult[map[string]interface{}]
	err := s.client.request(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// Comment adds a comment to a ticket.
func (s *TicketService) Comment(ctx context.Context, id string, text string) error {
	path := "/ticket/" + id + "/comment"
	payload := map[string]string{
		"Content":     text,
		"ContentType": "text/plain",
	}
	var result ActionResult
	err := s.client.request(ctx, "POST", path, payload, &result)
	return err
}

// Correspond adds correspondence to a ticket.
func (s *TicketService) Correspond(ctx context.Context, id string, text string) error {
	path := "/ticket/" + id + "/correspond"
	payload := map[string]string{
		"Content":     text,
		"ContentType": "text/plain",
	}
	var result ActionResult
	err := s.client.request(ctx, "POST", path, payload, &result)
	return err
}

// Take takes a ticket.
func (s *TicketService) Take(ctx context.Context, id string) error {
	path := "/ticket/" + id + "/take"
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, nil, &result)
	return err
}

// Untake untakes a ticket.
func (s *TicketService) Untake(ctx context.Context, id string) error {
	path := "/ticket/" + id + "/untake"
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, nil, &result)
	return err
}

// Steal steals a ticket.
func (s *TicketService) Steal(ctx context.Context, id string) error {
	path := "/ticket/" + id + "/steal"
	var result ActionResult
	err := s.client.request(ctx, "PUT", path, nil, &result)
	return err
}

// BulkCreate creates multiple tickets.
func (s *TicketService) BulkCreate(ctx context.Context, tickets []Ticket) error {
	for _, t := range tickets {
		_, err := s.Create(ctx, &t)
		if err != nil {
			return err
		}
	}
	return nil
}

// BulkUpdate updates multiple tickets.
func (s *TicketService) BulkUpdate(ctx context.Context, tickets []Ticket) error {
	for _, t := range tickets {
		if t.ID == "" {
			continue // Skip tickets without ID
		}
		err := s.Update(ctx, t.ID, &t)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ticket represents a Request Tracker ticket.
type Ticket struct {
	ID           string                 `json:"id,omitempty"`
	URL          string                 `json:"_url,omitempty"`
	Type         string                 `json:"type,omitempty"`
	Subject      string                 `json:"Subject,omitempty"`
	Status       string                 `json:"Status,omitempty"`
	Queue        string                 `json:"Queue,omitempty"`
	Owner        string                 `json:"Owner,omitempty"`
	Requestor    []string               `json:"Requestor,omitempty"`
	Cc           []string               `json:"Cc,omitempty"`
	AdminCc      []string               `json:"AdminCc,omitempty"`
	CustomFields map[string]interface{} `json:"CustomFields,omitempty"`
	Created      string                 `json:"Created,omitempty"`
	Resolved     string                 `json:"Resolved,omitempty"`
}

// UnmarshalJSON handles custom unmarshaling for Ticket to support RT's format
// where Requestor, Cc, and AdminCc can be arrays of strings or objects,
// Queue and Owner can be strings or objects, ID can be a number or string,
// and CustomFields can be an array or object
func (t *Ticket) UnmarshalJSON(data []byte) error {
	type TicketAlias Ticket
	aux := struct {
		ID           interface{} `json:"id"`
		Queue        interface{} `json:"Queue"`
		Owner        interface{} `json:"Owner"`
		Requestor    interface{} `json:"Requestor"`
		Cc           interface{} `json:"Cc"`
		AdminCc      interface{} `json:"AdminCc"`
		CustomFields interface{} `json:"CustomFields"`
		*TicketAlias
	}{
		TicketAlias: (*TicketAlias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle ID (can be number or string)
	t.ID = parseIDField(aux.ID)
	// Handle Queue
	t.Queue = parseStringField(aux.Queue)
	// Handle Owner
	t.Owner = parseStringField(aux.Owner)
	// Handle Requestor
	t.Requestor = parseContactField(aux.Requestor)
	// Handle Cc
	t.Cc = parseContactField(aux.Cc)
	// Handle AdminCc
	t.AdminCc = parseContactField(aux.AdminCc)
	// Handle CustomFields
	t.CustomFields = parseCustomFields(aux.CustomFields)

	return nil
}

// parseIDField handles ID which can be a number or string
func parseIDField(field interface{}) string {
	if field == nil {
		return ""
	}

	switch v := field.(type) {
	case string:
		return v
	case float64:
		// JSON unmarshals numbers to float64
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	}

	return ""
}

// parseStringField handles both strings and objects (extracting "id" field)
func parseStringField(field interface{}) string {
	if field == nil {
		return ""
	}

	switch v := field.(type) {
	case string:
		return v
	case map[string]interface{}:
		// Extract ID from object like {"id": "queue", "_url": "...", "type": "queue"}
		if id, ok := v["id"].(string); ok {
			return id
		}
	}

	return ""
}

// parseContactField handles both string arrays and object arrays
func parseContactField(field interface{}) []string {
	if field == nil {
		return nil
	}

	result := []string{}

	switch v := field.(type) {
	case []interface{}:
		for _, item := range v {
			switch itemVal := item.(type) {
			case string:
				result = append(result, itemVal)
			case map[string]interface{}:
				// Extract ID from object like {"id": "user", "_url": "...", "type": "user"}
				if id, ok := itemVal["id"].(string); ok {
					result = append(result, id)
				}
			}
		}
	case []string:
		result = v
	}

	return result
}

// parseCustomFields handles CustomFields which can be an array (when empty) or object
func parseCustomFields(field interface{}) map[string]interface{} {
	if field == nil {
		return nil
	}

	switch v := field.(type) {
	case map[string]interface{}:
		return v
	case []interface{}:
		// If it's an empty array, return nil
		if len(v) == 0 {
			return nil
		}
		// If it's an array with objects, try to convert to map
		result := make(map[string]interface{})
		for _, item := range v {
			if cfObj, ok := item.(map[string]interface{}); ok {
				if name, ok := cfObj["name"].(string); ok {
					result[name] = cfObj
				}
			}
		}
		return result
	}

	return nil
}

// DecodedContent returns the base64-decoded content of the attachment
func (a *Attachment) DecodedContent() (string, error) {
	if a.Content == "" {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(a.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode attachment content: %w", err)
	}

	return string(decoded), nil
}

// GetAttachmentIDs returns the IDs of all attachments linked to this transaction
func (t *Transaction) GetAttachmentIDs() []string {
	var ids []string
	for _, hl := range t.Hyperlinks {
		if hl.Ref == "attachment" {
			// Extract ID from URL like https://tickets.wc-12.com/REST/2.0/attachment/1206
			parts := strings.Split(strings.TrimSpace(hl.URL), "/")
			if len(parts) > 0 {
				id := parts[len(parts)-1]
				if id != "" {
					ids = append(ids, id)
				}
			}
		}
	}
	return ids
}
