package testrunner

import (
	"encoding/json"
	"fmt"
	"os"
)

// TestCaseLoader handles loading test cases from JSON files
type TestCaseLoader struct {
	FilePath string
}

// NewTestCaseLoader creates a new test case loader
func NewTestCaseLoader(filePath string) *TestCaseLoader {
	return &TestCaseLoader{
		FilePath: filePath,
	}
}

// LoadTestCases loads all test cases from the JSON file
func (tcl *TestCaseLoader) LoadTestCases() ([]TestCase, error) {
	// Read file
	data, err := os.ReadFile(tcl.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test cases file: %w", err)
	}

	// Unmarshal JSON
	var testCases []TestCase
	if err := json.Unmarshal(data, &testCases); err != nil {
		return nil, fmt.Errorf("failed to parse test cases JSON: %w", err)
	}

	// Validate loaded test cases
	if err := tcl.ValidateTestCases(testCases); err != nil {
		return nil, fmt.Errorf("test case validation failed: %w", err)
	}

	return testCases, nil
}

// ValidateTestCases checks that all required fields are present
func (tcl *TestCaseLoader) ValidateTestCases(testCases []TestCase) error {
	if len(testCases) == 0 {
		return fmt.Errorf("no test cases found in file")
	}

	for i, tc := range testCases {
		if err := tcl.ValidateTestCase(tc); err != nil {
			return fmt.Errorf("test case %d (%s) validation failed: %w", i, tc.ID, err)
		}
	}

	return nil
}

// ValidateTestCase checks that a single test case has all required fields
func (tcl *TestCaseLoader) ValidateTestCase(tc TestCase) error {
	if tc.ID == "" {
		return fmt.Errorf("missing required field: id")
	}
	if tc.MethodName == "" {
		return fmt.Errorf("missing required field: method")
	}
	if tc.ServiceType == "" {
		return fmt.Errorf("missing required field: service")
	}
	if tc.ValidatorType == "" {
		return fmt.Errorf("missing required field: validator_type")
	}
	if tc.Timeout == 0 {
		return fmt.Errorf("missing or invalid timeout value")
	}

	return nil
}

// GetTestCasesByService returns all test cases for a specific service
func GetTestCasesByService(testCases []TestCase, serviceType string) []TestCase {
	var filtered []TestCase
	for _, tc := range testCases {
		if tc.ServiceType == serviceType {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// GetTestCasesByValidator returns all test cases using a specific validator
func GetTestCasesByValidator(testCases []TestCase, validatorType string) []TestCase {
	var filtered []TestCase
	for _, tc := range testCases {
		if tc.ValidatorType == validatorType {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// GetEnabledTestCases returns only enabled test cases
func GetEnabledTestCases(testCases []TestCase) []TestCase {
	var enabled []TestCase
	for _, tc := range testCases {
		if tc.Enabled {
			enabled = append(enabled, tc)
		}
	}
	return enabled
}

// CountTestCasesByStatus returns counts of enabled vs disabled test cases
func CountTestCasesByStatus(testCases []TestCase) (enabled, disabled int) {
	for _, tc := range testCases {
		if tc.Enabled {
			enabled++
		} else {
			disabled++
		}
	}
	return
}

// GetServiceMethods returns a list of all unique services and their method counts
func GetServiceMethods(testCases []TestCase) map[string][]string {
	methods := make(map[string][]string)
	for _, tc := range testCases {
		methods[tc.ServiceType] = append(methods[tc.ServiceType], tc.MethodName)
	}
	return methods
}
