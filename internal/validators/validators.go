package validators

// ResultValidator defines the interface for validating test results
type ResultValidator interface {
	// Validate checks if a result contains meaningful/expected data
	// Returns error if validation fails with descriptive error message
	Validate(result interface{}) error

	// Name returns the name of this validator
	Name() string
}

// GenericValidator implements ResultValidator for basic result checks
type GenericValidator struct {
	ValidatorName string
}

func (g *GenericValidator) Name() string {
	return g.ValidatorName
}

func (g *GenericValidator) Validate(result interface{}) error {
	// Basic check: result is not nil
	if result == nil {
		return NewValidationError("GenericValidator", "result is nil")
	}
	return nil
}

// TicketValidator implements ResultValidator for Ticket types
type TicketValidator struct{}

func (tv *TicketValidator) Name() string {
	return "TicketValidator"
}

func (tv *TicketValidator) Validate(result interface{}) error {
	if result == nil {
		return NewValidationError("TicketValidator", "ticket is nil")
	}

	ticketMap, ok := result.(map[string]interface{})
	if !ok {
		return NewValidationError("TicketValidator", "result is not a map")
	}

	if _, hasID := ticketMap["id"]; !hasID {
		return NewValidationError("TicketValidator", "missing required field: id")
	}

	if _, hasStatus := ticketMap["status"]; !hasStatus {
		return NewValidationError("TicketValidator", "missing required field: status")
	}

	return nil
}

// UserValidator implements ResultValidator for User types
type UserValidator struct{}

func (uv *UserValidator) Name() string {
	return "UserValidator"
}

func (uv *UserValidator) Validate(result interface{}) error {
	if result == nil {
		return NewValidationError("UserValidator", "user is nil")
	}

	userMap, ok := result.(map[string]interface{})
	if !ok {
		return NewValidationError("UserValidator", "result is not a map")
	}

	if _, hasID := userMap["id"]; !hasID {
		return NewValidationError("UserValidator", "missing required field: id")
	}

	if _, hasName := userMap["name"]; !hasName {
		return NewValidationError("UserValidator", "missing required field: name")
	}

	return nil
}

// AssetValidator implements ResultValidator for Asset types
type AssetValidator struct{}

func (av *AssetValidator) Name() string {
	return "AssetValidator"
}

func (av *AssetValidator) Validate(result interface{}) error {
	if result == nil {
		return NewValidationError("AssetValidator", "asset is nil")
	}

	assetMap, ok := result.(map[string]interface{})
	if !ok {
		return NewValidationError("AssetValidator", "result is not a map")
	}

	if _, hasID := assetMap["id"]; !hasID {
		return NewValidationError("AssetValidator", "missing required field: id")
	}

	if _, hasName := assetMap["name"]; !hasName {
		return NewValidationError("AssetValidator", "missing required field: name")
	}

	return nil
}

// ValidationError represents a validation failure
type ValidationError struct {
	Validator string
	Message   string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func NewValidationError(validator, message string) *ValidationError {
	return &ValidationError{
		Validator: validator,
		Message:   message,
	}
}

// GetValidator returns the appropriate validator for the given type
func GetValidator(validatorType string) ResultValidator {
	switch validatorType {
	case "TicketValidator":
		return &TicketValidator{}
	case "UserValidator":
		return &UserValidator{}
	case "AssetValidator":
		return &AssetValidator{}
	default:
		return &GenericValidator{ValidatorName: validatorType}
	}
}
