# Data Model: Ticket Search Example

## Entities

### Ticket

- **Description**: Support ticket record surfaced in search results and detail view.
- **Key Fields**: ID, Subject, Status, Queue, Owner, Requestor, URL (`_url`).
- **Relationships**: Appears in `SearchResult.Items`.

### SearchResult

- **Description**: Paged list of search results.
- **Fields**: Items (Ticket[]), Page, PerPage, Total, NextPageURL.

### SearchFilters

- **Description**: User-provided filters for building TicketSQL.
- **Fields**: Queue, Status, Owner, Requestor, Subject.
- **Validation**:
  - Unknown filters are rejected.
  - Empty filter set is invalid.
  - Subject allows keyword match; Queue/Status/Owner/Requestor are exact match.

### Pagination

- **Description**: Paging controls for search.
- **Fields**: Page, PerPage.
- **Validation**:
  - Page defaults to 1 if not provided.
  - PerPage defaults to 20 if not provided.
  - PerPage maximum is 100.

## State Transitions

- Not applicable (read-only search example).
