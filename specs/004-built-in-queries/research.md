# Research: Built-in Query Methods for RT Data Types

## Decision 1: Query Syntax Per Service

- Decision: Use RT REST 2.0 JSON search for assets, TicketSQL for tickets, and RT user query syntax for users.
- Rationale: Matches existing service implementations and preserves compatibility with current RT API behavior.
- Alternatives considered: Unifying all searches under a custom query DSL; rejected due to added complexity and divergence from RT API semantics.

## Decision 2: Exact vs Partial Matching Variants

- Decision: Provide explicit Exact and Partial methods for asset name, asset custom fields, and user identifiers.
- Rationale: Removes ambiguity, aligns with clarified requirements, and keeps method behavior predictable.
- Alternatives considered: Boolean flags or single methods; rejected due to unclear semantics and discoverability.

## Decision 3: Pagination Behavior

- Decision: Built-in query methods return only the first page of results and rely on Search/SearchWithCriteria for additional pages.
- Rationale: Avoids timeouts and memory spikes for large result sets while preserving pagination metadata.
- Alternatives considered: Auto-fetch all pages; rejected due to performance and latency risks.

## Decision 4: Custom Field Not Found Errors

- Decision: Return a specific error when a custom field does not exist in the RT instance.
- Rationale: Distinguishes between invalid query configuration and valid queries with no results.
- Alternatives considered: Return empty results; rejected due to loss of diagnostic clarity.

## Decision 5: No New Runtime Dependencies

- Decision: Keep runtime dependencies to the Go standard library; use testify and httpmock only for tests.
- Rationale: Aligns with repository constitution and existing dependency policy.
- Alternatives considered: Adding a query builder library; rejected due to policy and unnecessary complexity.
