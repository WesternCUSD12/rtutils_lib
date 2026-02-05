# Research & Findings: RT Integration Test Suite

**Date**: February 5, 2026  
**Phase**: Phase 0 (Planning/Research)  
**Status**: Complete

## Executive Summary

This document consolidates research on building an integration test suite for the rtutils library. All major unknowns have been clarified through specification review and existing codebase analysis.

## Read-Only Method Inventory

### TicketService (4 read-only methods)

1. **Get(ctx context.Context, id string) (*Ticket, error)**
   - Fetches a single ticket by ID
   - Essential for verifying ticket retrieval works

2. **GetByURL(ctx context.Context, url string) (*Ticket, error)**
   - Fetches a ticket using full RT URL from search results
   - Tests URL-based retrieval mechanism

3. **Search(ctx context.Context, query string, page int, perPage int) (*SearchResult[Ticket], error)**
   - Searches tickets using TicketSQL with pagination
   - Core query functionality; expands summary results to full details

4. **SearchBySubject(ctx context.Context, query string) (*SearchResult[Ticket], error)**
   - Specialized search by subject with fuzzy matching
   - Tests convenience method and string escaping (single quotes)

### UserService (7 read-only methods)

1. **Get(ctx context.Context, id string) (*User, error)**
   - Fetches user by ID or Name
   - Basic user retrieval

2. **Search(ctx context.Context, query string) (*SearchResult[User], error)**
   - Base search method for users
   - Foundation for specialized searches

3. **SearchByUsernameExact(ctx context.Context, username string) (*SearchResult[User], error)**
   - Exact username match search

4. **SearchByUsernamePartial(ctx context.Context, query string) (*SearchResult[User], error)**
   - Partial username match (LIKE operator)

5. **SearchByEmailExact(ctx context.Context, email string) (*SearchResult[User], error)**
   - Exact email address search

6. **SearchByEmailPartial(ctx context.Context, query string) (*SearchResult[User], error)**
   - Partial email address search

7. **SearchByNameExact/Partial, RealName variants**
   - Full name search with exact and partial matching
   - Tests multiple query field types

### AssetService (5+ read-only methods)

1. **Get(ctx context.Context, id string) (*Asset, error)**
   - Fetches asset by ID
   - Basic asset retrieval

2. **Search(ctx context.Context, query string) (*SearchResult[Asset], error)**
   - Base asset search with automatic expansion of summary results
   - Detailed result fetching via `_url`

3. **SearchByNameExact(ctx context.Context, name string) (*SearchResult[Asset], error)**
   - Exact asset name match

4. **SearchByNamePartial(ctx context.Context, query string) (*SearchResult[Asset], error)**
   - Partial asset name search

5. **SearchByCustomFieldExact(ctx context.Context, fieldName string, value string) (*SearchResult[Asset], error)**
   - Custom field search (dynamic field support)

**Total Read-Only Methods**: ~20 methods across 3 services

## Test Data Requirements

### Current Assumption (Phase 1 MVP)
- No setup or teardown of test data
- All tests run against existing live RT data
- Tests are read-only queries only
- Assumes RT instance has:
  - At least 5+ tickets (various statuses)
  - At least 3+ users (email + name fields populated)
  - At least 2+ assets (with names)

### Future Phases (P2+)
- Could add fixture setup/teardown for deterministic testing
- Could create dedicated test database snapshot
- Could support environment-specific test configurations

## Validator Design

### Entity-Specific Validation Rules

Based on clarification Q5 ("usefulness validation"), validators check for minimum viable results:

#### Ticket Validation
```
Required fields: id (non-empty), status (non-empty)
Plus at least one additional field: subject, description, queue, owner, etc.
Passes if: id ≠ null && status ≠ null && len(other_fields) > 0
```

#### User Validation
```
Required fields: id (non-empty), name (non-empty)
Passes if: id ≠ null && name ≠ null
```

#### Asset Validation
```
Required fields: id (non-empty), name (non-empty)
Passes if: id ≠ null && name ≠ null
```

#### Generic Validation
```
Fallback for any other query type:
Required: id (non-empty) + at least one other field
Passes if: id ≠ null && len(other_fields) > 0
```

## Error Classification Strategy

Errors encountered during test execution must be classified for meaningful reporting:

### Error Categories

1. **Infrastructure Error**
   - Connection failed, RT instance unreachable
   - Authentication failed (invalid credentials)
   - Network timeout
   - TLS/certificate errors
   - **Action**: Report as infrastructure issue, don't count as test failure

