# Data Model: Built-in Query Methods for RT Data Types

## Entities

### Asset

- Fields:
  - ID: string or json.Number (RT asset identifier)
  - Name: string
  - Description: string (optional)
  - Status: string (optional)
  - CustomFields: list of custom field entries (name + values)
- Relationships:
  - Appears as items in SearchResult[Asset]
- Validation:
  - SearchByNameExact/Partial require non-empty name/query
  - SearchByCustomFieldExact/Partial require non-empty fieldName and value/query

### Ticket

- Fields:
  - ID: string
  - Subject: string
  - Status: string (optional)
  - Queue: string (optional)
  - Owner: string (optional)
  - Requestor: list of strings (optional)
  - CustomFields: map of string to dynamic values (optional)
- Relationships:
  - Appears as items in SearchResult[Ticket]
- Validation:
  - SearchBySubject requires non-empty query

### User

- Fields:
  - ID: string
  - Name: string
  - Email: string
  - Username: string
  - Disabled: bool (optional)
- Relationships:
  - Appears as items in SearchResult[User]
- Validation:
  - SearchByUsernameExact/Partial require non-empty username/query
  - SearchByEmailExact/Partial require non-empty email/query
  - SearchByNameExact/Partial require non-empty name/query

### SearchResult<T>

- Fields:
  - Total: int
  - Count: int
  - Page: int
  - Pages: int
  - PerPage: int
  - NextPage: string
  - Items: list of T
- Validation:
  - Built-in query methods return only the first page (RT API default pagination)

### APIError

- Fields:
  - StatusCode: int
  - Message: string
- Notes:
  - Used for HTTP 4xx/5xx and explicit custom field not found errors

## State Transitions

No state transitions are introduced by this feature. Query methods are read-only.
