package testrunner

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestCaseLoader(t *testing.T) {
	tcl := NewTestCaseLoader("test_cases.json")
	assert.Equal(t, "test_cases.json", tcl.FilePath)
}

func TestLoadTestCasesFileNotFound(t *testing.T) {
	tcl := NewTestCaseLoader("nonexistent.json")
	_, err := tcl.LoadTestCases()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read test cases file")
}

func TestLoadTestCasesValidJSON(t *testing.T) {
	// Create temporary test file
	tmpFile := createTempTestCasesFile(t, []TestCase{
		{
			ID:            "TEST001",
			MethodName:    "Get",
			ServiceType:   "TicketService",
			ValidatorType: "TicketValidator",
			Timeout:       5,
			Enabled:       true,
		},
		{
			ID:            "TEST002",
			MethodName:    "List",
			ServiceType:   "UserService",
			ValidatorType: "UserValidator",
			Timeout:       10,
			Enabled:       true,
		},
	})
	defer os.Remove(tmpFile)

	tcl := NewTestCaseLoader(tmpFile)
	testCases, err := tcl.LoadTestCases()

	require.NoError(t, err)
	assert.Equal(t, 2, len(testCases))
	assert.Equal(t, "TEST001", testCases[0].ID)
	assert.Equal(t, "TEST002", testCases[1].ID)
}

func TestLoadTestCasesInvalidJSON(t *testing.T) {
	// Create temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "invalid_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("{ invalid json }")
	tmpFile.Close()

	tcl := NewTestCaseLoader(tmpFile.Name())
	_, err = tcl.LoadTestCases()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse test cases JSON")
}

