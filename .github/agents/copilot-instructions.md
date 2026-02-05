# rtutils_lib Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-02-04

## Active Technologies
- Go 1.25+ + `rtutils_lib` (internal), `os`, `fmt`, `log` (002-user-assets-example)
- Go 1.25+ + `rtutils_lib` (internal), `os`, `flag`, `fmt`, `log`, `text/tabwriter` (002-user-assets-example)
- Go 1.25+ + `rtutils_lib` (internal), standard library (`flag`, `fmt`, `os`, `text/tabwriter`, `bufio`, `strings`). (003-asset-lookup-cli)
- Go 1.22+ + Standard library for runtime; `github.com/stretchr/testify` for tests only (001-ticket-search)
- Go 1.25.4 + Standard library; tests use `testify` and `httpmock` (004-built-in-queries)
- N/A (client library) (004-built-in-queries)

- Go 1.22+ + None (Standard Library for HTTP/JSON). `github.com/stretchr/testify` for testing. (001-rt-api-wrapper)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.22+

## Code Style

Go 1.22+: Follow standard conventions

## Recent Changes
- 004-built-in-queries: Added Go 1.25.4 + Standard library; tests use `testify` and `httpmock`
- 001-ticket-search: Added Go 1.22+ + Standard library for runtime; `github.com/stretchr/testify` for tests only
- 003-asset-lookup-cli: Added Go 1.25+ + `rtutils_lib` (internal), standard library (`flag`, `fmt`, `os`, `text/tabwriter`, `bufio`, `strings`).


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
