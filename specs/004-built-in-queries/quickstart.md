# Quickstart: Built-in Query Methods

This quickstart demonstrates the new exact/partial query methods for assets, tickets, and users, plus guidance for custom queries.

## Setup

```go
import "context"
import "rtutils_lib"

ctx := context.Background()
client := rtutils_lib.NewClient("https://your-rt-instance.example", "your-token")
```

## Asset Queries

Exact name lookup:

```go
assets, err := client.Assets.SearchByNameExact(ctx, "MacBook-1234")
if err != nil {
    // handle error
}
```

Partial name lookup:

```go
assets, err := client.Assets.SearchByNamePartial(ctx, "MacBook")
```

Exact custom field lookup:

```go
assets, err := client.Assets.SearchByCustomFieldExact(ctx, "Internal Name", "COMP-12345")
if err != nil {
    // If the custom field does not exist, an error is returned.
}
```

Partial custom field lookup:

```go
assets, err := client.Assets.SearchByCustomFieldPartial(ctx, "Internal Name", "COMP")
```

## Ticket Queries

Subject contains:

```go
tickets, err := client.Tickets.SearchBySubject(ctx, "Smartboard")
```

## User Queries

Exact username:

```go
users, err := client.Users.SearchByUsernameExact(ctx, "jsmith")
```

Partial username:

```go
users, err := client.Users.SearchByUsernamePartial(ctx, "smith")
```

Exact email:

```go
users, err := client.Users.SearchByEmailExact(ctx, "jsmith@example.com")
```

Partial email:

```go
users, err := client.Users.SearchByEmailPartial(ctx, "@example.com")
```

Exact name:

```go
users, err := client.Users.SearchByNameExact(ctx, "John Smith")
```

Partial name:

```go
users, err := client.Users.SearchByNamePartial(ctx, "Smith")
```

## Custom Queries

Use the existing methods for complex queries or pagination beyond the first page.

### Asset JSON Search (SearchWithCriteria)

```go
criteria := []map[string]interface{}{
    {
        "field":    "Name",
        "operator": "LIKE",
        "value":    "Laptop",
    },
    {
        "field":            "CustomField.{Serial Number}",
        "operator":         "LIKE",
        "value":            "SN-",
        "entry_aggregator": "AND",
    },
}

results, err := client.Assets.SearchWithCriteria(ctx, criteria)
```

### TicketSQL Query

```go
query := "Queue = 'Hardware' AND Status = 'open' AND Subject LIKE 'Smartboard'"
results, err := client.Tickets.Search(ctx, query, 1, 20)
```

### User Query

```go
query := "Email LIKE 'smith@%' AND Disabled = '0'"
results, err := client.Users.Search(ctx, query)
```

## Pagination Notes

Built-in query methods return only the first page of results. If you need additional pages, use `Search`, `SearchWithCriteria`, or the RT API pagination fields on the SearchResult struct.
