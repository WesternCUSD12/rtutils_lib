# Research: AssetSQL for User Assets

**Status**: Consolidated
**Date**: 2026-02-04

## Unknowns & Clarifications

### 1. AssetSQL Syntax for User Lookup
**Task**: Determine the correct AssetSQL queries for Username, Email, and Name.
**Findings**:
- Relationship: `HeldBy` and `Owner` are both relevant.
- Username: `Owner.Name = 'val' OR HeldBy.Name = 'val'`
- Email: `Owner.EmailAddress = 'val' OR HeldBy.EmailAddress = 'val'`
- Real Name: `Owner.RealName = 'val' OR HeldBy.RealName = 'val'`
- **Decision**: Use the OR logic combining both `Owner` and `HeldBy` for the specific field requested.

### 2. Authentication Environment Variables
**Task**: Confirm standard env var naming for RT tools.
**Findings**:
- Common conventions are `RT_USER`, `RT_PASS`, `RT_auth_token` or `RT_web_base_url`.
- For this library, we have been using `RT_BASE_URL` and `RT_TOKEN` in our own scripts?
- Checking `client.go` or tests: Tests use mock URLs.
- **Decision**: Use `RT_BASE_URL` and `RT_TOKEN` as defined in the spec `FR-002`.

## Technology Decisions

| Decision | Context | Choice | Rationale |
| :--- | :--- | :--- | :--- |
| **CLI Framework** | Arguments | `flag` package | Standard library, supports `--username`, `--email`, `--name` flags out of the box. |
| **Logging** | Output | `log` / `fmt` | Standard library is sufficient. `log` for stderr errors, `fmt` for stdout data. |
| **Formatting** | Output | `text/tabwriter` | Built-in Go table formatter for aligned output. |
