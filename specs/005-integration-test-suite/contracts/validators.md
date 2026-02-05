# Validator Specifications

**Date**: February 5, 2026  
**Phase**: Phase 1 (Design)

## Validator Interface Contract

All validators must implement this interface:

```go
type ResultValidator interface {
    // Validate checks if a result contains meaningful data
    // Returns (isValid bool, errorMessage string)
    Validate(result interface{}) (bool, string)
    
    // RequiredFields returns the list of field names that must be present
    RequiredFields() []string
}
```

## Validator Implementation Patterns

### Type Assertion Pattern

```go
func (v *TicketValidator) Validate(result interface{}) (bool, string) {
    ticket, ok := result.(*Ticket)
    if !ok {
        return false, "Result is not a Ticket type"
    }
    
    // Proceed with validation
    // ...
}
```

### Required Field Check Pattern

```go
if ticket.ID == "" {
    return false, "Required field 'id' is empty"
}
```

### Additional Field Count Pattern

```go
additionalFields := 0
if ticket.Subject != "" {
    additionalFields++
}
if ticket.Description != "" {
    additionalFields++
}
// ... check other fields

if additionalFields == 0 {
    return false, "No fields populated beyond required id and status"
}
```

## TicketValidator Specification

### Purpose
Validates that a Ticket result contains meaningful data for bug tracking workflows.

### Type
- Input: `*Ticket` (single) or `*SearchResult[Ticket]` (collection)
- Handles both direct Get results and Search results

### Required Fields
1. **id** (string)
   - Must be non-empty
   - Primary identifier
   - Error if empty: "Required field 'id' is empty"

2. **status** (string)
   - Must be non-empty
   - Indicates ticket state (new, open, resolved, etc.)
   - Error if empty: "Required field 'status' is empty"

### Additional Fields Requirement
Must have at least ONE of:
- `subject` (non-empty)
- `description` (non-empty)
- `queue` (non-empty)
- `owner` (non-empty)
- `priority` (non-empty)
- `created` (non-empty)
- `updated` (non-empty)
- Other custom fields

Error if none: "No fields populated beyond required id and status"

### Validation Algorithm

```
Input: interface{} (expected to be *Ticket or *SearchResult[Ticket])

1. Type assert to *Ticket
   IF not possible:
      RETURN (false, "Result is not a Ticket type")

2. Check required field: id
   IF id == "":
      RETURN (false, "Required field 'id' is empty")

3. Check required field: status
   IF status == "":
      RETURN (false, "Required field 'status' is empty")

4. Count additional populated fields
   additional_count = 0
   FOR each optional_field in [subject, description, queue, owner, priority, created, updated]:
      IF optional_field != "":
         additional_count++

5. Check additional field requirement
   IF additional_count == 0:
      RETURN (false, "No fields populated beyond required id and status")

6. All checks passed
   RETURN (true, "")
```

### Search Result Handling
- When validating SearchResult[Ticket], validate each Ticket in Items array
- If ANY ticket fails validation, entire result fails
- All must pass for overall PASS

## UserValidator Specification

### Purpose
Validates that a User result has essential identity information.

### Type
- Input: `*User` (single) or `*SearchResult[User]` (collection)
- Handles both direct Get results and Search results

### Required Fields
1. **id** (string)
   - Must be non-empty
   - Username or ID
   - Error if empty: "Required field 'id' is empty"

2. **name** (string)
   - Must be non-empty
   - Real name or display name
   - Error if empty: "Required field 'name' is empty"

### No Additional Fields Requirement
Unlike Ticket, User only requires the 2 core fields for basic functionality.

### Validation Algorithm

```
Input: interface{} (expected to be *User or *SearchResult[User])

1. Type assert to *User
   IF not possible:
      RETURN (false, "Result is not a User type")

2. Check required field: id
   IF id == "":
      RETURN (false, "Required field 'id' is empty")

3. Check required field: name
   IF name == "":
      RETURN (false, "Required field 'name' is empty")

4. All checks passed
   RETURN (true, "")
```

### Search Result Handling
- When validating SearchResult[User], validate each User in Items array
- If ANY user fails validation, entire result fails
- All must pass for overall PASS

## AssetValidator Specification

### Purpose
Validates that an Asset result has essential inventory information.

### Type
- Input: `*Asset` (single) or `*SearchResult[Asset]` (collection)
- Handles both direct Get results and Search results

### Required Fields
1. **id** (string)
   - Must be non-empty
   - Asset identifier
   - Error if empty: "Required field 'id' is empty"

2. **name** (string)
   - Must be non-empty
   - Asset descriptive name
   - Error if empty: "Required field 'name' is empty"

### No Additional Fields Requirement
Similar to User, Asset only requires 2 core identity fields.

### Validation Algorithm

```
Input: interface{} (expected to be *Asset or *SearchResult[Asset])

1. Type assert to *Asset
   IF not possible:
      RETURN (false, "Result is not an Asset type")

2. Check required field: id
   IF id == "":
      RETURN (false, "Required field 'id' is empty")

3. Check required field: name
   IF name == "":
      RETURN (false, "Required field 'name' is empty")

4. All checks passed
   RETURN (true, "")
```

