# Implementation Plan: RT Integration Test Suite

**Branch**: `005-integration-test-suite` | **Date**: February 5, 2026 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/005-integration-test-suite/spec.md`

## Summary

Build a robust CLI-based integration test suite that validates all read-only methods across the rtutils library against a live Request Tracker instance. The suite will execute ~20 query methods sequentially, validate results contain meaningful data per entity type (Tickets/Users/Assets), and generate JSON reports with error classification for CI/CD integration.

**Scope - Phase 1 MVP (Read-Only)**:

- Test TicketService methods: Get, GetByURL, Search, SearchBySubject
- Test UserService methods: Get, Search, SearchByUsername/Email/Name (exact and partial)
- Test AssetService methods: Get, Search, SearchByName/CustomField variants
- Smart validation: Entity-specific required fields (Tickets: id+status, Users/Assets: id+name)
- Sequential execution with detailed error classification (Infrastructure/Assertion/Method errors)
- JSON report output with pass/fail counts and execution metadata

## Technical Context

**Language/Version**: Go 1.22+ (matches rtutils_lib)
**Primary Dependencies**:

- `rtutils_lib` (the library being tested - zero external runtime deps)
- `testify/assert` and `testify/require` (testing only - aligns with constitution)
- Standard library: `encoding/json`, `net/http`, `context`

**Storage**: N/A (read-only queries to live RT instance)
**Testing**: Go `testing` package + `testify` (unit tests for validators; integration tests against live RT)
**Target Platform**: CLI tool executable on Linux/macOS (any platform Go supports)
**Project Type**: Single CLI tool with companion library
**Performance Goals**: Complete full test suite in <5 minutes against typical RT instance (~20 read methods)
**Constraints**:

- Read-only operation (no data modification)
- Sequential test execution
- <5 minute total execution time
- Per-method timeouts TBD during implementation

**Scale/Scope**:

- ~20 read-only methods to test (across 3 services)
- Configurable RT instance (host, credentials via env/config)
- Extensible test case structure for future custom tests (P3 feature)

## Constitution Check

✅ **GATE PASSED** - Implementation aligns with all constitution principles:

| Principle | Check | Status |
|-----------|-------|--------|
| **I. Idiomatic Go** | Uses standard lib (`net/http`, `encoding/json`), `context.Context` for IO, `gofmt` | ✅ |
| **II. Interface-First** | Tests use rtutils library's public interface (Client, Services) | ✅ |
| **III. Test-First (TDD)** | Validator functions will be tested before implementation; integration tests follow real query execution | ✅ |
| **IV. Type Safety** | Assertion validators use concrete Ticket/User/Asset structs; no map[string]interface{} | ✅ |
| **V. Error Transparency** | Test suite captures and reports structured APIErrors with classification | ✅ |
| **Constraint: Zero-dependency** | Only runtime dependencies are rtutils and stdlib; testify only in `_test.go` | ✅ |
| **Constraint: TDD Workflow** | Plan includes Phase 0 research, Phase 1 design, then Phase 2 implementation | ✅ |

**No violations detected.** Implementation can proceed.

## Project Structure

### Documentation (this feature)

```text
specs/005-integration-test-suite/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file (Phase 1 output)
├── research.md          # Phase 0 research findings (to be generated)
├── data-model.md        # Phase 1 design: entities, validators (to be generated)
├── quickstart.md        # Phase 1: setup and usage guide (to be generated)
├── contracts/           # Phase 1: test case definitions (to be generated)
│   ├── README.md        # Contract overview
│   ├── test-cases.json  # Structured test case definitions
│   └── validators.md    # Validator logic & rules
└── checklists/
    └── requirements.md  # Spec quality checklist
```

### Source Code (repository root)

```text
.
├── cmd/
│   └── rtutils-test-suite/           # Phase 2: CLI tool
│       └── main.go                   # Entry point
├── internal/
│   ├── testrunner/                   # Phase 2: Test execution engine
│   │   ├── runner.go                 # Orchestrates test execution
│   │   ├── reporter.go               # JSON report generation
│   │   └── runner_test.go            # Unit tests
│   ├── validators/                   # Phase 2: Smart validation for results
│   │   ├── ticket_validator.go       # Ticket result validation
│   │   ├── user_validator.go         # User result validation
│   │   ├── asset_validator.go        # Asset result validation
│   │   ├── validators.go             # Base validator interface
│   │   └── validators_test.go        # Unit tests
│   └── rtconfig/                     # Phase 2: RT instance configuration
│       ├── config.go                 # Load RT connection params
│       └── config_test.go            # Unit tests
├── examples/
│   ├── rt-instance-config.json       # Example RT config file
│   └── test-report-sample.json       # Example report output

# Existing test infrastructure
├── ticket_test.go                    # Existing unit tests (unchanged)
├── user_test.go                      # Existing unit tests (unchanged)
├── asset_test.go                     # Existing unit tests (unchanged)
├── test_helpers_test.go              # Existing test utilities (unchanged)
├── types.go                          # Data models (existing)
├── client.go                         # HTTP client (existing)
└── contracts/
    └── interfaces.go                 # Existing service interfaces
```

**Structure Decision**: 
- CLI tool located at `cmd/rtutils-test-suite/` following Go conventions
- Core testing logic in `internal/testrunner/` and `internal/validators/`
- Configuration in separate `internal/rtconfig/` package for clean separation
- No modifications to existing source files or tests
- Aligns with rtutils_lib's existing interface-first design in `contracts/`

## Implementation Roadmap

### Phase 0: Research (Complete)
- [x] Identify all read-only methods across TicketService, UserService, AssetService
- [x] Document validator requirements per entity type
- [x] Define JSON report schema

### Phase 1: Design & Contracts (In Progress)
- [ ] Generate `data-model.md` with TestCase, TestResult, TestReport models
- [ ] Define validator rules in `contracts/validators.md`
- [ ] Create `contracts/test-cases.json` with structured test definitions
- [ ] Write `quickstart.md` with setup instructions
- [ ] Update agent context with technology decisions
- **Deliverable**: design documents, contracts, quickstart guide

### Phase 2: Implementation (Next - use `/speckit.tasks`)
- [ ] Implement validator functions (ticket/user/asset) with unit tests
- [ ] Build TestRunner with sequential execution logic
- [ ] Implement JSON report generator with error classification
- [ ] Create RTConfig to load credentials and RT connection params
- [ ] Build CLI entry point and argument parsing
- [ ] Integration testing against live RT instance
- **Deliverable**: Executable CLI tool ready for use

## Clarifications Applied

The following clarifications from the specification session have been incorporated:
- **Report Format**: JSON structured output (ensures CI/CD integration capability)
- **Test Execution**: Sequential execution (simpler, deterministic debugging)
- **Test Data**: Use existing live RT data with read-only access (no setup complexity)
- **Error Handling**: Classify all errors (Infrastructure/Assertion/Method errors for visibility)
- **Usefulness Criteria**: Content-based validation per entity type (Tickets: id+status+other, Users/Assets: id+name)

| Violation                  | Why Needed         | Simpler Alternative Rejected Because |
| -------------------------- | ------------------ | ------------------------------------ |
| [e.g., 4th project]        | [current need]     | [why 3 projects insufficient]        |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient]  |
