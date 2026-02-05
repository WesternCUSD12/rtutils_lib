# Phase 0 Research: Ticket Search Example

## Decision 1: Ticket search query format

- **Decision**: Use RT TicketSQL with AND-combined filters (e.g., `Queue = 'Helpdesk' AND Status = 'open' AND Subject LIKE 'printer'`).
- **Rationale**: Matches RT search conventions and keeps results predictable for operators.
- **Alternatives considered**: OR-combined filters; user-selectable AND/OR.

## Decision 2: Pagination parameters

- **Decision**: Default to page=1 and per_page=20; enforce per_page <= 100.
- **Rationale**: Provides predictable paging and avoids large payloads.
- **Alternatives considered**: No defaults; higher default page sizes.

## Decision 3: Detail hydration

- **Decision**: Use the `_url` field from search results to retrieve full ticket details.
- **Rationale**: Aligns with RT hypermedia responses and uses existing `Client.request` behavior for full URLs.
- **Alternatives considered**: Reconstructing URLs from IDs only.

## Decision 4: Example output format

- **Decision**: Tabular summary for search results and optional detail view for a selected ticket.
- **Rationale**: Consistent with existing examples and easy to scan.
- **Alternatives considered**: JSON-only output.
