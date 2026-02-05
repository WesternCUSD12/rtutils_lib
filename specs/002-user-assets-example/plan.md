# Implementation Plan: Example Program - User Asset Lookup

**Branch**: `002-user-assets-example` | **Date**: 2026-02-04 | **Spec**: [specs/002-user-assets-example/spec.md](spec.md)
**Input**: Feature specification from `specs/002-user-assets-example/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Create a standalone CLI example program in `examples/user_assets/main.go` that authenticates with Request Tracker and searches for assets. It must support searching by Username, Email, or Real Name via mutually exclusive CLI flags, querying both `Owner` and `HeldBy` relationships, and outputting results in a table format.

## Technical Context

**Language/Version**: Go 1.25+
**Primary Dependencies**: `rtutils_lib` (internal), `os`, `flag`, `fmt`, `log`, `text/tabwriter`
**Storage**: N/A
**Testing**: Manual verification via `go run` against a live or mocked RT instance.
**Target Platform**: CLI (Cross-platform)
**Project Type**: Example Program
**Performance Goals**: N/A
**Constraints**: Must use standard library for CLI args parsing (`flag` package).
**Scale/Scope**: Single main file.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

*   **Idiomatic Go**: [PASS] Use `flag` package for CLI arguments.
*   **Interface-First**: [N/A] Consumer code.
*   **Test-First**: [PASS] Acceptance scenarios defined.
*   **Type Safety**: [PASS] Go strong typing.
*   **Error Transparency**: [PASS] Errors logged to stderr.

## Project Structure

### Documentation (this feature)

```text
specs/002-user-assets-example/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output
```

### Source Code (repository root)

```text
examples/
└── user_assets/
    └── main.go          # Main entry point for the example
```

**Structure Decision**: Single file example.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
