# Data Model: RT Integration Test Suite

**Date**: February 5, 2026  
**Phase**: Phase 1 (Design)  
**Status**: Complete

## Core Entities

### TestCase

Represents a single test that validates one library method.

```go
type TestCase struct {
    ID          string                  // Unique test identifier (e.g., "ticket-get-1")
    MethodName  string                  // Full method name (e.g., "TicketService.Get")
    ServiceType string                  // Service type: "Ticket" | "User" | "Asset"
    Description string                  // Human-readable test description
    
    // Execution configuration
    Timeout     int                     // Timeout in seconds (TBD during implementation)
    Retryable   bool                    // Whether this test can be retried on failure
    
    // Test parameters
    InputParams map[string]interface{}  // Method input params (e.g., {"id": "123"})
    
    // Validation rules
    ValidatorType string                // "ticket" | "user" | "asset" | "generic"
    ExpectedFields []string             // Fields that must be present in result
}
```

### TestResult

Represents the outcome of running a single test.

```go
type TestResult struct {
    TestID          string              // Reference to TestCase.ID
    MethodName      string              // Full method name tested
    Status          string              // "PASS" | "FAIL"
    ErrorType       string              // Error classification:
                                        //   "INFRASTRUCTURE_ERROR"
                                        //   "ASSERTION_FAILURE"
                                        //   "METHOD_ERROR"
                                        //   "UNKNOWN"
    
    // Timing
    StartTimeUTC    string              // ISO 8601 timestamp
    DurationMs      int                 // Execution time in milliseconds
    
    // Result data
    ResultData      interface{}         // Actual result returned by method (JSON-serializable)
    ResultSummary   string              // Human-readable result summary
    
    // Error details
    ErrorMessage    string              // Error message or assertion failure reason
    ExpectedValue   interface{}         // What was expected (for assertion failures)
    ActualValue     interface{}         // What was actually returned
    StackTrace      string              // Stack trace if applicable (method error)
}
```

### TestReport

Aggregates multiple TestResults with summary statistics.

```go
type TestReport struct {
    Execution ExecutionMetadata        // Execution context & statistics
    Results   []TestResult             // Individual test results
}

type ExecutionMetadata struct {
    Timestamp              string                  // ISO 8601 execution start time
    DurationSeconds        float64                 // Total execution time
    RTInstanceURL          string                  // RT instance tested against
    RTInstanceVersion      string                  // RT version (if available)
    LibraryVersion         string                  // rtutils_lib version
    
    // Statistics
    TotalTests             int                     // Total tests executed
    PassedTests            int                     // Number of passed tests
    FailedTests            int                     // Number of failed tests
    InfrastructureErrors   int                     // Count of infrastructure errors
    AssertionFailures      int                     // Count of assertion failures
    MethodErrors           int                     // Count of method execution errors
    
    // Summary
    SummaryMessage         string                  // Human-readable overall result
    AllTestsPassed         bool                    // true if all tests passed
}
```

### RTConnection

Configuration for connecting to a Request Tracker instance.

```go
type RTConnection struct {
    // Connection details
    URL            string              // RT instance URL (e.g., "https://rt.example.com")
    Username       string              // Username for authentication
    Password       string              // Password for authentication (from env var)
    
    // Optional
    Timeout        int                 // Request timeout in seconds (default: 30)
    InsecureTLS    bool                // Skip TLS verification (for testing only)
    
    // Validation
    Validate() error                   // Verify URL format, credentials present
}
```

### ResultValidator (Interface)

Defines the contract for validating test results.

```go
type ResultValidator interface {
    // Validate checks if a result contains meaningful data
    // Returns (isValid, errorMessage)
    Validate(result interface{}) (bool, string)
    
    // RequiredFields returns the list of fields that must be present
    RequiredFields() []string
}
```

### TicketValidator (implements ResultValidator)

Validates Ticket query results.

```go
type TicketValidator struct{
    // No fields needed - stateless validator
}

// Implementation rules:
// - Required: id (non-empty), status (non-empty)
// - Plus: at least one additional field (subject, description, queue, owner, etc.)
// - Passes if all constraints met
```

### UserValidator (implements ResultValidator)

Validates User query results.

```go
type UserValidator struct {
    // No fields needed - stateless validator
}

// Implementation rules:
// - Required: id (non-empty), name (non-empty)
// - Passes if both present
```

### AssetValidator (implements ResultValidator)

Validates Asset query results.

```go
type AssetValidator struct {
    // No fields needed - stateless validator
}

// Implementation rules:
// - Required: id (non-empty), name (non-empty)
// - Passes if both present
```

