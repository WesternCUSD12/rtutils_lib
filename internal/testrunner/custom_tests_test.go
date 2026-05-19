package testrunner

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rtutils_lib/internal/validators"
)

func TestNewCustomTestRegistry(t *testing.T) {
	registry := NewCustomTestRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.validators)
	assert.Equal(t, 0, len(registry.validators))
}

func TestRegisterValidator(t *testing.T) {
	registry := NewCustomTestRegistry()
	validator := &validators.GenericValidator{ValidatorName: "test"}

	err := registry.RegisterValidator("custom1", validator)
	assert.NoError(t, err)
	assert.True(t, registry.HasValidator("custom1"))
}

func TestRegisterValidatorEmptyName(t *testing.T) {
	registry := NewCustomTestRegistry()
	validator := &validators.GenericValidator{ValidatorName: "test"}

	err := registry.RegisterValidator("", validator)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestRegisterValidatorNil(t *testing.T) {
	registry := NewCustomTestRegistry()

	err := registry.RegisterValidator("invalid", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestGetValidator(t *testing.T) {
	registry := NewCustomTestRegistry()
	validator := &validators.GenericValidator{ValidatorName: "test"}

	registry.RegisterValidator("test", validator)
	retrieved, ok := registry.GetValidator("test")

	assert.True(t, ok)
	assert.Equal(t, validator, retrieved)
}

func TestGetValidatorNotFound(t *testing.T) {
	registry := NewCustomTestRegistry()

	validator, ok := registry.GetValidator("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, validator)
}

func TestHasValidator(t *testing.T) {
	registry := NewCustomTestRegistry()
	validator := &validators.GenericValidator{ValidatorName: "test"}

	registry.RegisterValidator("existing", validator)

	assert.True(t, registry.HasValidator("existing"))
	assert.False(t, registry.HasValidator("missing"))
}

func TestListValidators(t *testing.T) {
	registry := NewCustomTestRegistry()

	registry.RegisterValidator("val1", &validators.GenericValidator{ValidatorName: "val1"})
	registry.RegisterValidator("val2", &validators.GenericValidator{ValidatorName: "val2"})
	registry.RegisterValidator("val3", &validators.GenericValidator{ValidatorName: "val3"})

	validators := registry.ListValidators()
	assert.Equal(t, 3, len(validators))
	assert.Contains(t, validators, "val1")
	assert.Contains(t, validators, "val2")
	assert.Contains(t, validators, "val3")
}

func TestLoadCustomTests(t *testing.T) {
	// Create temporary test file
	testData := `[
    {
      "id": "CUSTOM-001",
      "method": "CustomSearch",
      "service": "TicketService",
      "validator_type": "TicketValidator",
      "enabled": true,
      "description": "Search with custom query",
      "assertions": ["Check result count > 0", "Verify ticket status"]
    }
  ]`

	tmpFile, err := os.CreateTemp("", "custom_tests_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(testData)
	require.NoError(t, err)
	tmpFile.Close()

	testCases, err := LoadCustomTests(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(testCases))
	assert.Equal(t, "CUSTOM-001", testCases[0].ID)
	assert.Equal(t, "CustomSearch", testCases[0].MethodName)
}

func TestLoadCustomTestsFileNotFound(t *testing.T) {
	testCases, err := LoadCustomTests("/nonexistent/file.json")

	assert.Error(t, err)
	assert.Nil(t, testCases)
	assert.Contains(t, err.Error(), "failed to read")
}

func TestLoadCustomTestsInvalidJSON(t *testing.T) {
	testData := `invalid json`

	tmpFile, err := os.CreateTemp("", "invalid_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(testData)
	tmpFile.Close()

	testCases, err := LoadCustomTests(tmpFile.Name())
	assert.Error(t, err)
	assert.Nil(t, testCases)
	assert.Contains(t, err.Error(), "parse")
}

func TestValidateCustomTestCaseValid(t *testing.T) {
	tc := CustomTestCase{
		TestCase: TestCase{
			ID:            "TEST-001",
			MethodName:    "Get",
			ServiceType:   "TicketService",
			ValidatorType: "TicketValidator",
			Enabled:       true,
		},
		Assertions: []string{"Check ID present"},
	}

	err := ValidateCustomTestCase(tc)
	assert.NoError(t, err)
}

func TestValidateCustomTestCaseMissingID(t *testing.T) {
	tc := CustomTestCase{
		TestCase: TestCase{
			MethodName:    "Get",
			ServiceType:   "TicketService",
			ValidatorType: "TicketValidator",
		},
		Assertions: []string{"Check ID present"},
	}

	err := ValidateCustomTestCase(tc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestValidateCustomTestCaseMissingAssertions(t *testing.T) {
	tc := CustomTestCase{
		TestCase: TestCase{
			ID:            "TEST-001",
			MethodName:    "Get",
			ServiceType:   "TicketService",
			ValidatorType: "TicketValidator",
		},
		Assertions: []string{},
	}

	err := ValidateCustomTestCase(tc)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "assertions")
}

func TestValidateCustomTestCases(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase: TestCase{
				ID:            "T001",
				MethodName:    "Get",
				ServiceType:   "TicketService",
				ValidatorType: "TicketValidator",
			},
			Assertions: []string{"Valid"},
		},
		{
			TestCase: TestCase{
				ID:            "T002",
				MethodName:    "Search",
				ServiceType:   "TicketService",
				ValidatorType: "TicketValidator",
			},
			Assertions: []string{"Valid"},
		},
	}

	errors := ValidateCustomTestCases(cases)
	assert.Empty(t, errors)
}

func TestValidateCustomTestCasesDuplicate(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase: TestCase{
				ID:            "DUP",
				MethodName:    "Get",
				ServiceType:   "TicketService",
				ValidatorType: "TicketValidator",
			},
			Assertions: []string{"Valid"},
		},
		{
			TestCase: TestCase{
				ID:            "DUP",
				MethodName:    "Search",
				ServiceType:   "TicketService",
				ValidatorType: "TicketValidator",
			},
			Assertions: []string{"Valid"},
		},
	}

	errors := ValidateCustomTestCases(cases)
	assert.Greater(t, len(errors), 0)
	assert.Contains(t, errors[0].Error(), "duplicate")
}

func TestGetCustomTestsByService(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase:   TestCase{ID: "T1", ServiceType: "TicketService"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T2", ServiceType: "UserService"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T3", ServiceType: "TicketService"},
			Assertions: []string{"Valid"},
		},
	}

	filtered := GetCustomTestsByService(cases, "TicketService")
	assert.Equal(t, 2, len(filtered))
	assert.Equal(t, "T1", filtered[0].ID)
	assert.Equal(t, "T3", filtered[1].ID)
}

func TestGetCustomTestsByTag(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase:   TestCase{ID: "T1"},
			Assertions: []string{"Valid"},
			Tags:       []string{"smoke", "fast"},
		},
		{
			TestCase:   TestCase{ID: "T2"},
			Assertions: []string{"Valid"},
			Tags:       []string{"integration"},
		},
		{
			TestCase:   TestCase{ID: "T3"},
			Assertions: []string{"Valid"},
			Tags:       []string{"smoke", "slow"},
		},
	}

	filtered := GetCustomTestsByTag(cases, "smoke")
	assert.Equal(t, 2, len(filtered))
	assert.Equal(t, "T1", filtered[0].ID)
	assert.Equal(t, "T3", filtered[1].ID)
}

