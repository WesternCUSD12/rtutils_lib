# Integration Test Contracts

**Date**: February 5, 2026  
**Phase**: Phase 1 (Design)

## Overview

This directory contains the contract specifications for the integration test suite. These contracts define:

1. **Test Case Definitions**: Structured JSON defining which methods to test and how
2. **Validator Rules**: Requirements for what constitutes "useful" results
3. **Expected Behaviors**: Standard response patterns and edge cases

## Files

### test-cases.json

Structured test case definitions for all ~20 read-only methods across the library.

**Structure**:
```json
{
  "test_cases": [
    {
      "id": "ticket-get-1",
      "method": "TicketService.Get",
      "service": "Ticket",
      "operation": "Get",
      "description": "Fetch a single ticket by ID",
      "inputs": {
        "id": "1"
      },
      "expected_result_type": "Ticket",
      "validator": "ticket",
      "enabled": true
    }
  ]
}
```

**Purpose**: 
- Machine-parseable test definitions
- Enables dynamic test execution without hardcoding test cases
- Supports future feature additions (custom tests, parametrized tests)
- Allows runtime configuration of which tests to run

### validators.md

Detailed validator specifications and implementation rules.

**Contents**:
- Validator interface contract
- Entity-specific validation rules
- Required vs optional field definitions
- Error classification rules
- Edge case handling

**Purpose**:
- Design specification for Phase 2 implementation
- Reference for validator unit tests
- Documentation of "useful data" criteria

## Test Case Coverage

### TicketService

| TestID | Method | InputExample | ValidatorType |
|--------|--------|-------------|---|
| ticket-get-1 | Get | {id: "1"} | ticket |
| ticket-getbyurl-1 | GetByURL | {url: "/ticket/1"} | ticket |
| ticket-search-1 | Search | {query: "Status='new'", page: 1, perPage: 20} | ticket |
| ticket-search-subject-1 | SearchBySubject | {query: "bug"} | ticket |

### UserService

| TestID | Method | InputExample | ValidatorType |
|--------|--------|-------------|---|
| user-get-1 | Get | {id: "admin"} | user |
| user-search-1 | Search | {query: "Name='admin'"} | user |
| user-search-un-exact-1 | SearchByUsernameExact | {username: "admin"} | user |
| user-search-un-partial-1 | SearchByUsernamePartial | {query: "adm"} | user |
| user-search-email-exact-1 | SearchByEmailExact | {email: "admin@rt.example.com"} | user |
| user-search-email-partial-1 | SearchByEmailPartial | {query: "@rt.example"} | user |
| user-search-name-exact-1 | SearchByNameExact | {name: "Administrator"} | user |
| user-search-name-partial-1 | SearchByNamePartial | {query: "admin"} | user |

### AssetService

| TestID | Method | InputExample | ValidatorType |
|--------|--------|-------------|---|
| asset-get-1 | Get | {id: "1"} | asset |
| asset-search-1 | Search | {query: "Name LIKE '%server%'"} | asset |
| asset-search-name-exact-1 | SearchByNameExact | {name: "Server-01"} | asset |
| asset-search-name-partial-1 | SearchByNamePartial | {query: "server"} | asset |
| asset-search-custom-exact-1 | SearchByCustomFieldExact | {fieldName: "Location", value: "DC1"} | asset |
| asset-search-custom-partial-1 | SearchByCustomFieldPartial | {fieldName: "Location", value: "DC"} | asset |

**Total**: ~20 test cases

## Validator Contracts

### TicketValidator

```
Input: Ticket struct (or SearchResult[Ticket])
Required Fields:
  - id: non-empty string
  - status: non-empty string
Additional Fields:
  - At least one of: subject, description, queue, owner, priority, created, etc.
Passes if: id ≠ "" AND status ≠ "" AND len(additional_fields) > 0
```

### UserValidator

```
Input: User struct (or SearchResult[User])
Required Fields:
  - id: non-empty string
  - name: non-empty string
Passes if: id ≠ "" AND name ≠ ""
```

### AssetValidator

