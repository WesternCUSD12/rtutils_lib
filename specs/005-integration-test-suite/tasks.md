# Tasks: RT Integration Test Suite

**Branch**: `005-integration-test-suite` | **Date**: February 5, 2026
**Spec**: [spec.md](spec.md) | **Plan**: [plan.md](plan.md)

## Overview

This document breaks the implementation plan into executable tasks organized by user story and phase. Each task is independent-testable and can be assigned to developers.

**Total Tasks**: 38  
**Estimated Duration**: 3-4 weeks for experienced Go developer  
**MVP MVP Scope (US1 + US2)**: ~2 weeks

---

## Dependencies & Execution Order

```
Phase 1: Setup
    ↓
Phase 2: Foundational Infrastructure
    ↓ (depends on Phase 2)
├─ Phase 3: US1 & US2 (Run & Validate) ◄── MVP Deliverable
    ├─ Phase 4: US3 (Reporting) 
    ├─ Phase 5: US4 (CI/CD)
    └─ Phase 6: US5 (Custom Tests)
    └─ Phase 7: Polish & Integration Testing
```

### Parallel Opportunity
- US3, US4, US5 can be developed in parallel after Phase 2 completes
- Each has independent test coverage and no cross-dependencies

---

## Phase 1: Setup & Project Initialization

### Project Structure

- [X] T001 Create `cmd/rtutils-test-suite/` directory structure with `main.go`
- [X] T002 Create `internal/testrunner/` package directory
- [X] T003 Create `internal/validators/` package directory
- [X] T004 Create `internal/rtconfig/` package directory
- [X] T005 Create `examples/` directory with sample files
- [X] T006 Create `.gitignore` entry for test binaries and reports
- [X] T007 Update `go.mod` with module declaration for internal packages
- [X] T008 Verify all packages build: `go build ./...`

---

## Phase 2: Foundational Infrastructure (Blocking Prerequisites)

### RTConfig Package

Handles RT instance connection configuration.

- [X] T009 [P] Implement `internal/rtconfig/config.go` with RTConnection struct
  - RTConnection with URL, Username, Password, Timeout fields
  - LoadFromEnv() function to read RT_URL, RT_USERNAME, RT_PASSWORD
  - Validate() method to check credentials present
  - GetClient() to create configured rtutils_lib Client

- [X] T010 [P] Implement `internal/rtconfig/config_test.go` 
  - Test LoadFromEnv with valid env vars
  - Test Validate() catches missing fields
  - Test error handling for invalid URLs

### Data Models

Define core entities for test execution.

- [X] T011 [P] Implement `internal/testrunner/models.go` with:
  - TestCase struct (ID, MethodName, ServiceType, InputParams, ValidatorType)
  - TestResult struct (TestID, Status, ErrorType, DurationMs, ErrorMessage, etc.)
  - TestReport struct (Execution ExecutionMetadata, Results []TestResult)
  - ExecutionMetadata struct (Timestamp, TotalTests, PassedTests, FailedTests, etc.)

### Validator Interface & Base Implementation

- [X] T012 [P] Implement `internal/validators/validators.go` with:
  - ResultValidator interface (Validate(result interface{}) (bool, string))
  - GetValidator(entityType string) ResultValidator function
  - Type checking utilities

- [X] T013 [P] Implement `internal/validators/validators_test.go` with unit tests for validator selection

### Git Integration

- [X] T014 [P] Commit Phase 2 deliverables: "feat: Phase 2 setup - project structure and foundational models"

---

## Phase 3: User Story 1 & 2 - Run Tests & Validate Results (MVP Core)

### MVP: US1 - Execute Test Suite Against Live RT

Test the ability to run all methods against a live RT instance.

#### Validators Implementation (Blocking for all tests)

- [ ] T015 [P] [US1] Implement `internal/validators/ticket_validator.go`
  - TicketValidator struct
  - Validate() method: check id + status + other field
  - RequiredFields() method
  - See contracts/validators.md for specifications

- [ ] T016 [P] [US1] Implement `internal/validators/ticket_validator_test.go`
  - Test valid ticket passes
  - Test missing id fails
  - Test missing status fails
  - Test missing additional field fails
  - Test empty strings treated as missing

- [ ] T017 [P] [US1] Implement `internal/validators/user_validator.go`
  - UserValidator struct
  - Validate() method: check id + name
  - RequiredFields() method

- [ ] T018 [P] [US1] Implement `internal/validators/user_validator_test.go`
  - Test valid user passes
  - Test missing id fails
  - Test missing name fails

- [ ] T019 [P] [US1] Implement `internal/validators/asset_validator.go`
  - AssetValidator struct
  - Validate() method: check id + name
  - RequiredFields() method

- [ ] T020 [P] [US1] Implement `internal/validators/asset_validator_test.go`
  - Test valid asset passes
  - Test missing id fails
  - Test missing name fails

- [ ] T021 [P] [US1] Implement `internal/validators/generic_validator.go`
  - GenericValidator struct (fallback for unknown types)
  - Validate() method: check id + other field
  - RequiredFields() method

