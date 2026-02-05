# Implementation Plan: Ticket Search Example

**Branch**: `001-ticket-search` | **Date**: February 5, 2026 | **Spec**: [specs/001-ticket-search/spec.md](spec.md)
**Input**: Feature specification from `/specs/001-ticket-search/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Create a Go example under `examples/` that demonstrates ticket search by Queue, Status, Owner, Requestor, and Subject keywords with manual pagination. The example uses `rtutils_lib` to build TicketSQL queries, retrieves paged `SearchResult` data, and hydrates ticket details via the RT `_url` link.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.22+  
**Primary Dependencies**: Standard library for runtime; `github.com/stretchr/testify` for tests only  
**Storage**: N/A  
**Testing**: `go test ./...` with `testify` for unit tests  
**Target Platform**: CLI example on macOS/Linux/Windows  
**Project Type**: Single library with examples  
**Performance Goals**: Search responses under 2 seconds for typical queries  
**Constraints**: Zero runtime deps; enforce page size max 100  
**Scale/Scope**: Typical RT instances (hundreds to thousands of tickets)

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- ✅ Idiomatic Go: use stdlib, context-first, gofmt
- ✅ Interface-first: update `contracts/` before implementation if signatures change
- ✅ Test-first: add contract/unit tests before implementation changes
- ✅ Type safety & validation: validate inputs before request
- ✅ Error transparency: return `APIError` details to caller

### Post-Design Re-check

- ✅ No deviations introduced by research/design artifacts

## Project Structure

### Documentation (this feature)

```text
specs/001-ticket-search/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
contracts/
examples/
├── asset_lookup/
├── user_assets/
└── ticket_search/            # New example (planned)
specs/
├── 001-rt-api-wrapper/
├── 001-ticket-search/
├── 002-user-assets-example/
└── 003-asset-lookup-cli/
ticket.go
client.go
types.go
```

**Structure Decision**: Single Go module with example programs under `examples/`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |
