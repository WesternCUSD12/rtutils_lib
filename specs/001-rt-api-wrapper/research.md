# Research: RT API Wrapper

**Topic**: RT REST 2.0 API & Go Implementation Strategy
**Date**: 2026-02-04

## 1. Hypermedia & `_url` Handling

The Request Tracker REST 2.0 API uses the `_url` field as a hypermedia control to uniquely identify resources and provide their canonical URI.

### Usage Pattern

- **Search Results**: When searching (e.g., `/tickets?query=...`), the returned list often contains summary objects. The `_url` field points to the full resource endpoint.
- **Creation**: A successful POST returns a JSON object with `id` and `_url`.
- **Navigation**: To get "more information" about an object (as requested), the client must perform a GET request against the value of `_url`.

### Implementation Decision

- **Struct Field**: All main resources (`Ticket`, `User`, `Asset`) will include a `URL string `json:"\_url,omitempty"` field.
- **Resource Expansion**: Implementing a `Fetch()` or `Load()` method on the resource struct/service that uses this URL is the correct pattern to "get more information".
- **ID vs URL**: While `id` is used for user-friendly referencing, `_url` is the precise machine API endpoint. The library will effectively abstract this by constructing URLs from IDs for initial requests but using `_url` from responses for subsequent actions if available.

## 2. Authentication

RT 6.0+ uses Token Authentication.
**Header**: `Authorization: token <token_value>`
**Decision**: The `Client` struct will hold the token and inject this header into every `http.Request`.

## 3. Data Structures & Custom Fields

RT objects are mix of static fields (Subject, Status) and dynamic `CustomFields`.
**Decision**:

- Static fields will be mapped to strict struct fields (e.g., `Type string`, `Subject string` for Tickets).
- `CustomFields` will be mapped to `map[string]interface{}` or a dedicated `map[string]CustomFieldValue` type to handle the loose schema of custom fields.
- **JSON Parsing**: `encoding/json` standard library is sufficient.

## 4. Error Handling

RT REST 2.0 errors can be returned as:

- Status codes (401, 404, etc.)
- JSON bodies with `message` or array of strings.

**Decision**:

- A generic `APIError` struct will be defined to capture non-2xx response bodies.
- The client will check `resp.StatusCode` and try to decode the body into `APIError` if not 200/201.

## 5. Alternatives Considered

- **Code Generation (OpenAPI)**: RT doesn't provide a rigorous OpenAPI spec that is easy to consume for Go without heavy modification. Native structs are cleaner for this scope.
- **Third-party Libraries**: None exist that are current for RT 6.0 REST 2.0 in Go. Building native is the correct path.