#### Test Runner Implementation

- [ ] T022 [US1] Implement `internal/testrunner/runner.go` with:
  - TestRunner struct (Client, Config, TestCases)
  - Run() method for sequential execution
  - executeTest() for single test execution
  - classifyError() for error classification logic
  - See plan.md execution flow diagram

- [ ] T023 [P] [US1] Implement `internal/testrunner/runner_test.go`
  - Mock Client for testing
  - Test sequential execution order
  - Test error classification (Infrastructure/Assertion/Method)
  - Test result collection

#### Test Case Loader

- [ ] T024 [P] [US1] Implement `internal/testrunner/test_case_loader.go`
  - LoadTestCases() from contracts/test-cases.json
  - Populate TestCase slice with all 18 read-only methods
  - Validation of loaded test cases

- [ ] T025 [P] [US1] Implement `internal/testrunner/test_case_loader_test.go`
  - Test loading JSON
  - Test correct number of test cases (18)
  - Test test case fields populated correctly

#### US2 Implementation (Tests run validation assertions)

- [ ] T026 [US2] Refine executeTest() in runner.go to:
  - Call appropriate Validator for result
  - Check validation results
  - Populate TestResult with assertion details
  - Handle assertion failures vs method errors

- [ ] T027 [P] [US2] Add integration tests for validator execution:
  - Run ticket tests and verify validation
  - Run user tests and verify validation
  - Run asset tests and verify validation

#### CLI Entry Point (MVP)

- [ ] T028 [US1] Implement `cmd/rtutils-test-suite/main.go` with:
  - Parse command-line flags (--config, --report, --debug)
  - Read config (env vars or JSON)
  - Create RTConnection and load Client
  - Initialize TestRunner
  - Call Run()
  - Return appropriate exit code (0 for pass, 1 for fail)

- [ ] T029 [P] [US1] Implement `cmd/rtutils-test-suite/main_test.go`
  - Test flag parsing
  - Test config loading
  - Test exit code logic

#### Build & Verify MVP

- [ ] T030 [US1] Build and test MVP:
  - `go build -o rtutils-test-suite cmd/rtutils-test-suite/main.go`
  - Run against live RT with export RT_URL, RT_USERNAME, RT_PASSWORD
  - Verify all 18 tests execute
  - Verify pass/fail reported correctly

- [ ] T031 [US2] Verify validation works end-to-end:
  - Confirm useful data requirement enforced
  - Confirm invalid results rejected
  - Confirm empty results fail with clear message

- [ ] T032 [US1] Commit MVP: "feat: US1+US2 MVP - test execution and result validation"

---

## Phase 4: User Story 3 - Generate Detailed Test Report

Generate machine-readable JSON reports.

- [ ] T033 [P] [US3] Implement `internal/testrunner/reporter.go` with:
  - GenerateReport(results []TestResult) TestReport function
  - JSON marshaling with proper formatting
  - Timestamp generation (ISO 8601)
  - Statistics calculation (passed, failed, infrastructure errors)
  - Summary message generation

- [ ] T034 [P] [US3] Implement `internal/testrunner/reporter_test.go`
  - Test report generation with mixed results
  - Test JSON validity
  - Test statistics accuracy
  - Test summary message

- [ ] T035 [P] [US3] Update `cmd/rtutils-test-suite/main.go`:
  - Call reporter to generate JSON
  - Write report to file (specify via --report flag or default: `results.json`)
  - Pretty-print report to stdout

- [ ] T036 [P] [US3] Generate example report output:
  - Run full test suite
  - Save to `examples/test-report-sample.json`
  - Update quickstart.md with example

- [ ] T037 [US3] Commit US3: "feat: US3 - JSON test report generation"

---

## Phase 5: User Story 4 - Integrate Tests into CI/CD Pipeline

Enable automated CI/CD integration.

- [ ] T038 [P] [US4] Update `cmd/rtutils-test-suite/main.go`:
  - Implement proper exit codes (0 for all pass, 1 for failures)
  - Env var support for CI/CD (RT_URL, RT_USERNAME, RT_PASSWORD)
  - Log output to both stdout and file
  - Handle missing env vars gracefully

- [ ] T039 [P] [US4] Create GitHub Actions workflow example:
  - File: `.github/workflows/integration-tests.yml`
  - Test on push and pull_request
  - Checkout, build, run tests
  - Store report as artifact
  - Comment on PR with results

- [ ] T040 [P] [US4] Create GitLab CI pipeline example:
  - File: `.gitlab-ci.yml` stage for integration tests
  - Run tests, save report
  - Set up CI/CD secrets

- [ ] T041 [US4] Commit US4: "feat: US4 - CI/CD pipeline integration"

---

## Phase 6: User Story 5 - Support Custom Test Cases

Extensible test framework for custom queries.

- [ ] T042 [P] [US5] Implement custom test case support:
  - `internal/testrunner/custom_tests.go` with:
    - CustomTestCase struct extending TestCase
    - LoadCustomTests(configFile string) []CustomTestCase
    - RegisterValidator(name string, validator ResultValidator) function

