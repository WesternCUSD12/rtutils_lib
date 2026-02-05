# RT API Wrapper Quickstart

## Installation

```bash
go get github.com/yourusername/rtutils_lib
```

## Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/rtutils_lib"
)

func main() {
	// Initialize Client
	client := rtutils_lib.NewClient("https://rt.example.com/REST/2.0", "your-api-token")

	ctx := context.Background()

	// 1. Create a Ticket
	newTicket := &rtutils_lib.Ticket{
		Subject: "Printer on fire",
		Queue:   "General",
		CustomFields: map[string]interface{}{
			"Severity": "Critical",
		},
	}
	
	createdTicket, err := client.Tickets().Create(ctx, newTicket)
	if err != nil {
		log.Fatalf("Failed to create ticket: %v", err)
	}
	fmt.Printf("Created Ticket ID: %s, URL: %s\n", createdTicket.ID, createdTicket.URL)

	// 2. Search for Tickets
	// "Status = 'new'" is TicketSQL
	results, err := client.Tickets().Search(ctx, "Status = 'new'", false)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	
	for _, t := range results.Items {
		fmt.Printf("Found ticket: %s - %s\n", t.ID, t.Subject)
		// Note provided _url field
		fmt.Printf("API Link: %s\n", t.URL) 
	}

	// 3. User Management
	user, err := client.Users().Get(ctx, "jdoe")
	if err == nil {
		fmt.Printf("User %s (RealName: %s)\n", user.Name, user.RealName)
	}
}
```
