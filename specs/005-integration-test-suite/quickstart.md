# Quickstart: RT Integration Test Suite

**Date**: February 5, 2026  
**Phase**: Phase 1 (Design)  
**Status**: Complete

## Overview

The RT Integration Test Suite is a Go CLI tool that validates all read-only query methods in the rtutils library against a live Request Tracker instance.

**Status**: Phase 2 (Implementation) - Coming Soon  
**Expected Completion**: [To be determined during planning]

## Quick Start (Once Implemented)

### Prerequisites

- Go 1.22+ installed
- Access to a live Request Tracker instance with test data
- Network connectivity to RT instance

### Setup

#### 1. Clone & Build

```bash
# Clone the repository
git clone https://github.com/your-org/rtutils_lib.git
cd rtutils_lib

# Build the test suite
go build -o rtutils-test-suite cmd/rtutils-test-suite/main.go
```

#### 2. Configure RT Connection

Set environment variables:

```bash
export RT_URL="https://rt.example.com"
export RT_USERNAME="testuser"
export RT_PASSWORD="testpass123"
```

Or create a configuration file:

```json
{
  "rt_url": "https://rt.example.com",
  "rt_username": "testuser",
  "rt_password": "testpass123",
  "timeout_seconds": 30,
  "insecure_tls": false
}
```

#### 3. Run Tests

```bash
# Using environment variables
./rtutils-test-suite

# Using config file
./rtutils-test-suite --config rt-config.json

# With output report
./rtutils-test-suite --report results.json
```

### Expected Output

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
      "result_summary": "Retrieved ticket #123 with status 'new', queue 'General'"
    },
    {
      "test_id": "user-search-email-2",
      "method": "UserService.SearchByEmailExact",
      "status": "FAIL",
      "error_type": "ASSERTION_FAILURE",
      "error_message": "Expected 'name' field to be non-empty",
      "expected": { "name": "non-empty" },
      "actual": { "name": "" }
    }
  ],
  "summary": "18 of 20 tests passed. 2 assertion failures detected in user search methods."
}
```

### Exit Codes

- `0` — All tests passed
- `1` — One or more tests failed
- `2` — Infrastructure error (cannot connect to RT)

### What Gets Tested

#### TicketService Methods

| Method | Type | Status |
|--------|------|--------|
| Get(id) | Read | ✓ Included |
| GetByURL(url) | Read | ✓ Included |
| Search(query, page, perPage) | Read | ✓ Included |
| SearchBySubject(query) | Read | ✓ Included |
| Create | Write | ✗ Not included (Phase 1 MVP) |
| Update | Write | ✗ Not included (Phase 1 MVP) |

#### UserService Methods

| Method | Type | Status |
|--------|------|--------|
| Get(id) | Read | ✓ Included |
| Search(query) | Read | ✓ Included |
| SearchByUsernameExact | Read | ✓ Included |
| SearchByUsernamePartial | Read | ✓ Included |
| SearchByEmailExact | Read | ✓ Included |
| SearchByEmailPartial | Read | ✓ Included |
| SearchByNameExact | Read | ✓ Included |
| SearchByNamePartial | Read | ✓ Included |
| Create | Write | ✗ Not included (Phase 1 MVP) |
| Update | Write | ✗ Not included (Phase 1 MVP) |

#### AssetService Methods

| Method | Type | Status |
|--------|------|--------|
| Get(id) | Read | ✓ Included |
| Search(query) | Read | ✓ Included |
| SearchByNameExact | Read | ✓ Included |
| SearchByNamePartial | Read | ✓ Included |
| SearchByCustomFieldExact | Read | ✓ Included |
| SearchByCustomFieldPartial | Read | ✓ Included |
| Create | Write | ✗ Not included (Phase 1 MVP) |
| Update | Write | ✗ Not included (Phase 1 MVP) |

**Total**: ~20 read-only methods tested

## How Results Are Validated

### Meaningful Data Validation

Each test verifies that results contain useful data, not just syntactically valid responses.

#### Ticket Results Must Have:
- `id` field (non-empty)
- `status` field (non-empty)
- At least one additional field (`subject`, `description`, `queue`, `owner`, etc.)

#### User Results Must Have:
- `id` field (non-empty)
- `name` field (non-empty)

#### Asset Results Must Have:
- `id` field (non-empty)
- `name` field (non-empty)

### Error Classification

Failures are classified to help identify root causes:

| Error Type | Meaning | Example |
|-----------|---------|---------|
| INFRASTRUCTURE_ERROR | Cannot connect to RT instance | Connection timeout, auth failure |
| ASSERTION_FAILURE | Method worked but returned incomplete/invalid data | Missing required field in result |
| METHOD_ERROR | Method execution failed | Unhandled exception, panic |

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      
      - name: Build Test Suite
        run: go build -o rtutils-test-suite cmd/rtutils-test-suite/main.go
      
      - name: Run Integration Tests
        env:
          RT_URL: ${{ secrets.RT_TEST_INSTANCE_URL }}
          RT_USERNAME: ${{ secrets.RT_TEST_USERNAME }}
          RT_PASSWORD: ${{ secrets.RT_TEST_PASSWORD }}
        run: ./rtutils-test-suite --report report.json
      
      - name: Parse Results
        if: always()
        run: |
          # Check exit code or parse JSON report
          passed=$(jq .execution.passed report.json)
          total=$(jq .execution.total_tests report.json)
          echo "Passed: $passed/$total"
```

