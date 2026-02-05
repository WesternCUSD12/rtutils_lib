# Data Model: RT API Wrapper

## Core Entities

### Ticket

Represents a Request Tracker ticket.

```go
type Ticket struct {
    ID           string                 `json:"id,omitempty"`
    URL          string                 `json:"_url,omitempty"`
    Type         string                 `json:"type,omitempty"`
    Subject      string                 `json:"Subject,omitempty"`
    Status       string                 `json:"Status,omitempty"`
    Queue        string                 `json:"Queue,omitempty"` // Name or ID
    Owner        string                 `json:"Owner,omitempty"` // Username or ID
    Requestor    []string               `json:"Requestor,omitempty"`
    Cc           []string               `json:"Cc,omitempty"`
    AdminCc      []string               `json:"AdminCc,omitempty"`
    CustomFields map[string]interface{} `json:"CustomFields,omitempty"`
    Created      string                 `json:"Created,omitempty"`
    Resolved     string                 `json:"Resolved,omitempty"`
    // ... potentially others
}
```

### User

Represents a Request Tracker user.

```go
type User struct {
    ID           string `json:"id,omitempty"`
    URL          string `json:"_url,omitempty"`
    Name         string `json:"Name,omitempty"`
    RealName     string `json:"RealName,omitempty"`
    EmailAddress string `json:"EmailAddress,omitempty"`
    Disabled     int    `json:"Disabled,omitempty"` // 0 or 1
    Privileged   int    `json:"Privileged,omitempty"`
}
```

### Asset

Represents a Request Tracker asset.

```go
type Asset struct {
    ID           string                 `json:"id,omitempty"`
    URL          string                 `json:"_url,omitempty"`
    Name         string                 `json:"Name,omitempty"`
    Catalog      string                 `json:"Catalog,omitempty"`
    Content      string                 `json:"Content,omitempty"`
    Status       string                 `json:"Status,omitempty"`
    CustomFields map[string]interface{} `json:"CustomFields,omitempty"`
}
```

### Responses

**SearchResult**

```go
type SearchResult[T any] struct {
    Total      int    `json:"total"`
    Count      int    `json:"count"`
    Page       int    `json:"page"`
    Pages      int    `json:"pages"`
    PerPage    int    `json:"per_page"`
    NextPage   string `json:"next_page"`
    Items      []T    `json:"items"`
}
```

**OperationResponse**

```go
// For simple messages like "Ticket 123 created" or array of messages
type ActionResult struct {
    ID      string `json:"id"`
    URL     string `json:"_url"`
    Type    string `json:"type"`
    Message string // Consolidates message strings
}
```
