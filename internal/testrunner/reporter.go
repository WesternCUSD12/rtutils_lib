package testrunner

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ReportFormatter handles JSON report generation and formatting
type ReportFormatter struct {
	Report *TestReport
}

// NewReportFormatter creates a new report formatter
func NewReportFormatter(report *TestReport) *ReportFormatter {
	return &ReportFormatter{
		Report: report,
	}
}

// FormatJSON returns the report as formatted JSON
func (rf *ReportFormatter) FormatJSON(pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(rf.Report, "", "  ")
	}
	return json.Marshal(rf.Report)
}

// FormatText returns a human-readable text summary
func (rf *ReportFormatter) FormatText() string {
	if rf.Report == nil {
		return "No report data available"
	}

	output := ""
	output += "╔═══════════════════════════════════════════════════════════╗\n"
	output += "║        RT Integration Test Suite - Execution Report        ║\n"
	output += "╚═══════════════════════════════════════════════════════════╝\n\n"

	// Execution metadata
	exec := rf.Report.Execution
	output += fmt.Sprintf("Timestamp:          %s\n", exec.Timestamp)
	output += fmt.Sprintf("RT Instance:        %s\n", exec.RTInstanceURL)
	output += fmt.Sprintf("Duration:           %.2f seconds\n\n", exec.DurationSeconds)

	// Test results summary
	output += "Test Results:\n"
	output += fmt.Sprintf("  Total Tests:       %d\n", exec.TotalTests)
	output += fmt.Sprintf("  ✓ Passed:          %d\n", exec.PassedTests)
	output += fmt.Sprintf("  ✗ Failed:          %d\n", exec.FailedTests)

	if exec.FailedTests > 0 {
		output += "\n  Error Breakdown:\n"
		if exec.InfrastructureErrors > 0 {
			output += fmt.Sprintf("    • Infrastructure: %d\n", exec.InfrastructureErrors)
		}
		if exec.AssertionFailures > 0 {
			output += fmt.Sprintf("    • Assertions:     %d\n", exec.AssertionFailures)
		}
		if exec.MethodErrors > 0 {
			output += fmt.Sprintf("    • Method Errors:  %d\n", exec.MethodErrors)
		}
	}

	// Summary message
	output += fmt.Sprintf("\nSummary:            %s\n", exec.SummaryMessage)

	// Status indicator
	output += "\n"
	if exec.FailedTests == 0 {
		output += "Status: ✓ ALL TESTS PASSED\n"
	} else {
		output += "Status: ✗ SOME TESTS FAILED\n"
	}

	return output
}

// WriteJSON writes the JSON report to a writer
func (rf *ReportFormatter) WriteJSON(w io.Writer, pretty bool) error {
	formatted, err := rf.FormatJSON(pretty)
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}
	_, err = w.Write(formatted)
	return err
}

// WriteText writes the text report to a writer
func (rf *ReportFormatter) WriteText(w io.Writer) error {
	text := rf.FormatText()
	_, err := w.Write([]byte(text))
	return err
}

// GenerateReport creates a TestReport from test results
// This is a high-level function that could be used with alternative report sources
func GenerateReport(rtURL string, results []TestResult, startTime time.Time) *TestReport {
	report := NewTestReport(rtURL)
	report.Results = results
	report.Execution.TotalTests = len(results)

	// Calculate pass/fail
	passed := 0
	infrastructureErrors := 0
	assertionFailures := 0
	methodErrors := 0

	for _, result := range results {
		if result.Status == "PASS" {
			passed++
		} else {
			switch result.ErrorType {
			case "InfrastructureError":
				infrastructureErrors++
			case "AssertionFailure":
				assertionFailures++
			case "MethodError":
				methodErrors++
			}
		}
	}

	report.Execution.PassedTests = passed
	report.Execution.FailedTests = len(results) - passed
	report.Execution.InfrastructureErrors = infrastructureErrors
	report.Execution.AssertionFailures = assertionFailures
	report.Execution.MethodErrors = methodErrors

	// Calculate duration
	duration := time.Since(startTime)
	report.Execution.DurationSeconds = duration.Seconds()

	// Generate summary
	report.Execution.SummaryMessage = fmt.Sprintf(
		"%d passed, %d failed (Infrastructure: %d, Assertion: %d, Method: %d)",
		passed,
		len(results)-passed,
		infrastructureErrors,
		assertionFailures,
		methodErrors,
	)

	return report
}

// FilterResultsByStatus returns results matching the given status
func FilterResultsByStatus(results []TestResult, status string) []TestResult {
	var filtered []TestResult
	for _, r := range results {
		if r.Status == status {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// FilterResultsByService returns results for a specific service
func FilterResultsByService(results []TestResult, service string) []TestResult {
	var filtered []TestResult
	for _, r := range results {
		if r.MethodName != "" {
			// This is a simple filter based on method name
			// In a more complete implementation, we could store service type in TestResult
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// GetFailedResults returns all failed test results
func GetFailedResults(results []TestResult) []TestResult {
	return FilterResultsByStatus(results, "FAIL")
}

// GetPassedResults returns all passed test results
func GetPassedResults(results []TestResult) []TestResult {
	return FilterResultsByStatus(results, "PASS")
}

// GetResultsWithErrorType returns results matching a specific error type
func GetResultsWithErrorType(results []TestResult, errorType string) []TestResult {
	var filtered []TestResult
	for _, r := range results {
		if r.ErrorType == errorType {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