### GitLab CI Example

```yaml
integration-tests:
  stage: test
  image: golang:1.22
  script:
    - go build -o rtutils-test-suite cmd/rtutils-test-suite/main.go
    - ./rtutils-test-suite --report report.json
  artifacts:
    reports:
      dotenv: report.json
    paths:
      - report.json
  only:
    - merge_requests
    - main
```

## Troubleshooting

### Connection Failed

**Problem**: "Cannot connect to RT instance"

**Solution**:
- Verify `RT_URL` is correct and accessible
- Check network connectivity: `curl -I https://rt.example.com`
- Verify TLS certificates (use `--insecure-tls` for testing if needed)

### Authentication Failed

**Problem**: "Invalid credentials" or "Unauthorized"

**Solution**:
- Verify `RT_USERNAME` and `RT_PASSWORD` are correct
- Confirm user has API access in RT
- Check RT logs for auth errors

### Tests Timeout

**Problem**: "Request timeout after 30 seconds"

**Solution**:
- Increase timeout: `RT_TIMEOUT=60` (in seconds)
- Check RT instance performance
- Verify network latency

### Assertion Failures

**Problem**: "Expected field 'status' to be non-empty"

**Solution**:
- Check that test RT instance has sufficient data
- Verify tickets/users/assets have populated fields
- Review RT instance configuration

## Example Test Run

```bash
$ export RT_URL=https://rt.example.com
$ export RT_USERNAME=admin
$ export RT_PASSWORD=secret
$ ./rtutils-test-suite --report test-report.json

[INFO] Connecting to RT instance: https://rt.example.com
[INFO] Running 20 tests...
[INFO] ✓ TicketService.Get (145ms)
[INFO] ✓ TicketService.GetByURL (132ms)
[INFO] ✓ TicketService.Search (487ms)
[INFO] ✓ TicketService.SearchBySubject (412ms)
[INFO] ✓ UserService.Get (98ms)
[INFO] ✓ UserService.Search (156ms)
[INFO] ✓ UserService.SearchByUsernameExact (187ms)
[INFO] ✓ UserService.SearchByUsernamePartial (203ms)
[INFO] ✓ UserService.SearchByEmailExact (175ms)
[INFO] ✓ UserService.SearchByEmailPartial (192ms)
[INFO] ✓ UserService.SearchByNameExact (168ms)
[INFO] ✓ UserService.SearchByNamePartial (201ms)
[INFO] ✓ AssetService.Get (127ms)
[INFO] ✓ AssetService.Search (356ms)
[INFO] ✓ AssetService.SearchByNameExact (215ms)
[INFO] ✓ AssetService.SearchByNamePartial (198ms)
[INFO] ✓ AssetService.SearchByCustomFieldExact (289ms)
[INFO] ✓ AssetService.SearchByCustomFieldPartial (267ms)
[INFO] ✗ UserService.SearchByEmailExact (234ms) - ASSERTION_FAILURE
[INFO] ✗ AssetService.SearchByNameExact (156ms) - ASSERTION_FAILURE
[INFO] Completed in 4.2 seconds
[INFO] Results: 18 PASS, 2 FAIL
[INFO] Report saved to: test-report.json

$ echo $?
1  # Exit with failure code due to 2 failed tests
```

## Next Steps

### For Developers

1. Understanding the validation rules [→ See data-model.md](data-model.md)
2. Running against your RT instance [→ See setup instructions above]
3. Debugging assertion failures [→ Examine report.json for details]
4. Extending with custom tests [→ Phase 2 feature]

### For CI/CD Teams

1. Integrate into your pipeline [→ See CI/CD examples above]
2. Set up secrets for RT credentials
3. Parse JSON report for dashboard integration
4. Set up alerts for test failures

### For Contributors

1. Review the specification [→ spec.md](spec.md)
2. Check the data model [→ data-model.md](data-model.md)
3. Examine the implementation plan [→ plan.md](plan.md)
4. Follow the TDD workflow from the constitution

## Additional Information

- **Constitution**: [rtutils_lib Constitution](.specify/memory/constitution.md)
- **Specification**: [Full feature spec](spec.md)
- **Implementation Plan**: [Detailed plan](plan.md)
- **rtutils_lib Repository**: [GitHub](https://github.com/your-org/rtutils_lib)
