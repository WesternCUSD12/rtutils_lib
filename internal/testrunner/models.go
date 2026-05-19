package testrunner

import "time"

// TestCase represents a single test that validates one library method
type TestCase struct {
	ID            string                 `json:"id"`
	MethodName    string                 `json:"method"`
	ServiceType   string                 `json:"service"`
	Operation     string                 `json:"operation"`
	Description   string                 `json:"description"`
	InputParams   map[string]interface{} `json:"inputs"`
	ExpectedType  string                 `json:"expected_result_type"`
	ValidatorType string                 `json:"validator_type"`
	Timeout       int                    `json:"timeout"`
	Enabled       bool                   `json:"enabled"`
}

// TestResult represents the outcome of running a single test
type TestResult struct {
	TestID        string      `json:"test_id"`
	MethodName    string      `json:"method"`
	Status        string      `json:"status"`
	ErrorType     string      `json:"error_type,omitempty"`
	StartTime     string      `json:"start_time"`
	DurationMs    int64       `json:"duration_ms"`
	ResultSummary string      `json:"result_summary,omitempty"`
	ErrorMessage  string      `json:"error_message,omitempty"`
	ExpectedValue interface{} `json:"expected,omitempty"`
	ActualValue   interface{} `json:"actual,omitempty"`
}

// ExecutionMetadata contains overall execution statistics
type ExecutionMetadata struct {
	Timestamp            string  `json:"timestamp"`
	DurationSeconds      float64 `json:"duration_seconds"`
	RTInstanceURL        string  `json:"rt_instance"`
	TotalTests           int     `json:"total_tests"`
	PassedTests          int     `json:"passed"`
	FailedTests          int     `json:"failed"`
	InfrastructureErrors int     `json:"infrastructure_errors"`
	AssertionFailures    int     `json:"assertion_failures"`
	MethodErrors         int     `json:"method_errors"`
	SummaryMessage       string  `json:"summary"`
}

// TestReport aggregates multiple TestResults with summary statistics
type TestReport struct {
	Execution ExecutionMetadata `json:"execution"`
	Results   []TestResult      `json:"results"`
}

// NewTestResult creates a new test result with common fields populated
func NewTestResult(testID, methodName string) *TestResult {
	return &TestResult{
		TestID:     testID,
		MethodName: methodName,
		StartTime:  time.Now().UTC().Format(time.RFC3339),
	}
}

// NewTestReport creates a new test report with execution metadata
func NewTestReport(rtURL string) *TestReport {
	return &TestReport{
		Execution: ExecutionMetadata{
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			RTInstanceURL: rtURL,
		},
		Results: []TestResult{},
	}
}
