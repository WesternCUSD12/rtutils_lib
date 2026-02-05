<!--
  Sync Impact Report:
  - Version: Template -> 1.0.0
  - Modified Principles: All (Initial Definition)
  - Templates Checked: plan.md (Aligned), spec.md (Aligned)
  - Pending: N/A
-->

# rtutils_lib Constitution

## Core Principles

### I. Idiomatic Go

The library MUST use standard Go patterns and the standard library (`net/http`, `encoding/json`) wherever possible. Functions MUST return `(Result, error)`. `context.Context` MUST be the first argument for all IO-bound operations. Code MUST be formatted with `gofmt`.

### II. Interface-First Design

All primary service interactions (Tickets, Users, Assets) MUST be defined as Go interfaces in a dedicated `contracts/` package (or equivalent) before implementation. This ensures testability and clear API boundaries. The `Client` struct implements these interfaces.

### III. Test-First (NON-NEGOTIABLE)

Test Driven Development (TDD) is mandatory. Tests against the `contracts` MUST be written before the implementation. Unit tests MUST use `testify/assert` or `testify/require`. Mocking strategies MUST be used for unit tests to avoid live API calls.

### IV. Type Safety & Validation

Leverage Go's type system to ensure correctness. RT Resources (Tickets, Users) MUST be mapped to concrete structs. `map[string]interface{}` is permitted ONLY for `CustomFields` due to their dynamic nature. Input validation (e.g. required ID) MUST occur before the network call.

### V. Error Transparency

Errors MUST NOT be swallowed. HTTP 4xx/5xx responses MUST be wrapped in a structured `APIError` type that exposes the underlying Request Tracker error message or status code to the consumer.

## Implementation Constraints

**Stack**: Go 1.22+.
**Dependencies**: Zero-dependency policy for the core runtime (HTTP client, JSON parsing). `testify` is allowed for `_test.go` files only.
**Hypermedia**: The library MUST abstract the `_url` mechanics; users interact with IDs, the library handles the URL resolution internally where advantageous.

## Development Workflow

1. **Spec**: Define the feature in `specs/`.
2. **Contract**: Define the interface.
3. **Test**: Write the test case failing.
4. **Implement**: Write the code to pass.
5. **Verify**: Run `go test ./...`.

## Governance

All Pull Requests MUST include tests covering the new functionality. Changes to `contracts/` are considered breaking changes if they alter existing method signatures and require a MAJOR version bump.

**Version**: 1.0.0 | **Ratified**: 2026-02-04 | **Last Amended**: 2026-02-04