func TestValidateTestCasesEmpty(t *testing.T) {
	tcl := NewTestCaseLoader("dummy.json")
	err := tcl.ValidateTestCases([]TestCase{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no test cases found")
}

func TestValidateTestCasesMissingFields(t *testing.T) {
	tcl := NewTestCaseLoader("dummy.json")

	// Missing ID
	err := tcl.ValidateTestCase(TestCase{
		MethodName:    "Get",
		ServiceType:   "TicketService",
		ValidatorType: "TicketValidator",
		Timeout:       5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "id")

	// Missing method
	err = tcl.ValidateTestCase(TestCase{
		ID:            "TEST001",
		ServiceType:   "TicketService",
		ValidatorType: "TicketValidator",
		Timeout:       5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "method")

	// Missing service type
	err = tcl.ValidateTestCase(TestCase{
		ID:            "TEST001",
		MethodName:    "Get",
		ValidatorType: "TicketValidator",
		Timeout:       5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service")

	// Missing validator type
	err = tcl.ValidateTestCase(TestCase{
		ID:          "TEST001",
		MethodName:  "Get",
		ServiceType: "TicketService",
		Timeout:     5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validator_type")

	// Missing timeout
	err = tcl.ValidateTestCase(TestCase{
		ID:            "TEST001",
		MethodName:    "Get",
		ServiceType:   "TicketService",
		ValidatorType: "TicketValidator",
		Timeout:       0,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

func TestValidateTestCaseValid(t *testing.T) {
	tcl := NewTestCaseLoader("dummy.json")

	err := tcl.ValidateTestCase(TestCase{
		ID:            "TEST001",
		MethodName:    "Get",
		ServiceType:   "TicketService",
		ValidatorType: "TicketValidator",
		Timeout:       5,
		Enabled:       true,
	})

	assert.NoError(t, err)
}

func TestGetTestCasesByService(t *testing.T) {
	testCases := []TestCase{
		{ID: "T001", ServiceType: "TicketService"},
		{ID: "T002", ServiceType: "UserService"},
		{ID: "T003", ServiceType: "TicketService"},
		{ID: "T004", ServiceType: "AssetService"},
	}

	tickets := GetTestCasesByService(testCases, "TicketService")
	assert.Equal(t, 2, len(tickets))
	assert.Equal(t, "T001", tickets[0].ID)
	assert.Equal(t, "T003", tickets[1].ID)

	users := GetTestCasesByService(testCases, "UserService")
	assert.Equal(t, 1, len(users))
	assert.Equal(t, "T002", users[0].ID)

	assets := GetTestCasesByService(testCases, "AssetService")
	assert.Equal(t, 1, len(assets))
	assert.Equal(t, "T004", assets[0].ID)

	// Non-existent service
	none := GetTestCasesByService(testCases, "NonExistent")
	assert.Equal(t, 0, len(none))
}

func TestGetTestCasesByValidator(t *testing.T) {
	testCases := []TestCase{
		{ID: "T001", ValidatorType: "TicketValidator"},
		{ID: "T002", ValidatorType: "UserValidator"},
		{ID: "T003", ValidatorType: "TicketValidator"},
		{ID: "T004", ValidatorType: "GenericValidator"},
	}

	tickets := GetTestCasesByValidator(testCases, "TicketValidator")
	assert.Equal(t, 2, len(tickets))

	users := GetTestCasesByValidator(testCases, "UserValidator")
	assert.Equal(t, 1, len(users))

	generic := GetTestCasesByValidator(testCases, "GenericValidator")
	assert.Equal(t, 1, len(generic))
}

func TestGetEnabledTestCases(t *testing.T) {
	testCases := []TestCase{
		{ID: "T001", Enabled: true},
		{ID: "T002", Enabled: false},
		{ID: "T003", Enabled: true},
		{ID: "T004", Enabled: false},
		{ID: "T005", Enabled: true},
	}

	enabled := GetEnabledTestCases(testCases)
	assert.Equal(t, 3, len(enabled))
	assert.Equal(t, "T001", enabled[0].ID)
	assert.Equal(t, "T003", enabled[1].ID)
	assert.Equal(t, "T005", enabled[2].ID)
}

func TestCountTestCasesByStatus(t *testing.T) {
	testCases := []TestCase{
		{ID: "T001", Enabled: true},
		{ID: "T002", Enabled: false},
		{ID: "T003", Enabled: true},
		{ID: "T004", Enabled: true},
	}

	enabled, disabled := CountTestCasesByStatus(testCases)
	assert.Equal(t, 3, enabled)
	assert.Equal(t, 1, disabled)
}

func TestCountTestCasesByStatusEmpty(t *testing.T) {
	testCases := []TestCase{}

	enabled, disabled := CountTestCasesByStatus(testCases)
	assert.Equal(t, 0, enabled)
	assert.Equal(t, 0, disabled)
}

func TestGetServiceMethods(t *testing.T) {
	testCases := []TestCase{
		{ID: "T001", ServiceType: "TicketService", MethodName: "Get"},
		{ID: "T002", ServiceType: "TicketService", MethodName: "Search"},
		{ID: "T003", ServiceType: "UserService", MethodName: "List"},
		{ID: "T004", ServiceType: "AssetService", MethodName: "Get"},
		{ID: "T005", ServiceType: "TicketService", MethodName: "GetByURL"},
	}

	serviceMethods := GetServiceMethods(testCases)

	assert.Equal(t, 3, len(serviceMethods))
	assert.Equal(t, 3, len(serviceMethods["TicketService"]))
	assert.Equal(t, 1, len(serviceMethods["UserService"]))
	assert.Equal(t, 1, len(serviceMethods["AssetService"]))

	// Check method names
	assert.Contains(t, serviceMethods["TicketService"], "Get")
	assert.Contains(t, serviceMethods["TicketService"], "Search")
	assert.Contains(t, serviceMethods["TicketService"], "GetByURL")
}

// Helper function to create temporary test cases file
func createTempTestCasesFile(t *testing.T, testCases []TestCase) string {
	data, err := json.Marshal(testCases)
	require.NoError(t, err)

	tmpFile, err := os.CreateTemp("", "test_cases_*.json")
	require.NoError(t, err)
	defer tmpFile.Close()

	_, err = tmpFile.Write(data)
	require.NoError(t, err)

	return tmpFile.Name()
}
