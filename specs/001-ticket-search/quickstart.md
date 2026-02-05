# Quickstart: Ticket Search Example

## Prerequisites

- Go 1.22+
- `RT_BASE_URL` and `RT_TOKEN` set in environment or `.env`

## Run

```bash
go run examples/ticket_search/main.go --queue "Helpdesk" --status "open" --subject "printer"
```

## Paginate

```bash
go run examples/ticket_search/main.go --queue "Helpdesk" --status "open" --page 2 --per-page 20
```

## View Details

```bash
go run examples/ticket_search/main.go --queue "Helpdesk" --subject "printer" --details-id 12345
```
