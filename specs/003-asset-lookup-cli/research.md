# Research: Asset Lookup CLI

**Branch**: `003-asset-lookup-cli`
**Date**: 2026-02-05

## Unknowns & Clarifications

### 1. AssetSQL Syntax for Custom Fields
- **Question**: What is the correct syntax for querying a Custom Field by name in AssetSQL?
- **Finding**: The standard RT syntax is `'CF.{FieldName}' = 'Value'`. Note that curly braces are often required if the name contains spaces (e.g., "Internal Name").
- **Decision**: Use `fmt.Sprintf("'CF.{%s}' = '%s'", "Internal Name", value)` for the query construction.

### 2. .env File Loading (Zero Dependency)
- **Question**: How to load `.env` files without adding `godotenv` dependency?
- **Finding**: A simple parser can read the file line-by-line using `bufio.Scanner`, ignore comments (`#`), split by `=`, and set `os.Setenv`.
- **Decision**: Implement `loadEnv(filename string)` helper function in `main.go`.

### 3. Output Formatting
- **Question**: How to align multiple Key-Value pairs nicely?
- **Finding**: `text/tabwriter` is perfect for this. We can use it with a tabstop of 0 and padding of 2.
- **Decision**: Use `text/tabwriter` for both the summary listing (table) and the detail view (aligned key-value pairs).

## Technology Decisions

| Area | Choice | Rationale |
|------|--------|-----------|
| **CLI Flags** | `flag` (std lib) | Sufficient for 3 mutually exclusive flags; no need for `cobra`/`urfave` deps. |
| **Env Loading** | Custom implementation | Adheres to "Zero-dependency" constitution constraint. |
| **Parsing** | `json.Number` | Existing `Asset` struct already handles this safely. |

## Alternatives Considered

- **`godotenv`**: Rejected to avoid adding a dependency to the core `go.mod` (or managing a separate module) for a simple example.
- **Interactive Prompts for Multiple Matches**: Rejected to keep the CLI stateless and pipe-friendly (per spec).