### Search Result Handling
- When validating SearchResult[Asset], validate each Asset in Items array
- If ANY asset fails validation, entire result fails
- All must pass for overall PASS

## GenericValidator Specification

### Purpose
Fallback validator for unknown or future entity types.

### Type
- Input: Any interface{} result

### Required Fields
1. **id** (extracted from result)
   - Must be non-empty
   - Error if not present or empty

### Additional Fields Requirement
Must have at least ONE additional field populated (beyond id)

### Validation Algorithm

```
Input: interface{} (unknown type)

1. IF result has 'id' field:
      IF id is empty or missing:
         RETURN (false, "Required field 'id' is empty")
   ELSE:
      RETURN (false, "No 'id' field found in result")

2. Count non-id fields
   additional_count = 0
   FOR each field in result:
      IF field_name != "id" AND field is populated:
         additional_count++

3. Check additional field requirement
   IF additional_count == 0:
      RETURN (false, "No fields populated beyond required id")

4. All checks passed
   RETURN (true, "")
```

## Error Message Standards

All error messages must follow this pattern:

**For missing required field:**
```
"Required field '{field_name}' is empty"
```

Example: `"Required field 'status' is empty"`

**For missing additional field:**
```
"No fields populated beyond required {field_list}"
```

Example: `"No fields populated beyond required id and status"`

**For type mismatch:**
```
"Result is not a {ExpectedType} type"
```

Example: `"Result is not a Ticket type"`

## Edge Cases & Handling

### Empty String Values
Empty strings (`""`) are treated as missing/unpopulated fields.
- `null` in Go: handled by type assertion (would fail)
- Empty string: fails validation
- Whitespace-only strings: treated as valid (not trimmed)

### Zero/Null Numbers
- `0` for numeric IDs: considered populated (not empty)
- `nil` pointer: type assertion fails

### Collection Handling (SearchResult[T])
- Empty collection (no items): passes validation if Items array is properly structured
- Collection with some invalid items: entire result fails (all-or-nothing)
- Collection with all valid items: passes

### Custom Fields
- CustomFields map: ignored in validation (too dynamic to standardize)
- Only check structured fields (id, status, name, etc.)
- Can be extended in Phase 2 with configurable validators

## Implementation Guidelines for Phase 2

1. **Test-First Approach**
   - Write validator unit tests BEFORE implementing validators
   - Test with various edge cases (empty strings, missing fields, etc.)
   - Use `testify/assert` for clear assertions

2. **No External Dependencies**
   - Validators use only stdlib and rtutils_lib types
   - No reflection required (use type assertions)
   - Keep code simple and fast

3. **Performance**
   - Validators should complete in <1ms
   - No network calls or I/O
   - Pure in-memory validation only

4. **Error Messages**
   - Include field name in error
   - Include what was expected
   - Include what was actual
   - Keep messages under 200 characters

5. **Extensibility**
   - Consider a ValidatorFactory for dynamic validator selection
   - Allow registering custom validators for Phase 3
   - Validate validator itself (ensure required methods work)

## Phase 2 Test Cases

### TicketValidator Unit Tests

```go
func TestTicketValidator(t *testing.T) {
    v := &TicketValidator{}
    
    // Happy path
    valid, msg := v.Validate(&Ticket{ID: "1", Status: "new", Subject: "test"})
    assert.True(t, valid)
    
    // Missing id
    valid, msg := v.Validate(&Ticket{ID: "", Status: "new", Subject: "test"})
    assert.False(t, valid)
    assert.Contains(t, msg, "id")
    
    // Missing additional field
    valid, msg := v.Validate(&Ticket{ID: "1", Status: "new"})
    assert.False(t, valid)
    assert.Contains(t, msg, "No fields populated")
}
```

### UserValidator Unit Tests

```go
func TestUserValidator(t *testing.T) {
    v := &UserValidator{}
    
    // Happy path
    valid, msg := v.Validate(&User{ID: "admin", Name: "Administrator"})
    assert.True(t, valid)
    
    // Missing id
    valid, msg := v.Validate(&User{ID: "", Name: "Administrator"})
    assert.False(t, valid)
    
    // Missing name
    valid, msg := v.Validate(&User{ID: "admin", Name: ""})
    assert.False(t, valid)
}
```

## FAQ

**Q: Why require additional fields for Tickets but not Users/Assets?**  
A: Tickets often have extensive metadata; this catches incomplete expansions. Users/Assets have simpler requirements, so just id+name is sufficient.

**Q: Why not validate all fields on an entity?**  
A: To be lenient and catch only broken methods, not over-enforce. Can always add stricter validators later.

**Q: Should validators check data types?**  
A: Type assertion handles Go's type checking. Runtime JSON parsing ensures types match struct definitions.

**Q: Can validators be stateful?**  
A: No. Validators are stateless (no receiver fields except receiver pointer). Each call is independent.

**Q: What about custom fields?**  
A: Phase 1 ignores custom fields. Phase 2 can support configurable field requirements.