```
Input: Asset struct (or SearchResult[Asset])
Required Fields:
  - id: non-empty string
  - name: non-empty string
Passes if: id ≠ "" AND name ≠ ""
```

### GenericValidator

```
Input: Any result
Required Fields:
  - id: non-empty (required)
Additional Fields:
  - At least one other field present
Passes if: id ≠ "" AND len(other_fields) > 0
```

## Error Classification Contract

### Infrastructure Error

**Conditions**:
- Connection refused to RT_URL
- TLS certificate validation failed
- DNS resolution failed
- Authentication failed (401/403 from RT)
- Network timeout
- RT instance unavailable (5xx response)

**Handling**: Report as INFRASTRUCTURE_ERROR, do not count test as failure (system issue, not code issue)

### Assertion Failure

**Conditions**:
- Method executed successfully (no exception)
- Result returned but validation failed
- Required fields missing from result
- Field contains empty/null/default value when non-empty expected

**Handling**: Fail test, include expected vs actual in report

### Method Error

**Conditions**:
- Method panicked/threw exception
- Method returned non-nil error
- Unhandled exception in method call

**Handling**: Fail test, include error message and stack trace

## Test Execution Contract

### Sequential Execution

```
for each TestCase:
  1. Connect to RT instance (if not connected)
  2. Call method with TestCase.inputs
  3. Catch any errors/exceptions → classify as METHOD_ERROR or INFRASTRUCTURE_ERROR
  4. If successful:
     a. Get appropriate Validator
     b. Call Validator.Validate(result)
     c. If invalid → TestResult{ASSERTION_FAILURE, error_message}
     d. If valid → TestResult{PASS}
  5. Record timing, result data, error details
  6. Continue to next TestCase
7. Compile all TestResults into TestReport
8. Output JSON report
9. Return exit code based on results
```

### Timing Guarantees

- Overall suite: <5 minutes
- Per-method: ~200-300ms average (TBD during implementation)
- Report generation: <1 second

## JSON Report Contract

See [data-model.md](../data-model.md#json-report-schema) for the complete JSON schema.

**Requirements**:
- Valid JSON output
- ISO 8601 timestamps
- All numeric values included
- Error messages included for failures
- Expected vs actual for assertion failures
- Human-readable summary at end

## Phase 2 Implementation Checklist

- [ ] Implement TestCase loader from test-cases.json
- [ ] Implement ResultValidator interface
- [ ] Implement TicketValidator with unit tests
- [ ] Implement UserValidator with unit tests
- [ ] Implement AssetValidator with unit tests
- [ ] Implement GenericValidator with unit tests
- [ ] Implement TestRunner with sequential execution
- [ ] Implement error classification logic
- [ ] Implement JSON report generation
- [ ] Implement RTConnection configuration
- [ ] Implement CLI entry point
- [ ] Integration tests against live RT
- [ ] Update this document with Phase 2 completion notes

## Future Extensions

### P2 Features
- [ ] Parallel execution mode
- [ ] Retry logic for transient failures
- [ ] Performance benchmarking
- [ ] Custom validator plugins

### P3 Features
- [ ] Test data setup/teardown fixtures
- [ ] Parametrized test cases
- [ ] Custom assertion rules
- [ ] Dashboard integration

## Questions & Decisions

### Q: Should test cases be version-specific?
**A**: No. Phase 1 targets current/recent RT versions only. Version detection and fallback validators can be added in P2 if needed.

### Q: What about custom ticket/user/asset fields?
**A**: CustomFields in rtutils types are `map[string]interface{}`. Phase 1 validates core fields only; extension for custom field validation can be added in P2.

### Q: Can tests modify RT data?
**A**: No. Phase 1 MVP is read-only. Write operations can be tested in Phase 2 with proper fixture management.

### Q: Should validators be strict or lenient?
**A**: Lenient (minimal required fields). Goal is to catch completely broken methods, not enforce all fields. Additional validation can be added per-deployment in config.

## Contact & References

- **Specification**: [../spec.md](../spec.md)
- **Data Model**: [../data-model.md](../data-model.md)
- **Implementation Plan**: [../plan.md](../plan.md)
- **Quickstart**: [../quickstart.md](../quickstart.md)