### GenericValidator (implements ResultValidator)

Fallback validator for unknown entity types.

```go
type GenericValidator struct {
    // No fields needed - stateless validator
}

// Implementation rules:
// - Required: id (non-empty)
// - Plus: at least one additional field
// - Passes if constraints met
```

## Test Execution Flow

### Data Flow Diagram

```
┌─────────────────┐
│  RTConnection   │
│   (config)      │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│   TestRunner            │
│  (sequential executor)  │
└────────┬────────────────┘
         │
         ├─► For each TestCase:
         │       │
         │       ├─► Execute method via rtutils_lib Client
         │       │       │
         │       │       ├─→ Infrastructure error? → TestResult{INFRASTRUCTURE_ERROR}
         │       │       ├─→ Method error? → TestResult{METHOD_ERROR}
         │       │       └─→ Success? → Continue validation
         │       │
         │       ├─► Get appropriate Validator (ticket/user/asset)
         │       │
         │       ├─► Validator.Validate(result)
         │       │       │
         │       │       ├─→ Invalid? → TestResult{ASSERTION_FAILURE}
         │       │       └─→ Valid? → TestResult{PASS}
         │       │
         │       └─► Collect TestResult
         │
         ▼
┌─────────────────────────┐
│   TestReport            │
│  (accumulated results)  │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  JSON Reporter          │
│  (outputs report)       │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  report.json            │
│  (machine-readable)     │
└─────────────────────────┘
```

## Validator Implementation Pattern

All validators follow this pattern:

```
func (v *TicketValidator) Validate(result interface{}) (bool, string) {
    ticket, ok := result.(*Ticket)
    if !ok {
        return false, "Result is not a Ticket"
    }
    
    // Check required fields
    if ticket.ID == "" {
        return false, "Required field 'id' is empty"
    }
    
    if ticket.Status == "" {
        return false, "Required field 'status' is empty"
    }
    
    // Check for at least one additional field
    otherFieldCount := 0
    if ticket.Subject != "" { otherFieldCount++ }
    if ticket.Description != "" { otherFieldCount++ }
    if ticket.Queue != "" { otherFieldCount++ }
    if ticket.Owner != "" { otherFieldCount++ }
    // ... check other fields
    
    if otherFieldCount == 0 {
        return false, "No additional fields populated beyond id and status"
    }
    
    return true, ""
}
```

## Error States & Transitions

### Test Execution Error States

```
    ┌──────────────────┐
    │  Test Queued     │
    └────────┬─────────┘
             │
             ▼
    ┌──────────────────┐      Connection        ┌──────────────────┐
    │ Attempting       │──────Failed───────────►│ INFRASTRUCTURE   │
    │ Execution        │                        │ ERROR            │
    └────────┬─────────┘                        └──────────────────┘
             │
             ├─ Auth Failed ──────────────────────►  INFRASTRUCTURE ERROR
             │
             ├─ Timeout ──────────────────────────►  INFRASTRUCTURE ERROR
             │
             │
             ▼
    ┌──────────────────┐
    │ Method Executed  │
    └────────┬─────────┘
             │
             ├─ Panicked/Uncaught Error ──────────►  METHOD ERROR
             │
             ▼
    ┌──────────────────┐
    │ Result Returned  │
    └────────┬─────────┘
             │
             ▼
    ┌──────────────────┐
    │ Validate Result  │
    └────────┬─────────┘
             │
             ├─ Invalid ──────────────────────────►  ASSERTION FAILURE
             │
             └─ Valid ────────────────────────────►  PASS
```

## JSON Serialization

All TestResults and TestReport must be JSON-serializable:

```go
import "encoding/json"

// Example serialization
report := TestReport{
    Execution: ExecutionMetadata{
        Timestamp: "2026-02-05T14:23:45Z",
        TotalTests: 20,
        PassedTests: 18,
        FailedTests: 2,
        // ...
    },
    Results: []TestResult{
        {
            TestID: "ticket-get-1",
            MethodName: "TicketService.Get",
            Status: "PASS",
            // ...
        },
        // ... more results
    },
}

bytes, _ := json.MarshalIndent(report, "", "  ")
// Write bytes to file
```

## Type Safety Guarantees

- All validators operate on concrete types (not `interface{}`)
- Type assertions with checks prevent nil dereferences
- Result validation uses field presence checks before access
- Error messages include type information for debugging

## Future Extensions

### P2 Enhancements
- Per-test metadata (tags, categories)
- Custom validator plugins
- Retry logic for flaky tests
- Performance benchmarking

### P3 Enhancements
- Test data setup/teardown fixtures
- Parametrized test cases
- Custom assertion rules