- [ ] T043 [P] [US5] Update runner.go to include custom tests:
  - Initialize custom tests if provided via flag
  - Merge custom tests with standard tests
  - Execute custom tests in same sequential flow

- [ ] T044 [P] [US5] Create custom test config example:
  - File: `examples/custom-tests.json`
  - Example: custom query for specific ticket state
  - Document custom test structure

- [ ] T045 [P] [US5] Update `cmd/rtutils-test-suite/main.go`:
  - Add --custom-tests flag for custom test file
  - Load and integrate custom tests

- [ ] T046 [US5] Commit US5: "feat: US5 - custom test case support"

---

## Phase 7: Polish & Integration Testing

Final validation and release preparation.

### Documentation

- [ ] T047 Create `cmd/rtutils-test-suite/README.md`:
  - Build instructions
  - Usage examples
  - Troubleshooting guide

### Integration Testing

- [ ] T048 [P] Run full integration test suite:
  - Test all 18 methods against live RT
  - Verify validation for each entity type
  - Generate and verify JSON report
  - Test with invalid credentials
  - Test with unreachable RT instance
  - Test timeout handling

- [ ] T049 [P] Performance testing:
  - Measure execution time (target: <5 minutes)
  - Identify slow methods
  - Verify <1ms per validator call

- [ ] T050 [P] Error scenario testing:
  - RT offline → infrastructure error reported
  - Auth failure → infrastructure error reported
  - Missing required field → assertion failure with details
  - Method exception → method error with stack trace

### CI/CD Integration Verification

- [ ] T051 [P] Test workflows in CI system:
  - Push to branch, verify tests run
  - Verify exit codes work correctly
  - Verify report artifacts created

### Final Cleanup

- [ ] T052 [P] Code review & refactoring:
  - Run `gofmt` on all files
  - Run `go vet` for linting
  - Remove debug logging
  - Add package-level comments (doc strings)

- [ ] T053 [P] Update go.mod:
  - Verify no unnecessary dependencies
  - Document dependency justification
  - Run `go mod tidy`

- [ ] T054 Commit Phase 7: "feat: final polish and integration testing"

- [ ] T055 Create release notes:
  - Feature summary
  - Usage instructions
  - Known limitations
  - Future enhancements

- [ ] T056 Merge to main: "Release: RT Integration Test Suite MVP"

---

## Testing Strategy

### Unit Tests
- Task T010, T013, T016, T018, T020, T022, T025, T034, T049

### Integration Tests  
- Task T023, T027, T030, T031, T035, T048, T050, T051

### Performance Tests
- Task T049

### Manual Testing
- Task T030, T031, T048, T050, T051

---

## Parallel Execution Opportunities

**Critical Path** (Sequential, blocking):
- T001-T008 (Setup)
- T009-T014 (Foundational)
- T015-T032 (MVP US1+US2)

**After T032, Deploy in Parallel**:
- Branch A: T033-T037 (US3 Reporting)
- Branch B: T038-T041 (US4 CI/CD)
- Branch C: T042-T046 (US5 Custom Tests)

**Recommended Team Assignment**:
- Developer 1: US1+US2 (core MVP) - 2 weeks
- Developer 2: Setup + US3 (reporting) - 1.5 weeks
- Developer 3: US4+US5 (integration+extensions) - 1.5 weeks

---

## Success Metrics

- [ ] All 38 tasks complete
- [ ] All unit tests passing (`go test ./...`)
- [ ] All integration tests passing against live RT
- [ ] Test suite executes in <5 minutes
- [ ] 95%+ tests pass on healthy RT instance
- [ ] JSON report validates against schema
- [ ] CI/CD integration working
- [ ] Zero code style warnings (gofmt, go vet)

---

## Task Template for Implementation

Each developer should follow this for each task:

1. **Create branch** (if not already on 005-integration-test-suite):
   ```bash
   git checkout -b 005-integration-test-suite
   ```

2. **Create/update files** as specified in task description

3. **Write tests first** (TDD per constitution):
   - Create `_test.go` file with failing tests
   - Run `go test` to verify tests fail
   - Implement code to pass tests

4. **Verify code quality**:
   ```bash
   go fmt ./internal/...
   go vet ./internal/...
   go test ./... -v
   ```

5. **Commit with clear message**:
   ```bash
   git commit -m "feat: Tnnn [description]"
   ```

6. **Push and wait for review**

---

## Rollback Plan

If a task or phase fails:
1. Identify blocking issue
2. Create hotfix branch from latest stable
3. Fix issue in isolation
4. Test fix thoroughly
5. Merge back to main branch

---

## Next Steps

1. **Assign tasks** to team members
2. **Start Phase 1** (Setup) immediately
3. **Complete Phase 2** (Foundational) before US1/US2 development
4. **Deploy MVP** (US1+US2) first for early validation
5. **Parallelize** US3, US4, US5 development
6. **Integrate** and test before final release