2. **Test Assertion Failure**
   - Method executed successfully but returned unexpected data
   - Missing required fields in result
   - Wrong data type in response
   - Empty result when data expected
   - **Action**: Fail test, include expected vs actual in report

3. **Method Error**
   - Method execution returned error (not nil)
   - Method panicked
   - Unhandled exception during method call
   - **Action**: Fail test, include error message/stack trace

## JSON Report Schema

The test suite will output a JSON report following this schema:

```json
{
  "execution": {
    "timestamp": "2026-02-05T14:23:45Z",
    "duration_seconds": 187.5,
    "rt_instance": "https://rt.example.com",
    "total_tests": 20,
    "passed": 18,
    "failed": 2,
    "infrastructure_errors": 0,
    "assertion_failures": 2,
    "method_errors": 0
  },
  "results": [
    {
      "test_id": "ticket-get-1",
      "method": "TicketService.Get",
      "status": "PASS",
      "duration_ms": 145,
      "result_summary": "Retrieved ticket #123 with 8 fields"
    },
    {
      "test_id": "user-search-by-email-2",
      "method": "UserService.SearchByEmailPartial",
      "status": "FAIL",
      "error_type": "ASSERTION_FAILURE",
      "duration_ms": 234,
      "error_message": "Expected 'name' field to be non-empty, got empty string",
      "expected": { "name": "non-empty" },
      "actual": { "name": "" }
    }
  ],
  "summary": "18 of 20 tests passed. 2 assertion failures: user methods returned incomplete data."
}
```

## Testing Framework Choice

**Recommendation**: Use Go's built-in `testing` package with `testify/assert`

**Why**:
- Standard Go tool, no external runtime dependencies
- Already used by rtutils_lib
- Aligns with constitution requirement for zero-dependency runtime
- `testify` permitted for `_test.go` files only per constitution
- Clear assertion syntax for validator tests

## Sequential Execution Model

**Rationale** (from Clarification Q2):
- Simpler initial implementation
- Deterministic debugging (no race conditions)
- Clear error propagation
- ~20 tests * ~200-300ms per test = ~4-6 minutes total
- Fits within 5-minute success criterion

**Future Optimization**: Parallel execution can be added as P2 enhancement if needed.

## Configuration Approach

**Phase 1 MVP**: Load RT connection details from environment variables

Required environment variables:
```
RT_URL=https://rt.example.com
RT_USERNAME=admin
RT_PASSWORD=secret123
```

Optional:
```
RT_TIMEOUT=30          # seconds per request
RT_DEBUG=false         # verbose logging
```

**Phase 2+**: Support JSON config files and CLI flags for flexibility

## Success Metrics (Measurable)

Based on spec SC-001 through SC-007:

- ✅ All ~20 read-only methods can be tested in single execution
- ✅ Validator rules provide entity-specific meaningful data checks
- ✅ Full test suite completes in <5 minutes
- ✅ JSON report clearly shows method-by-method results
- ✅ Exit codes enable CI/CD integration (0=success, non-zero=failure)
- ✅ Error messages provide sufficient diagnostic detail

## Dependencies & Integration Points

### Library Dependencies
- `rtutils_lib` - the library being tested (imported directly)
- Standard Go libraries only for runtime
- `testify` for validator unit tests

### Integration Points
- Expects `Client` struct and service interfaces from rtutils_lib
- Requires live RT instance with typical data
- No modifications to existing rtutils_lib code needed

### CI/CD Integration
- Tool should be runnable as: `./rtutils-test-suite --config rt-config.json --report report.json`
- Exit code 0 on all tests pass, non-zero on any failure
- JSON report parseable by CI/CD dashboards

## Open Questions for Phase 2

1. Per-method timeout values - should these be configurable?
2. Should test suite support multiple RT instances in sequence?
3. How to handle custom fields on Tickets/Users/Assets?
4. Should validators be pluggable/extensible?
5. For P2: How to support test data setup/teardown?

## Conclusion

Research phase provides clear direction on:
- ✅ Scope: ~20 read-only methods to test
- ✅ Validation: Entity-specific required field checking
- ✅ Execution: Sequential model selected
- ✅ Reporting: JSON schema defined with error classification
- ✅ Integration: Environment-based configuration approach
- ✅ Testing: Go testing + testify following constitution

**Ready for Phase 1 design work.**
