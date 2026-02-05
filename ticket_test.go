package rtutils_lib

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestTicketService_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock response
	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/ticket",
		httpmock.NewStringResponder(201, `{"id": "123", "_url": "http://rt.example.com/REST/2.0/ticket/123", "type": "ticket"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	ticket := &Ticket{
		Queue:   "General",
		Subject: "New Ticket",
	}

	id, err := client.Tickets.Create(context.Background(), ticket)
	assert.NoError(t, err)
	assert.Equal(t, "123", id)

	// Verify request body
	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST http://rt.example.com/REST/2.0/ticket"])
}

func TestTicketService_Get(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/ticket/123",
		httpmock.NewStringResponder(200, `{"id": "123", "Subject": "Existing Ticket", "type": "ticket"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	ticket, err := client.Tickets.Get(context.Background(), "123")
	assert.NoError(t, err)
	assert.Equal(t, "123", ticket.ID)
	assert.Equal(t, "Existing Ticket", ticket.Subject)
}

func TestTicketService_Search(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/tickets",
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query().Get("query")
			page := req.URL.Query().Get("page")
			if query == "Queue = 'General'" && page == "1" {
				return httpmock.NewStringResponse(200, `{"total": 5, "page": 1, "items": [{"id": "1", "Subject": "T1"}, {"id": "2", "Subject": "T2"}]}`), nil
			}
			return httpmock.NewStringResponse(400, "Bad Request"), nil
		})

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	result, err := client.Tickets.Search(context.Background(), "Queue = 'General'", 1)
	assert.NoError(t, err)
	assert.Equal(t, 5, result.Total)
	assert.Equal(t, 2, len(result.Items))
	assert.Equal(t, "T1", result.Items[0].Subject)
}

func TestTicketService_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/ticket/123",
		httpmock.NewStringResponder(200, `{"id": "123", "type": "ticket", "message": "Ticket 123 updated"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	updates := &Ticket{
		Status: "resolved",
	}

	err := client.Tickets.Update(context.Background(), "123", updates)
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["PUT http://rt.example.com/REST/2.0/ticket/123"])
}

func TestTicketService_Delete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "http://rt.example.com/REST/2.0/ticket/123",
		httpmock.NewStringResponder(200, `{"id": "123", "type": "ticket", "message": "Ticket 123 deleted"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	err := client.Tickets.Delete(context.Background(), "123")
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["DELETE http://rt.example.com/REST/2.0/ticket/123"])
}

func TestTicketService_GetHistory(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://rt.example.com/REST/2.0/ticket/123/history",
		httpmock.NewStringResponder(200, `{"items": [{"id": "1", "Type": "Create"}]}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	history, err := client.Tickets.GetHistory(context.Background(), "123")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(history))
	assert.Equal(t, "Create", history[0].Type)
}

func TestTicketService_Comment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/ticket/123/comment",
		httpmock.NewStringResponder(200, `{"message": "Comment recorded"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	err := client.Tickets.Comment(context.Background(), "123", "Test Comment")
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST http://rt.example.com/REST/2.0/ticket/123/comment"])
}

func TestTicketService_Correspond(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/ticket/123/correspond",
		httpmock.NewStringResponder(200, `{"message": "Correspondence recorded"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	err := client.Tickets.Correspond(context.Background(), "123", "Test Reply")
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST http://rt.example.com/REST/2.0/ticket/123/correspond"])
}

func TestTicketService_Actions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Take
	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/ticket/123/take",
		httpmock.NewStringResponder(200, `{"message": "Ticket taken"}`))
	// Untake
	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/ticket/123/untake",
		httpmock.NewStringResponder(200, `{"message": "Ticket untaken"}`))
	// Steal
	httpmock.RegisterResponder("PUT", "http://rt.example.com/REST/2.0/ticket/123/steal",
		httpmock.NewStringResponder(200, `{"message": "Ticket stolen"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	assert.NoError(t, client.Tickets.Take(context.Background(), "123"))
	assert.NoError(t, client.Tickets.Untake(context.Background(), "123"))
	assert.NoError(t, client.Tickets.Steal(context.Background(), "123"))
}

func TestTicketService_Bulk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Bulk Create (simulated serial)
	httpmock.RegisterResponder("POST", "http://rt.example.com/REST/2.0/ticket",
		httpmock.NewStringResponder(201, `{"id": "new", "type": "ticket"}`))

	client := NewClient("http://rt.example.com/REST/2.0", "test-token")
	httpmock.ActivateNonDefault(client.client)

	tickets := []Ticket{{Subject: "T1"}, {Subject: "T2"}}
	err := client.Tickets.BulkCreate(context.Background(), tickets)
	assert.NoError(t, err)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 2, info["POST http://rt.example.com/REST/2.0/ticket"])
}
