# Implementation Plan: Built-in Query Methods for RT Data Types

**Branch**: `004-built-in-queries` | **Date**: February 5, 2026 | **Spec**: [specs/004-built-in-queries/spec.md](specs/004-built-in-queries/spec.md)
**Input**: Feature specification from `/specs/004-built-in-queries/spec.md`

**Note**: This plan follows the `/speckit.plan` workflow and the rtutils_lib constitution.

## Summary

Add built-in query methods for assets, tickets, and users with explicit exact/partial variants, preserve custom query support, and document query construction. Implement interface updates, new methods, tests, and documentation while returning only the first page of results.

## Technical Context

**Language/Version**: Go 1.25.4  
**Primary Dependencies**: Standard library; tests use `testify` and `httpmock`  
**Storage**: N/A (client library)  
**Testing**: `go test ./...` with `testify/assert` or `testify/require`, `httpmock`  
**Target Platform**: Go library targeting RT REST 2.0  
**Project Type**: Single library module  
**Performance Goals**: First page of search results returned in under 2 seconds for typical RT instances  
**Constraints**: Zero-dependency runtime; interface-first; TDD required; APIError wrapping for HTTP failures  
**Scale/Scope**: Library additions to Assets/Tickets/Users services and contracts

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- **I. Idiomatic Go**: PASS - Methods are Go-style with `context.Context` and `(Result, error)` signatures.
- **II. Interface-First Design**: PASS - Contracts in `contracts/` will be updated before implementation.
- **III. Test-First (NON-NEGOTIABLE)**: PASS - Plan includes contract tests before implementation.
- **IV. Type Safety & Validation**: PASS - Concrete structs used; validation required before requests.
- **V. Error Transparency**: PASS - APIError wrapping maintained; add explicit custom field error.

## Project Structure

### Documentation (this feature)

```text
specs/004-built-in-queries/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
asset.go
asset_test.go
client.go
contracts/
  interfaces.go
ticket.go
ticket_test.go
types.go
user.go
user_test.go
```

**Structure Decision**: Single Go module with root-level services and tests.

## Phase 0: Research

- Produce [specs/004-built-in-queries/research.md](specs/004-built-in-queries/research.md) with decisions on query syntax, matching variants, pagination behavior, and error handling.

## Phase 1: Design & Contracts

- Produce [specs/004-built-in-queries/data-model.md](specs/004-built-in-queries/data-model.md) documenting Asset, Ticket, User, and SearchResult structures and validation rules.
- Produce OpenAPI contracts in [specs/004-built-in-queries/contracts](specs/004-built-in-queries/contracts) describing the logical query operations.
- Produce [specs/004-built-in-queries/quickstart.md](specs/004-built-in-queries/quickstart.md) with examples for built-in methods and custom queries.
- Run agent context update script.

**Post-Design Constitution Check**: PASS - No violations introduced.

## Phase 2: Planning

- `/speckit.tasks` will generate tasks.md after this plan.

## Complexity Tracking

No constitution violations.
