package testrunner

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestCaseJSON(t *testing.T) {
	testCaseJSON := `{
		"id": "TICKET001",
		"method": "Get",
		"service": "TicketService",
		"description": "Get a ticket by ID",
		"inputs": {"id": "123"},
		"expected_result_type": "Ticket",
		"validator_type": "TicketValidator",
		"timeout": 5,
		"enabled": true
	}`

	var tc TestCase
	err := json.Unmarshal([]byte(testCaseJSON), &tc)
	require.NoError(t, err)

	assert.Equal(t, "TICKET001", tc.ID)
	assert.Equal(t, "Get", tc.MethodName)
	assert.Equal(t, "TicketService", tc.ServiceType)
	assert.Equal(t, "TicketValidator", tc.ValidatorType)
	assert.True(t, tc.Enabled)
	assert.Equal(t, 5, tc.Timeout)
}

func TestNewTestResult(t *testing.T) {
	result := NewTestResult("TICKET001", "Get")

	assert.Equal(t, "TICKET001", result.TestID)
	assert.Equal(t, "Get", result.MethodName)
	assert.NotEmpty(t, result.StartTime)
	assert.Equal(t, "", result.Status)
	assert.Equal(t, int64(0), result.DurationMs)
}

func TestTestResultJSON(t *testing.T) {
	result := &TestResult{
		TestID:        "TICKET001",
		MethodName:    "Get",
		Status:        "PASS",
		StartTime:     "2024-01-15T10:30:00Z",
		DurationMs:    125,
		ResultSummary: "Successfully retrieved ticket with id 123",
	}

	resultJSON, err := json.Marshal(result)
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(resultJSON, &data)
	require.NoError(t, err)

	assert.Equal(t, "TICKET001", data["test_id"])
	assert.Equal(t, "Get", data["method"])
	assert.Equal(t, "PASS", data["status"])
	assert.Equal(t, float64(125), data["duration_ms"])
}

func TestTestResultWithError(t *testing.T) {
	result := NewTestResult("ERROR001", "Search")
	result.Status = "FAIL"
	result.ErrorType = "AssertionFailure"
	result.ErrorMessage = "Expected 'active' but got 'stalled'"
	result.DurationMs = 450

	assert.Equal(t, "FAIL", result.Status)
	assert.Equal(t, "AssertionFailure", result.ErrorType)
	assert.Contains(t, result.ErrorMessage, "Expected")
}

func TestNewTestReport(t *testing.T) {
	report := NewTestReport("https://rt.example.com")

	assert.Equal(t, "https://rt.example.com", report.Execution.RTInstanceURL)
	assert.NotEmpty(t, report.Execution.Timestamp)
	assert.Equal(t, 0, report.Execution.TotalTests)
	assert.Equal(t, 0, report.Execution.PassedTests)
	assert.Empty(t, report.Results)
}

func TestTestReportAggregation(t *testing.T) {
	report := NewTestReport("https://rt.example.com")

	result1 := NewTestResult("TEST001", "Get")
	result1.Status = "PASS"
	result1.DurationMs = 100

	result2 := NewTestResult("TEST002", "Search")
	result2.Status = "FAIL"
	result2.ErrorType = "AssertionFailure"
	result2.DurationMs = 200

	result3 := NewTestResult("TEST003", "List")
	result3.Status = "PASS"
	result3.DurationMs = 150

	report.Results = append(report.Results, *result1, *result2, *result3)

	report.Execution.TotalTests = 3
	report.Execution.PassedTests = 2
	report.Execution.FailedTests = 1

	assert.Equal(t, 3, report.Execution.TotalTests)
	assert.Equal(t, 2, report.Execution.PassedTests)
	assert.Equal(t, 1, report.Execution.FailedTests)
	assert.Equal(t, 3, len(report.Results))
}

func TestExecutionMetadata(t *testing.T) {
	metadata := ExecutionMetadata{
		Timestamp:            time.Now().UTC().Format(time.RFC3339),
		DurationSeconds:      2.45,
		RTInstanceURL:        "https://rt.example.com",
		TotalTests:           5,
		PassedTests:          3,
		FailedTests:          2,
		InfrastructureErrors: 1,
		AssertionFailures:    1,
		MethodErrors:         0,
		SummaryMessage:       "3 passed, 2 failed",
	}

	assert.NotEmpty(t, metadata.Timestamp)
	assert.Equal(t, 2.45, metadata.DurationSeconds)
	assert.Equal(t, 5, metadata.TotalTests)
	assert.Equal(t, 3, metadata.PassedTests)
	assert.Equal(t, 2, metadata.FailedTests)
	assert.Equal(t, 1, metadata.InfrastructureErrors)
	assert.Equal(t, 1, metadata.AssertionFailures)
	assert.Equal(t, 0, metadata.MethodErrors)
}

func TestTestReportJSON(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 1
	report.Execution.PassedTests = 1

	result := NewTestResult("TEST001", "Get")
	result.Status = "PASS"
	result.DurationMs = 100
	report.Results = append(report.Results, *result)

	reportJSON, err := json.Marshal(report)
	require.NoError(t, err)

	var data map[string]interface{}
	err = json.Unmarshal(reportJSON, &data)
	require.NoError(t, err)

	assert.NotNil(t, data["execution"])
	assert.NotNil(t, data["results"])

	execution := data["execution"].(map[string]interface{})
	assert.Equal(t, "https://rt.example.com", execution["rt_instance"])
	assert.Equal(t, float64(1), execution["total_tests"])
	assert.Equal(t, float64(1), execution["passed"])
}

func TestTestResultStartTimeFormat(t *testing.T) {
	result := NewTestResult("TEST001", "Get")

	_, err := time.Parse(time.RFC3339, result.StartTime)
	require.NoError(t, err)
}

func TestExecutionMetadataErrorCounting(t *testing.T) {
	metadata := ExecutionMetadata{
		Timestamp:            time.Now().UTC().Format(time.RFC3339),
		TotalTests:           10,
		PassedTests:          7,
		FailedTests:          3,
		InfrastructureErrors: 1,
		AssertionFailures:    2,
		MethodErrors:         0,
	}

	totalErrors := metadata.InfrastructureErrors + metadata.AssertionFailures + metadata.MethodErrors
	assert.Equal(t, 3, totalErrors)
	assert.Equal(t, metadata.FailedTests, totalErrors)
}
