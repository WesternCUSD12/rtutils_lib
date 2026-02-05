# rtutils_lib

`rtutils_lib` is a robust Go client library for interacting with the Request Tracker (RT) REST API 2.0. It provides a clean, service-oriented interface for managing Assets, Tickets, and Users.

## Installation

```bash
go get rtutils_lib
```

## Configuration

To use the client, you'll need your RT instance's base URL and an authentication token.

```go
import "rtutils_lib"

baseURL := "https://your-rt-instance.com"
token := "your-rt-api-token"

client := rtutils_lib.NewClient(baseURL, token)
```

> [!TIP]
> The `NewClient` function automatically appends `/REST/2.0` to the base URL if it's not already present.

## API Reference

The client is organized into several services: `Assets`, `Tickets`, and `Users`.

### Assets Service

Handles all asset-related operations.

#### Methods
- `Create(ctx, asset)`: Creates a new asset. Returns the new asset's ID.
- `Get(ctx, id)`: Fetches an asset by its ID.
- `Search(ctx, query)`: Searches for assets using AssetSQL.
- `SearchByNameExact(ctx, name)`: Searches for assets with an exact name match.
- `SearchByNamePartial(ctx, query)`: Searches for assets with a partial name match.
- `SearchByCustomFieldExact(ctx, fieldName, value)`: Searches for assets with an exact custom field match.
- `SearchByCustomFieldPartial(ctx, fieldName, query)`: Searches for assets with a partial custom field match.
- `SearchWithCriteria(ctx, criteria)`: Searches for assets using JSON search criteria.
- `Update(ctx, id, asset)`: Updates an existing asset.
- `Delete(ctx, id)`: Deletes an asset.

#### Asset Struct
```go
type Asset struct {
    ID           json.Number        `json:"id,omitempty"`
    Name         string             `json:"Name,omitempty"`
    Description  string             `json:"Description,omitempty"`
    Status       string             `json:"Status,omitempty"`
    CustomFields []AssetCustomField `json:"CustomFields,omitempty"`
    // ... other fields
}
```

### Tickets Service

Handles ticket management, history, and communication.

#### Methods
- `Create(ctx, ticket)`: Creates a new ticket. Returns the new ticket's ID.
- `Get(ctx, id)`: Fetches a ticket by its ID.
- `Search(ctx, query, page, perPage)`: Searches for tickets using TicketSQL.
- `SearchBySubject(ctx, query)`: Searches for tickets by subject (partial match).
- `Update(ctx, id, ticket)`: Updates a ticket.
- `Comment(ctx, id, text)`: Adds a private comment to a ticket.
- `Correspond(ctx, id, text)`: Adds public correspondence to a ticket.
- `GetHistory(ctx, id)`: Fetches the transaction history for a ticket.
- `Take(ctx, id)` / `Untake(ctx, id)` / `Steal(ctx, id)`: Manage ticket ownership.

#### Ticket Struct
```go
type Ticket struct {
    ID           string                 `json:"id,omitempty"`
    Subject      string                 `json:"Subject,omitempty"`
    Status       string                 `json:"Status,omitempty"`
    Queue        string                 `json:"Queue,omitempty"`
    Owner        string                 `json:"Owner,omitempty"`
    Requestor    []string               `json:"Requestor,omitempty"`
    CustomFields map[string]interface{} `json:"CustomFields,omitempty"`
    // ... handles complex RT JSON formats automatically
}
```

### Users Service

Handles user discovery and management.

#### Methods
- `Get(ctx, id)`: Fetches a user by ID or username.
- `Search(ctx, query)`: Searches for users.
- `SearchByUsernameExact(ctx, username)`: Searches for users by exact username.
- `SearchByUsernamePartial(ctx, query)`: Searches for users by partial username.
- `SearchByEmailExact(ctx, email)`: Searches for users by exact email address.
- `SearchByEmailPartial(ctx, query)`: Searches for users by partial email address.
- `SearchByNameExact(ctx, name)`: Searches for users by exact full name.
- `SearchByNamePartial(ctx, query)`: Searches for users by partial full name.
- `GetGroupMemberships(ctx, id)`: Returns a list of group names the user belongs to.
- `AddToGroup(ctx, userID, groupID)` / `RemoveFromGroup(ctx, userID, groupID)`: Manage group memberships.

## LLM Agent Guide

If you are an AI agent using this library, follow these patterns for maximum reliability:

1.  **Client Initialization**: Always use `rtutils_lib.NewClient(url, token)`.
2.  **Searching Assets**: `client.Assets.Search(ctx, "query")` is powerful as it searches across multiple fields (Name, Serial, Tag) and automatically fetches full details for each result.
3.  **Custom Fields**:
    *   For **Assets**: Use `asset.GetCustomField("Field Name")` and `asset.SetCustomField("Field Name", "Value")`.
    *   For **Tickets**: Use the `CustomFields` map. Note that RT returns CFs as objects; if you only need the value, you may need to type-assert or inspect the map entry.
4.  **TicketSQL**: When searching tickets, use standard TicketSQL (e.g., `Queue = 'General' AND Status = 'open'`).
5.  **Error Handling**: Check if `err` is non-nil. The library returns `*rtutils_lib.APIError` which includes the `StatusCode` and the raw `Message` from RT.

### Example: Finding an Asset and Creating a Ticket
```go
// 1. Find asset by serial number
results, _ := client.Assets.Search(ctx, "SN12345")
if results.Count > 0 {
    asset := results.Items[0]
    
    // 2. Create ticket linked to asset
    ticket := &rtutils_lib.Ticket{
        Subject: "Repair Request: " + asset.Name,
        Queue:   "Hardware",
        Status:  "new",
    }
    ticketID, _ := client.Tickets.Create(ctx, ticket)
    
    // 3. Add initial comment
    client.Tickets.Comment(ctx, ticketID, "Found asset at " + asset.URL)
}
```

## Error Handling

All methods return an `error` interface. To access specific API error details:

```go
if err != nil {
    if apiErr, ok := err.(*rtutils_lib.APIError); ok {
        fmt.Printf("Status: %d, Response: %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```
