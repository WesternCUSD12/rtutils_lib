package testrunner

import (
	"encoding/json"
	"fmt"
	"os"

	"rtutils_lib/internal/validators"
)

// CustomTestCase extends TestCase with additional fields for custom tests
type CustomTestCase struct {
	TestCase
	Assertions  []string               `json:"assertions"`
	Tags        []string               `json:"tags,omitempty"`
	Preconditions []string             `json:"preconditions,omitempty"`
	CustomFields map[string]interface{} `json:"customFields,omitempty"`
}

// CustomTestRegistry maintains custom validators for user-defined test cases
type CustomTestRegistry struct {
	validators map[string]validators.ResultValidator
}

// NewCustomTestRegistry creates a new custom test registry
func NewCustomTestRegistry() *CustomTestRegistry {
	return &CustomTestRegistry{
		validators: make(map[string]validators.ResultValidator),
	}
}

// RegisterValidator registers a custom validator for a specific test type
func (r *CustomTestRegistry) RegisterValidator(name string, validator validators.ResultValidator) error {
	if name == "" {
		return fmt.Errorf("validator name cannot be empty")
	}
	if validator == nil {
		return fmt.Errorf("validator cannot be nil")
	}
	
	r.validators[name] = validator
	return nil
}

// GetValidator retrieves a registered validator by name
func (r *CustomTestRegistry) GetValidator(name string) (validators.ResultValidator, bool) {
	validator, ok := r.validators[name]
	return validator, ok
}

// HasValidator checks if a validator is registered
func (r *CustomTestRegistry) HasValidator(name string) bool {
	_, ok := r.validators[name]
	return ok
}

// ListValidators returns all registered validator names
func (r *CustomTestRegistry) ListValidators() []string {
	names := make([]string, 0, len(r.validators))
	for name := range r.validators {
		names = append(names, name)
	}
	return names
}

// LoadCustomTests loads custom test cases from a JSON file
func LoadCustomTests(filePath string) ([]CustomTestCase, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom tests file: %w", err)
	}

	var testCases []CustomTestCase
	if err := json.Unmarshal(data, &testCases); err != nil {
		return nil, fmt.Errorf("failed to parse custom tests JSON: %w", err)
	}

	return testCases, nil
}

// ValidateCustomTestCase validates a custom test case has required fields
func ValidateCustomTestCase(tc CustomTestCase) error {
	if tc.ID == "" {
		return fmt.Errorf("id is required")
	}
	if tc.MethodName == "" {
		return fmt.Errorf("methodName is required")
	}
	if tc.ServiceType == "" {
		return fmt.Errorf("serviceType is required")
	}
	if tc.ValidatorType == "" {
		return fmt.Errorf("validatorType is required")
	}
	if len(tc.Assertions) == 0 {
		return fmt.Errorf("assertions array cannot be empty")
	}

	return nil
}

// ValidateCustomTestCases validates all custom test cases
func ValidateCustomTestCases(testCases []CustomTestCase) []error {
	var errors []error
	
	seen := make(map[string]bool)
	for i, tc := range testCases {
		// Check for duplicate IDs
		if seen[tc.ID] {
			errors = append(errors, fmt.Errorf("duplicate id at index %d: %s", i, tc.ID))
		}
		seen[tc.ID] = true

		// Validate individual test case
		if err := ValidateCustomTestCase(tc); err != nil {
			errors = append(errors, fmt.Errorf("invalid test case at index %d: %w", i, err))
		}
	}

	return errors
}

// GetCustomTestsByService filters custom tests by service type
func GetCustomTestsByService(tests []CustomTestCase, service string) []CustomTestCase {
	var filtered []CustomTestCase
	for _, tc := range tests {
		if tc.ServiceType == service {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// GetCustomTestsByTag filters custom tests by tag
func GetCustomTestsByTag(tests []CustomTestCase, tag string) []CustomTestCase {
	var filtered []CustomTestCase
	for _, tc := range tests {
		for _, t := range tc.Tags {
			if t == tag {
				filtered = append(filtered, tc)
				break
			}
		}
	}
	return filtered
}

// GetCustomTestsByValidator filters custom tests by validator type
func GetCustomTestsByValidator(tests []CustomTestCase, validator string) []CustomTestCase {
	var filtered []CustomTestCase
	for _, tc := range tests {
		if tc.ValidatorType == validator {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// GetEnabledCustomTests returns only enabled custom tests
func GetEnabledCustomTests(tests []CustomTestCase) []CustomTestCase {
	var enabled []CustomTestCase
	for _, tc := range tests {
		if tc.Enabled {
			enabled = append(enabled, tc)
		}
	}
	return enabled
}

// MergeTestCases combines standard and custom test cases
// Custom tests take precedence if there are duplicates by ID
func MergeTestCases(standard []TestCase, custom []CustomTestCase) []TestCase {
	merged := make([]TestCase, 0, len(standard)+len(custom))
	seen := make(map[string]bool)

	// Add custom tests first (they have priority)
	for _, ct := range custom {
		merged = append(merged, ct.TestCase)
		seen[ct.ID] = true
	}

	// Add standard tests that don't conflict
	for _, st := range standard {
		if !seen[st.ID] {
			merged = append(merged, st)
			seen[st.ID] = true
		}
	}

	return merged
}

// GetCustomTestsByDescription filters by partial description match
func GetCustomTestsByDescription(tests []CustomTestCase, keyword string) []CustomTestCase {
	var filtered []CustomTestCase
	for _, tc := range tests {
		if contains(tc.TestCase.Description, keyword) {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// CountCustomTestsByValidator counts custom tests grouped by validator type
func CountCustomTestsByValidator(tests []CustomTestCase) map[string]int {
	counts := make(map[string]int)
	for _, tc := range tests {
		counts[tc.ValidatorType]++
	}
	return counts
}

// CountCustomTestsByService counts custom tests grouped by service
func CountCustomTestsByService(tests []CustomTestCase) map[string]int {
	counts := make(map[string]int)
	for _, tc := range tests {
		counts[tc.ServiceType]++
	}
	return counts
}

// contains is a helper function for string containment check (case-sensitive)
func contains(s, substring string) bool {
	for i := 0; i < len(s)-len(substring)+1; i++ {
		if s[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

