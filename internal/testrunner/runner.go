package testrunner

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"rtutils_lib"
	"rtutils_lib/internal/rtconfig"
	"rtutils_lib/internal/validators"
)

// TestRunner orchestrates the execution of test cases against a live RT instance
type TestRunner struct {
	Config    *rtconfig.RTConnection
	TestCases []TestCase
	Validator validators.ResultValidator
	Client    *rtutils_lib.Client
	Report    *TestReport
	StartTime time.Time
}

// NewTestRunner creates a new test runner with configuration and test cases
func NewTestRunner(config *rtconfig.RTConnection, testCases []TestCase) *TestRunner {
	return &TestRunner{
		Config:    config,
		TestCases: testCases,
		Report:    NewTestReport(config.URL),
	}
}

// Run executes all test cases sequentially and returns a test report
func (tr *TestRunner) Run() (*TestReport, error) {
	tr.StartTime = time.Now()

	// Validate configuration
	if err := tr.Config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create RT client using token-based authentication
	// For now, use username as a simple token placeholder
	// In production, this would use RT_TOKEN env var or actual token generation
	token := fmt.Sprintf("%s:%s", tr.Config.Username, tr.Config.Password)
	
	client := rtutils_lib.NewClient(tr.Config.URL, token)
	tr.Client = client

	// Execute each test case sequentially
	for _, testCase := range tr.TestCases {
		if !testCase.Enabled {
			continue
		}

		result := tr.executeTest(testCase)
		tr.Report.Results = append(tr.Report.Results, *result)

		// Update execution metadata
		tr.Report.Execution.TotalTests++
		if result.Status == "PASS" {
			tr.Report.Execution.PassedTests++
		} else {
			tr.Report.Execution.FailedTests++
			// Classify the error type
			if result.ErrorType != "" {
				switch result.ErrorType {
				case "InfrastructureError":
					tr.Report.Execution.InfrastructureErrors++
				case "AssertionFailure":
					tr.Report.Execution.AssertionFailures++
				case "MethodError":
					tr.Report.Execution.MethodErrors++
				}
			}
		}
	}

	// Calculate total execution time
	duration := time.Since(tr.StartTime)
	tr.Report.Execution.DurationSeconds = duration.Seconds()

	// Generate summary message
	tr.Report.Execution.SummaryMessage = fmt.Sprintf(
		"%d passed, %d failed (Infrastructure: %d, Assertion: %d, Method: %d)",
		tr.Report.Execution.PassedTests,
		tr.Report.Execution.FailedTests,
		tr.Report.Execution.InfrastructureErrors,
		tr.Report.Execution.AssertionFailures,
		tr.Report.Execution.MethodErrors,
	)

	return tr.Report, nil
}

// executeTest runs a single test case and returns the result
func (tr *TestRunner) executeTest(testCase TestCase) *TestResult {
	result := NewTestResult(testCase.ID, testCase.MethodName)
	startTime := time.Now()

	// Create validator for this test
	validator := validators.GetValidator(testCase.ValidatorType)

	// Simulate calling the RT client method
	// In real implementation, this would use reflection or method routing
	methodResult, err := tr.callMethod(testCase)

	if err != nil {
		// Classify error
		result.Status = "FAIL"
		result.ErrorType, result.ErrorMessage = tr.classifyError(err, testCase)
		duration := time.Since(startTime)
		result.DurationMs = duration.Milliseconds()
		return result
	}

	// Validate the result
	if err := validator.Validate(methodResult); err != nil {
		result.Status = "FAIL"
		result.ErrorType = "AssertionFailure"
		result.ErrorMessage = err.Error()
		duration := time.Since(startTime)
		result.DurationMs = duration.Milliseconds()
		return result
	}

	// Success
	result.Status = "PASS"
	result.ResultSummary = fmt.Sprintf("Successfully executed %s", testCase.MethodName)
	duration := time.Since(startTime)
	result.DurationMs = duration.Milliseconds()

	return result
}

// callMethod invokes the appropriate method on the RT client
func (tr *TestRunner) callMethod(testCase TestCase) (interface{}, error) {
	// This is a placeholder for method invocation
	// In reality, this would use reflection or a method routing map
	// For now, return nil to indicate method not implemented
	return nil, fmt.Errorf("method %s.%s not implemented", testCase.ServiceType, testCase.MethodName)
}

// classifyError determines the error type and returns appropriate classification
func (tr *TestRunner) classifyError(err error, testCase TestCase) (string, string) {
	errMsg := err.Error()

	// Infrastructure errors: connection, timeout, HTTP errors
	if strings.Contains(errMsg, "connection") ||
		strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "dial") ||
		strings.Contains(errMsg, "EOF") {
		return "InfrastructureError", errMsg
	}

	// Check for HTTP status code errors (5xx = infrastructure, 4xx = method)
	if strings.Contains(errMsg, "404") || strings.Contains(errMsg, "400") {
		return "MethodError", errMsg
	}
	if strings.Contains(errMsg, "500") || strings.Contains(errMsg, "502") || strings.Contains(errMsg, "503") {
		return "InfrastructureError", errMsg
	}

	// Default to method error if we can't determine
	return "MethodError", errMsg
}

// ReportJSON returns the test report as formatted JSON
func (tr *TestRunner) ReportJSON() ([]byte, error) {
	return json.MarshalIndent(tr.Report, "", "  ")
}

// WriteReport writes the test report to a writer
func (tr *TestRunner) WriteReport(w io.Writer) error {
	reportJSON, err := tr.ReportJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(reportJSON)
	return err
}

// GetSummary returns a text summary of the test execution
func (tr *TestRunner) GetSummary() string {
	return tr.Report.Execution.SummaryMessage
}

// ParseMethodResult converts interface{} result to map for validation
func ParseMethodResult(result interface{}) map[string]interface{} {
	if result == nil {
		return map[string]interface{}{
			"type": "nil",
		}
	}

	switch v := result.(type) {
	case map[string]interface{}:
		return v
	case string:
		// Try to unmarshal JSON string
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err == nil {
			return m
		}
	}

	// If we can't parse, create a minimal map with type info
	return map[string]interface{}{
		"type": reflect.TypeOf(result).String(),
	}
}
