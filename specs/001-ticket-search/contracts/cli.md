# Contract: Ticket Search Example CLI

## Command

```text
go run examples/ticket_search/main.go [flags]
```

## Flags

- `--queue` (string, optional): Exact match for Queue
- `--status` (string, optional): Exact match for Status
- `--owner` (string, optional): Exact match for Owner
- `--requestor` (string, optional): Exact match for Requestor
- `--subject` (string, optional): Keyword match for Subject
- `--page` (int, optional): Page number (default 1)
- `--per-page` (int, optional): Results per page (default 20, max 100)
- `--details-id` (string, optional): Ticket ID to hydrate and display full details

## Behavior

- At least one of `--queue`, `--status`, `--owner`, `--requestor`, or `--subject` is required.
- Filters are combined with AND semantics.
- Invalid filters or invalid pagination values produce a validation error and non-zero exit.
- Search results are displayed in a table (ID, Subject, Status, Queue, Owner, Requestor, URL).
- If `--details-id` is provided, the example fetches and displays full details for that ticket using the `_url` field when available.

## Exit Codes

- `0`: Success
- `1`: Validation or request error
- `2`: No results or ticket not found
