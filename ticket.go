package rtutils_lib

import (
	"context"
	"fmt"
	"net/url"
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
