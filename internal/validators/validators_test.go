package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenericValidatorName(t *testing.T) {
	gv := &GenericValidator{ValidatorName: "TestValidator"}
	assert.Equal(t, "TestValidator", gv.Name())
}

func TestGenericValidatorValidateNil(t *testing.T) {
	gv := &GenericValidator{ValidatorName: "Generic"}
	err := gv.Validate(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "result is nil")
}

func TestGenericValidatorValidateNotNil(t *testing.T) {
	gv := &GenericValidator{ValidatorName: "Generic"}
	err := gv.Validate("some result")
	assert.NoError(t, err)
}

func TestTicketValidatorName(t *testing.T) {
	tv := &TicketValidator{}
	assert.Equal(t, "TicketValidator", tv.Name())
}

func TestTicketValidatorValidateNil(t *testing.T) {
	tv := &TicketValidator{}
	err := tv.Validate(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ticket is nil")
}

func TestTicketValidatorValidateWithRequiredFields(t *testing.T) {
	tv := &TicketValidator{}
	ticket := map[string]interface{}{
		"id":      "123",
		"status":  "open",
		"subject": "Test ticket",
	}
	err := tv.Validate(ticket)
	assert.NoError(t, err)
}

func TestTicketValidatorValidateMissingID(t *testing.T) {
	tv := &TicketValidator{}
	ticket := map[string]interface{}{
		"status": "open",
	}
	err := tv.Validate(ticket)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestTicketValidatorValidateMissingStatus(t *testing.T) {
	tv := &TicketValidator{}
	ticket := map[string]interface{}{
		"id": "123",
	}
	err := tv.Validate(ticket)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status")
}

func TestTicketValidatorValidateInvalidType(t *testing.T) {
	tv := &TicketValidator{}
	err := tv.Validate("not a map")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a map")
}

func TestUserValidatorName(t *testing.T) {
	uv := &UserValidator{}
	assert.Equal(t, "UserValidator", uv.Name())
}

func TestUserValidatorValidateNil(t *testing.T) {
	uv := &UserValidator{}
	err := uv.Validate(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "user is nil")
}

func TestUserValidatorValidateWithRequiredFields(t *testing.T) {
	uv := &UserValidator{}
	user := map[string]interface{}{
		"id":    "user123",
		"name":  "John Doe",
		"email": "john@example.com",
	}
	err := uv.Validate(user)
	assert.NoError(t, err)
}

func TestUserValidatorValidateMissingID(t *testing.T) {
	uv := &UserValidator{}
	user := map[string]interface{}{
		"name": "John Doe",
	}
	err := uv.Validate(user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestUserValidatorValidateMissingName(t *testing.T) {
	uv := &UserValidator{}
	user := map[string]interface{}{
		"id": "user123",
	}
	err := uv.Validate(user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestUserValidatorValidateInvalidType(t *testing.T) {
	uv := &UserValidator{}
	err := uv.Validate([]string{"not", "a", "map"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a map")
}

func TestAssetValidatorName(t *testing.T) {
	av := &AssetValidator{}
	assert.Equal(t, "AssetValidator", av.Name())
}

func TestAssetValidatorValidateNil(t *testing.T) {
	av := &AssetValidator{}
	err := av.Validate(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "asset is nil")
}

func TestAssetValidatorValidateWithRequiredFields(t *testing.T) {
	av := &AssetValidator{}
	asset := map[string]interface{}{
		"id":          "asset456",
		"name":        "Laptop",
		"description": "Dell Latitude",
	}
	err := av.Validate(asset)
	assert.NoError(t, err)
}

func TestAssetValidatorValidateMissingID(t *testing.T) {
	av := &AssetValidator{}
	asset := map[string]interface{}{
		"name": "Laptop",
	}
	err := av.Validate(asset)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestAssetValidatorValidateMissingName(t *testing.T) {
	av := &AssetValidator{}
	asset := map[string]interface{}{
		"id": "asset456",
	}
	err := av.Validate(asset)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestAssetValidatorValidateInvalidType(t *testing.T) {
	av := &AssetValidator{}
	err := av.Validate(123)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a map")
}

func TestValidationErrorInterface(t *testing.T) {
	ve := NewValidationError("TestValidator", "test error message")
	assert.Equal(t, "test error message", ve.Error())
	assert.Equal(t, "TestValidator", ve.Validator)
}

func TestGetValidatorTicketValidator(t *testing.T) {
	validator := GetValidator("TicketValidator")
	assert.NotNil(t, validator)
	assert.Equal(t, "TicketValidator", validator.Name())
	_, ok := validator.(*TicketValidator)
	assert.True(t, ok)
}

func TestGetValidatorUserValidator(t *testing.T) {
	validator := GetValidator("UserValidator")
	assert.NotNil(t, validator)
	assert.Equal(t, "UserValidator", validator.Name())
	_, ok := validator.(*UserValidator)
	assert.True(t, ok)
}

func TestGetValidatorAssetValidator(t *testing.T) {
	validator := GetValidator("AssetValidator")
	assert.NotNil(t, validator)
	assert.Equal(t, "AssetValidator", validator.Name())
	_, ok := validator.(*AssetValidator)
	assert.True(t, ok)
}

func TestGetValidatorUnknown(t *testing.T) {
	validator := GetValidator("UnknownValidator")
	assert.NotNil(t, validator)
	assert.Equal(t, "UnknownValidator", validator.Name())
	_, ok := validator.(*GenericValidator)
	assert.True(t, ok)
}

func TestValidatorResultValidator(t *testing.T) {
	// Ensure all validators implement ResultValidator interface
	var validators []ResultValidator

	validators = append(validators, &GenericValidator{})
	validators = append(validators, &TicketValidator{})
	validators = append(validators, &UserValidator{})
	validators = append(validators, &AssetValidator{})

	for _, v := range validators {
		assert.NotNil(t, v.Name())
		// All should handle nil without panic
		_ = v.Validate(nil)
	}
}
