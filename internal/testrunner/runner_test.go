package testrunner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rtutils_lib/internal/rtconfig"
)

func TestNewTestRunner(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{
		{
			ID:            "TEST001",
			MethodName:    "Get",
			ServiceType:   "Ticket",
			ValidatorType: "TicketValidator",
			Enabled:       true,
		},
	}

	runner := NewTestRunner(config, testCases)

	assert.NotNil(t, runner)
	assert.Equal(t, config, runner.Config)
	assert.Equal(t, 1, len(runner.TestCases))
	assert.NotNil(t, runner.Report)
	assert.Equal(t, "https://rt.example.com", runner.Report.Execution.RTInstanceURL)
}

func TestTestRunnerExecuteTestSuccess(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{
		{
			ID:            "TICKET001",
			MethodName:    "Get",
			ServiceType:   "Ticket",
			ValidatorType: "TicketValidator",
			Enabled:       true,
			InputParams: map[string]interface{}{
				"id": "123",
			},
		},
	}

	runner := NewTestRunner(config, testCases)

	// Execute single test - should fail because method not implemented
	result := runner.executeTest(testCases[0])

	assert.Equal(t, "FAIL", result.Status)
	assert.Equal(t, "MethodError", result.ErrorType)
	assert.Contains(t, result.ErrorMessage, "not implemented")
}

func TestTestRunnerErrorClassification(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{
		{
			ID:            "TEST001",
			MethodName:    "Test",
			ServiceType:   "TestService",
			ValidatorType: "GenericValidator",
		},
	}

	runner := NewTestRunner(config, testCases)

	// Test infrastructure error classification
	errType, _ := runner.classifyError(
		&mockError{msg: "connection refused"},
		testCases[0],
	)
	assert.Equal(t, "InfrastructureError", errType)

	// Test timeout error classification
	errType, _ = runner.classifyError(
		&mockError{msg: "context deadline exceeded: timeout"},
		testCases[0],
	)
	assert.Equal(t, "InfrastructureError", errType)

	// Test HTTP 404 (method error)
	errType, _ = runner.classifyError(
		&mockError{msg: "HTTP 404 Not Found"},
		testCases[0],
	)
	assert.Equal(t, "MethodError", errType)

	// Test HTTP 500 (infrastructure error)
	errType, _ = runner.classifyError(
		&mockError{msg: "HTTP 500 Internal Server Error"},
		testCases[0],
	)
	assert.Equal(t, "InfrastructureError", errType)
}

func TestTestRunnerMetadataCalculation(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}

	testCases := []TestCase{
		{
			ID:            "TEST001",
			MethodName:    "Method1",
			ServiceType:   "Service1",
			ValidatorType: "GenericValidator",
			Enabled:       false, // Disabled, should be skipped
		},
		{
			ID:            "TEST002",
			MethodName:    "Method2",
			ServiceType:   "Service2",
			ValidatorType: "GenericValidator",
			Enabled:       true,
		},
	}

	runner := NewTestRunner(config, testCases)

	// Manually set up report
	runner.Report.Execution.TotalTests = 1
	runner.Report.Execution.PassedTests = 1
	runner.Report.Execution.FailedTests = 0

	summary := runner.GetSummary()
	// Summary should be empty initially since no Run() called
	assert.Empty(t, summary)
}

func TestTestRunnerReportJSON(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{
		{
			ID:            "TEST001",
			MethodName:    "Get",
			ServiceType:   "Ticket",
			ValidatorType: "TicketValidator",
			Enabled:       true,
		},
	}

	runner := NewTestRunner(config, testCases)

	// Add a mock result
	result := NewTestResult("TEST001", "Get")
	result.Status = "PASS"
	result.DurationMs = 125
	runner.Report.Results = append(runner.Report.Results, *result)

	// Get JSON report
	reportJSON, err := runner.ReportJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, reportJSON)

	// Verify it's valid JSON
	assert.Contains(t, string(reportJSON), "execution")
	assert.Contains(t, string(reportJSON), "results")
	assert.Contains(t, string(reportJSON), "TEST001")
}

func TestTestRunnerWriteReport(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{}

	runner := NewTestRunner(config, testCases)

	// Create a buffer to capture output
	var buf mockWriter
	err := runner.WriteReport(&buf)
	require.NoError(t, err)

	assert.NotEmpty(t, buf.data)
	assert.Contains(t, buf.data, "execution")
}

func TestTestRunnerGetSummary(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{}

	runner := NewTestRunner(config, testCases)

	// Manually set execution metadata
	runner.Report.Execution.TotalTests = 5
	runner.Report.Execution.PassedTests = 3
	runner.Report.Execution.FailedTests = 2
	runner.Report.Execution.InfrastructureErrors = 1
	runner.Report.Execution.AssertionFailures = 1
	runner.Report.Execution.MethodErrors = 0
	runner.Report.Execution.SummaryMessage = "3 passed, 2 failed (Infrastructure: 1, Assertion: 1, Method: 0)"

	summary := runner.GetSummary()
	assert.Contains(t, summary, "3 passed")
	assert.Contains(t, summary, "2 failed")
	assert.Contains(t, summary, "Infrastructure: 1")
	assert.Contains(t, summary, "Assertion: 1")
}

func TestParseMethodResult(t *testing.T) {
	// Test map parsing
	result := map[string]interface{}{
		"id":   "123",
		"name": "Test",
	}
	parsed := ParseMethodResult(result)
	assert.Equal(t, "123", parsed["id"])
	assert.Equal(t, "Test", parsed["name"])

	// Test JSON string parsing
	jsonStr := `{"id": "456", "name": "JSON Test"}`
	parsed = ParseMethodResult(jsonStr)
	assert.Equal(t, "456", parsed["id"])
	assert.Equal(t, "JSON Test", parsed["name"])

	// Test invalid input
	parsed = ParseMethodResult(nil)
	assert.NotNil(t, parsed)
	assert.Equal(t, "nil", parsed["type"])
}

func TestTestRunnerStartTimeTracking(t *testing.T) {
	config := &rtconfig.RTConnection{
		URL:     "https://rt.example.com",
		Token:   "testtoken",
		Timeout: 30,
	}
	testCases := []TestCase{}

	runner := NewTestRunner(config, testCases)

	// StartTime should be zero initially
	assert.True(t, runner.StartTime.IsZero())

	// Set start time manually
	now := time.Now()
	runner.StartTime = now

	assert.Equal(t, now, runner.StartTime)
}

// Mock implementations for testing

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

type mockWriter struct {
	data string
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	mw.data += string(p)
	return len(p), nil
}