func TestGetCustomTestsByValidator(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase:   TestCase{ID: "T1", ValidatorType: "TicketValidator"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T2", ValidatorType: "UserValidator"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T3", ValidatorType: "TicketValidator"},
			Assertions: []string{"Valid"},
		},
	}

	filtered := GetCustomTestsByValidator(cases, "TicketValidator")
	assert.Equal(t, 2, len(filtered))
	assert.Equal(t, "T1", filtered[0].ID)
	assert.Equal(t, "T3", filtered[1].ID)
}

func TestGetEnabledCustomTests(t *testing.T) {
	cases := []CustomTestCase{
		{TestCase: TestCase{ID: "T1", Enabled: true}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T2", Enabled: false}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T3", Enabled: true}, Assertions: []string{"Valid"}},
	}

	enabled := GetEnabledCustomTests(cases)
	assert.Equal(t, 2, len(enabled))
	assert.Equal(t, "T1", enabled[0].ID)
	assert.Equal(t, "T3", enabled[1].ID)
}

func TestMergeTestCases(t *testing.T) {
	standard := []TestCase{
		{ID: "STD1", MethodName: "GetStandard"},
		{ID: "STD2", MethodName: "SearchStandard"},
	}

	custom := []CustomTestCase{
		{TestCase: TestCase{ID: "CUST1", MethodName: "GetCustom"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "STD1", MethodName: "GetOverride"}, Assertions: []string{"Valid"}},
	}

	merged := MergeTestCases(standard, custom)

	assert.Equal(t, 3, len(merged))
	// Custom tests come first in order they appear
	assert.Equal(t, "GetCustom", merged[0].MethodName)
	assert.Equal(t, "GetOverride", merged[1].MethodName)
	assert.Equal(t, "SearchStandard", merged[2].MethodName)
}

func TestGetCustomTestsByDescription(t *testing.T) {
	cases := []CustomTestCase{
		{
			TestCase:   TestCase{ID: "T1", Description: "Test ticket retrieval by ID"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T2", Description: "Test user search by name"},
			Assertions: []string{"Valid"},
		},
		{
			TestCase:   TestCase{ID: "T3", Description: "Test ticket search with filters"},
			Assertions: []string{"Valid"},
		},
	}

	filtered := GetCustomTestsByDescription(cases, "ticket")
	assert.Equal(t, 2, len(filtered))
	assert.Equal(t, "T1", filtered[0].ID)
	assert.Equal(t, "T3", filtered[1].ID)
}

func TestCountCustomTestsByValidator(t *testing.T) {
	cases := []CustomTestCase{
		{TestCase: TestCase{ID: "T1", ValidatorType: "TicketValidator"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T2", ValidatorType: "UserValidator"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T3", ValidatorType: "TicketValidator"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T4", ValidatorType: "TicketValidator"}, Assertions: []string{"Valid"}},
	}

	counts := CountCustomTestsByValidator(cases)
	assert.Equal(t, 3, counts["TicketValidator"])
	assert.Equal(t, 1, counts["UserValidator"])
	assert.Equal(t, 0, counts["AssetValidator"])
}

func TestCountCustomTestsByService(t *testing.T) {
	cases := []CustomTestCase{
		{TestCase: TestCase{ID: "T1", ServiceType: "TicketService"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T2", ServiceType: "UserService"}, Assertions: []string{"Valid"}},
		{TestCase: TestCase{ID: "T3", ServiceType: "TicketService"}, Assertions: []string{"Valid"}},
	}

	counts := CountCustomTestsByService(cases)
	assert.Equal(t, 2, counts["TicketService"])
	assert.Equal(t, 1, counts["UserService"])
	assert.Equal(t, 0, counts["AssetService"])
}
