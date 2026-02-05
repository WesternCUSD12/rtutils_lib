# Implementation Plan: Asset Lookup CLI Tool

**Branch**: `003-asset-lookup-cli` | **Date**: 2026-02-05 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/003-asset-lookup-cli/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a command-line interface (CLI) tool in `examples/asset_lookup/` to query Request Tracker assets. The tool will support looking up assets by specific ID, Name, or "Internal Name" (Custom Field). It will output full asset details in a vertical key-value format for readability and handle multiple matches by displaying a summary table. Configuration will be handled via environment variables with support for a local `.env` file.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: `rtutils_lib` (internal), standard library (`flag`, `fmt`, `os`, `text/tabwriter`, `bufio`, `strings`).
**Storage**: N/A
**Testing**: Manual CLI verification; unit tests for query builder logic if extracted.
**Target Platform**: CLI
**Project Type**: CLI Tool Example
**Performance Goals**: <2s response time.
**Constraints**: Zero external dependencies (will implement simple .env parser).
**Scale/Scope**: Single file or small package example.

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

- [x] **Idiomatic Go**: Uses standard library `flag` and `text/tabwriter`.
- [x] **Interface-First**: Leverages existing `AssetService` interface.
- [x] **Test-First**: Will verify core logic; primary validation is via CLI execution.
- [x] **Type Safety**: Uses typed structs for asset data.
- [x] **Error Transparency**: Reports errors to stderr with exit codes.
- [x] **Zero-dependency**: No external libs for .env or table formatting.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
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
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |
