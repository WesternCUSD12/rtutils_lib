package testrunner

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReportFormatter(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	formatter := NewReportFormatter(report)

	assert.NotNil(t, formatter)
	assert.Equal(t, report, formatter.Report)
}

func TestFormatJSONPretty(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 1
	report.Execution.PassedTests = 1

	formatter := NewReportFormatter(report)
	jsonBytes, err := formatter.FormatJSON(true)

	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON
	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, &data)
	require.NoError(t, err)

	// Pretty formatting includes newlines
	assert.True(t, strings.Contains(string(jsonBytes), "\n"))
}

func TestFormatJSONCompact(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 1
	report.Execution.PassedTests = 1

	formatter := NewReportFormatter(report)
	jsonBytes, err := formatter.FormatJSON(false)

	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON
	var data map[string]interface{}
	err = json.Unmarshal(jsonBytes, &data)
	require.NoError(t, err)
}

func TestFormatText(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 5
	report.Execution.PassedTests = 3
	report.Execution.FailedTests = 2
	report.Execution.InfrastructureErrors = 1
	report.Execution.AssertionFailures = 1
	report.Execution.MethodErrors = 0
	report.Execution.DurationSeconds = 1.5
	report.Execution.SummaryMessage = "3 passed, 2 failed"

	formatter := NewReportFormatter(report)
	text := formatter.FormatText()

	assert.NotEmpty(t, text)
	assert.Contains(t, text, "RT Integration Test Suite - Execution Report")
	assert.Contains(t, text, "Total Tests:")
	assert.Contains(t, text, "Passed:")
	assert.Contains(t, text, "Failed:")
	assert.Contains(t, text, "3 passed, 2 failed")
	assert.Contains(t, text, "SOME TESTS FAILED")
}

func TestFormatTextAllPassed(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 3
	report.Execution.PassedTests = 3
	report.Execution.FailedTests = 0

	formatter := NewReportFormatter(report)
	text := formatter.FormatText()

	assert.Contains(t, text, "ALL TESTS PASSED")
	assert.NotContains(t, text, "Error Breakdown")
}

func TestFormatTextNilReport(t *testing.T) {
	formatter := NewReportFormatter(nil)
	text := formatter.FormatText()

	assert.Equal(t, "No report data available", text)
}

func TestWriteJSON(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 1
	report.Execution.PassedTests = 1

	formatter := NewReportFormatter(report)
	var buf bytes.Buffer

	err := formatter.WriteJSON(&buf, true)
	require.NoError(t, err)

	// Verify JSON written
	jsonData := buf.String()
	assert.NotEmpty(t, jsonData)
	assert.Contains(t, jsonData, "execution")
	assert.Contains(t, jsonData, "results")
}

func TestWriteText(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 1
	report.Execution.PassedTests = 1

	formatter := NewReportFormatter(report)
	var buf bytes.Buffer

	err := formatter.WriteText(&buf)
	require.NoError(t, err)

	// Verify text written
	text := buf.String()
	assert.NotEmpty(t, text)
	assert.Contains(t, text, "Execution Report")
}

func TestGenerateReport(t *testing.T) {
	startTime := time.Now()

	results := []TestResult{
		{
			TestID:     "TEST001",
			MethodName: "Get",
			Status:     "PASS",
			DurationMs: 100,
			StartTime:  time.Now().UTC().Format(time.RFC3339),
		},
		{
			TestID:     "TEST002",
			MethodName: "Search",
			Status:     "FAIL",
			ErrorType:  "AssertionFailure",
			DurationMs: 50,
			StartTime:  time.Now().UTC().Format(time.RFC3339),
		},
	}

	report := GenerateReport("https://rt.example.com", results, startTime)

	assert.NotNil(t, report)
	assert.Equal(t, 2, report.Execution.TotalTests)
	assert.Equal(t, 1, report.Execution.PassedTests)
	assert.Equal(t, 1, report.Execution.FailedTests)
	assert.Equal(t, 1, report.Execution.AssertionFailures)
	assert.Equal(t, 0, report.Execution.InfrastructureErrors)
	assert.Contains(t, report.Execution.SummaryMessage, "1 passed, 1 failed")
}

func TestGenerateReportWithMultipleErrors(t *testing.T) {
	startTime := time.Now()

	results := []TestResult{
		{TestID: "T001", Status: "PASS"},
		{TestID: "T002", Status: "PASS"},
		{TestID: "T003", Status: "FAIL", ErrorType: "InfrastructureError"},
		{TestID: "T004", Status: "FAIL", ErrorType: "AssertionFailure"},
		{TestID: "T005", Status: "FAIL", ErrorType: "MethodError"},
	}

	report := GenerateReport("https://rt.example.com", results, startTime)

	assert.Equal(t, 5, report.Execution.TotalTests)
	assert.Equal(t, 2, report.Execution.PassedTests)
	assert.Equal(t, 3, report.Execution.FailedTests)
	assert.Equal(t, 1, report.Execution.InfrastructureErrors)
	assert.Equal(t, 1, report.Execution.AssertionFailures)
	assert.Equal(t, 1, report.Execution.MethodErrors)
}

func TestFilterResultsByStatus(t *testing.T) {
	results := []TestResult{
		{TestID: "T001", Status: "PASS"},
		{TestID: "T002", Status: "FAIL"},
		{TestID: "T003", Status: "PASS"},
		{TestID: "T004", Status: "FAIL"},
	}

	passed := FilterResultsByStatus(results, "PASS")
	assert.Equal(t, 2, len(passed))
	assert.Equal(t, "T001", passed[0].TestID)
	assert.Equal(t, "T003", passed[1].TestID)

	failed := FilterResultsByStatus(results, "FAIL")
	assert.Equal(t, 2, len(failed))
	assert.Equal(t, "T002", failed[0].TestID)
	assert.Equal(t, "T004", failed[1].TestID)
}

func TestGetFailedResults(t *testing.T) {
	results := []TestResult{
		{TestID: "T001", Status: "PASS"},
		{TestID: "T002", Status: "FAIL"},
		{TestID: "T003", Status: "PASS"},
		{TestID: "T004", Status: "FAIL"},
	}

	failed := GetFailedResults(results)
	assert.Equal(t, 2, len(failed))
	for _, r := range failed {
		assert.Equal(t, "FAIL", r.Status)
	}
}

func TestGetPassedResults(t *testing.T) {
	results := []TestResult{
		{TestID: "T001", Status: "PASS"},
		{TestID: "T002", Status: "FAIL"},
		{TestID: "T003", Status: "PASS"},
	}

	passed := GetPassedResults(results)
	assert.Equal(t, 2, len(passed))
	for _, r := range passed {
		assert.Equal(t, "PASS", r.Status)
	}
}

func TestGetResultsWithErrorType(t *testing.T) {
	results := []TestResult{
		{TestID: "T001", Status: "FAIL", ErrorType: "InfrastructureError"},
		{TestID: "T002", Status: "FAIL", ErrorType: "AssertionFailure"},
		{TestID: "T003", Status: "FAIL", ErrorType: "InfrastructureError"},
		{TestID: "T004", Status: "PASS"},
	}

	infra := GetResultsWithErrorType(results, "InfrastructureError")
	assert.Equal(t, 2, len(infra))
	assert.Equal(t, "T001", infra[0].TestID)
	assert.Equal(t, "T003", infra[1].TestID)

	assertion := GetResultsWithErrorType(results, "AssertionFailure")
	assert.Equal(t, 1, len(assertion))
	assert.Equal(t, "T002", assertion[0].TestID)
}

func TestFormatTextErrorBreakdown(t *testing.T) {
	report := NewTestReport("https://rt.example.com")
	report.Execution.TotalTests = 3
	report.Execution.PassedTests = 0
	report.Execution.FailedTests = 3
	report.Execution.InfrastructureErrors = 1
	report.Execution.AssertionFailures = 1
	report.Execution.MethodErrors = 1

	formatter := NewReportFormatter(report)
	text := formatter.FormatText()

	assert.Contains(t, text, "Error Breakdown:")
	assert.Contains(t, text, "Infrastructure: 1")
	assert.Contains(t, text, "Assertions:     1")
	assert.Contains(t, text, "Method Errors:  1")
}

func TestReportDurationCalculation(t *testing.T) {
	startTime := time.Now().Add(-2 * time.Second)

	results := []TestResult{
		{TestID: "T001", Status: "PASS"},
	}

	report := GenerateReport("https://rt.example.com", results, startTime)

	// Duration should be approximately 2 seconds
	assert.True(t, report.Execution.DurationSeconds >= 1.9)
	assert.True(t, report.Execution.DurationSeconds <= 2.5)
}
