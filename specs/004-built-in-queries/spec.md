# Feature Specification: Built-in Query Methods for RT Data Types

**Feature Branch**: `004-built-in-queries`  
**Created**: February 5, 2026  
**Status**: Draft  
**Input**: User description: "Each data type (ticket, asset, user) should provide built in queries so a full custom query is not required every time. This should: allow for querying tickets subjects with like/contains (all tickets that contain the phrase Smartboard in their subject), allow for querying assets by Name, ID, or custom field (Internal Name for us), allow querying users by username, name or email. This will provide consistency in applications using the library with clear expectation of queries and returned data. Custom queries should be robustly supported with clear documentation on building an useful query"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Asset Queries by Standard Fields (Priority: P1)

A developer building an application needs to look up assets in their RT instance. They want to find assets by common identifiers: the asset name, the asset ID, or a custom field like "Internal Name" that their organization uses as a unique identifier. They should not need to construct complex JSON query structures for these common operations.

**Why this priority**: Asset lookups are the most frequently requested operation according to current usage patterns. Organizations use assets extensively for hardware tracking, and quick lookups by name or internal identifier are essential for daily operations.

**Independent Test**: Can be fully tested by calling `SearchByNameExact("laptop-001")` for exact lookups, `SearchByNamePartial("laptop")` for broader searches, `SearchByCustomFieldExact("Internal Name", "COMP-12345")` for exact custom field matches, or `SearchByCustomFieldPartial("Internal Name", "COMP")` for partial matches, verifying that the correct assets are returned with all details. Delivers immediate value for asset management applications.

**Acceptance Scenarios**:

1. **Given** an asset exists with Name "MacBook-1234", **When** developer calls `AssetService.SearchByNameExact("MacBook-1234")`, **Then** only that exact asset is returned with all fields populated
2. **Given** assets exist with names "MacBook-1234", "My MacBook", and "MacBook Pro", **When** developer calls `AssetService.SearchByNamePartial("MacBook")`, **Then** all three assets are returned
3. **Given** an asset exists with custom field "Internal Name" = "COMP-12345", **When** developer calls `AssetService.SearchByCustomFieldExact("Internal Name", "COMP-12345")`, **Then** only that exact asset is returned
4. **Given** assets exist with "Internal Name" values "COMP-123", "COMP-1234", and "COMP-12345", **When** developer calls `SearchByCustomFieldPartial("Internal Name", "COMP-123")`, **Then** all three assets are returned
5. **Given** no asset exists with exact custom field match, **When** developer calls `SearchByCustomFieldExact("Internal Name", "NonExistent")`, **Then** an empty SearchResult is returned with Total = 0

---

### User Story 2 - Ticket Queries by Subject Content (Priority: P2)

A developer building a help desk dashboard needs to find all tickets whose subject contains specific keywords. For example, they want to find all tickets about "Smartboard" issues, or all tickets mentioning "password reset". They need a simple method that handles the query syntax without requiring knowledge of TicketSQL.

**Why this priority**: Subject-based searching is a common requirement for reporting, dashboards, and automated ticket routing. While less frequent than asset lookups, it's essential for support workflows.

**Independent Test**: Can be fully tested by calling `SearchBySubject("Smartboard")` and verifying that all tickets with "Smartboard" in their subject are returned. Delivers value for ticket monitoring and reporting tools.

**Acceptance Scenarios**:

1. **Given** multiple tickets exist with "Smartboard" in their subject, **When** developer calls `TicketService.SearchBySubject("Smartboard")`, **Then** all matching tickets are returned
2. **Given** tickets exist with subjects "Password Reset Request" and "Password Change", **When** developer calls `SearchBySubject("Password")`, **Then** both tickets are returned
3. **Given** a ticket exists with subject in mixed case "SmartBoard Issues", **When** developer searches for "smartboard" (lowercase), **Then** the ticket is found (case-insensitive)
4. **Given** no tickets contain "XYZ123" in subject, **When** developer searches for "XYZ123", **Then** empty SearchResult is returned

---

### User Story 3 - User Queries by Identifier Fields (Priority: P3)

A developer building a user directory or authentication integration needs to look up RT users by common identifiers: username, full name, or email address. They need simple methods that abstract the RT query syntax and return consistent results.

**Why this priority**: User lookups are important for authentication, permissions, and user management features, but are typically less frequent than asset and ticket operations in most RT workflows.

**Independent Test**: Can be fully tested by calling exact match methods like `SearchByUsernameExact("jsmith")` or `SearchByEmailExact("jsmith@example.com")` for precise lookups, or partial match methods like `SearchByNamePartial("Smith")` for broader searches, verifying the correct users are returned. Delivers value for user management and authentication features.

**Acceptance Scenarios**:

1. **Given** a user exists with username "jsmith", **When** developer calls `UserService.SearchByUsernameExact("jsmith")`, **Then** only that user is returned with all fields populated
2. **Given** users exist with usernames "jsmith", "jsmith2", and "ajsmith", **When** developer calls `UserService.SearchByUsernamePartial("jsmith")`, **Then** all three users are returned
3. **Given** a user exists with email "jsmith@example.com", **When** developer calls `UserService.SearchByEmailExact("jsmith@example.com")`, **Then** only that user is returned
4. **Given** users exist with emails "john.smith@example.com" and "jane.smith@example.com", **When** developer calls `UserService.SearchByEmailPartial("smith@example.com")`, **Then** both users are returned
5. **Given** a user exists with name "John Smith", **When** developer calls `UserService.SearchByNameExact("John Smith")`, **Then** only that user is returned
6. **Given** multiple users have "Smith" in their name, **When** developer calls `UserService.SearchByNamePartial("Smith")`, **Then** all matching users are returned

---

### User Story 4 - Custom Query Support for Complex Requirements (Priority: P4)

A developer encounters a use case not covered by the built-in query methods. They need to construct a complex query with multiple conditions, custom fields, or specific RT query syntax. They should be able to use the existing `SearchWithCriteria` methods with clear documentation and examples showing how to build various query types.

**Why this priority**: While built-in methods should cover common cases, custom queries are essential for power users and complex applications. This is lower priority because the functionality already exists; we're primarily improving documentation.

**Independent Test**: Can be fully tested by consulting documentation, constructing a multi-condition query using `SearchWithCriteria`, and verifying results. Delivers value for advanced use cases and complex reporting requirements.

**Acceptance Scenarios**:

1. **Given** documentation includes examples of multi-field queries, **When** developer reads the documentation, **Then** they can construct a working custom query in under 10 minutes
2. **Given** developer needs to query assets by multiple custom fields, **When** they use `SearchWithCriteria` with documented syntax, **Then** the query executes successfully and returns expected results
3. **Given** documentation shows TicketSQL examples, **When** developer needs to query tickets by status AND subject, **Then** they can construct the query without trial-and-error
4. **Given** developer encounters a query error, **When** they check documentation, **Then** common error patterns and solutions are explained

---

## Clarifications

### Session 2026-02-05

- Q: Exact vs Partial Name Matching - Should SearchByName() use exact matching or partial/LIKE matching for asset names? → A: Separate methods - SearchByNameExact() and SearchByNamePartial()
- Q: Custom Field Matching Behavior - Should SearchByCustomField() use exact or partial matching for custom field values? → A: Separate methods - SearchByCustomFieldExact() and SearchByCustomFieldPartial()
- Q: User Name Matching Behavior - Should user search methods (SearchByName, SearchByUsername, SearchByEmail) use exact or partial matching? → A: Separate methods - Exact and partial variants for each (SearchByUsernameExact/Partial, SearchByEmailExact/Partial, SearchByNameExact/Partial)
- Q: Pagination for Large Result Sets - How should built-in query methods handle queries returning thousands of results? → A: Return first page only - Built-in methods return first page; user must handle pagination explicitly
- Q: Error Handling for Non-Existent Custom Fields - What happens when querying by a custom field that doesn't exist in the RT instance? → A: Return error - Return a specific error indicating the custom field doesn't exist in the RT instance

### Edge Cases

- What happens when a custom field name contains special characters (e.g., "Internal (Legacy) Name")?
- Queries that return thousands of results will only return the first page (RT API default page size). The SearchResult structure includes pagination information (total count, page info) that users can use to fetch additional pages via the generic Search or SearchWithCriteria methods.
- When querying by a custom field that doesn't exist in the RT instance, the method MUST return a specific error indicating the custom field was not found, allowing applications to handle the invalid field name appropriately.
- How are null or empty values in fields handled (e.g., asset with no Name set)?
- What happens when network timeout occurs during a search operation?
- How are special characters in search strings handled (e.g., searching for "laptop*" or "name with \" quote")?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: AssetService MUST provide a `SearchByNameExact(name string)` method that returns assets matching the exact name (case-insensitive equality)
- **FR-001b**: AssetService MUST provide a `SearchByNamePartial(query string)` method that returns assets whose name contains the query string (case-insensitive partial match)
- **FR-002**: AssetService MUST provide a `SearchByCustomFieldExact(fieldName string, value string)` method that queries any custom field with exact value matching
- **FR-002b**: AssetService MUST provide a `SearchByCustomFieldPartial(fieldName string, query string)` method that queries any custom field with partial/LIKE matching
- **FR-003**: TicketService MUST provide a `SearchBySubject(query string)` method that finds tickets whose subject contains the query string (case-insensitive)
- **FR-004**: UserService MUST provide `SearchByUsernameExact(username string)` and `SearchByUsernamePartial(query string)` methods for exact and partial username matching
- **FR-004b**: UserService MUST provide `SearchByEmailExact(email string)` and `SearchByEmailPartial(query string)` methods for exact and partial email matching
- **FR-004c**: UserService MUST provide `SearchByNameExact(name string)` and `SearchByNamePartial(query string)` methods for exact and partial full name matching
- **FR-005**: All built-in query methods MUST return the same `SearchResult[T]` structure used by existing search methods for consistency
- **FR-006**: All services MUST retain their existing `SearchWithCriteria` (or equivalent) methods for custom queries
- **FR-007**: Each built-in query method MUST handle empty results gracefully by returning an empty SearchResult with Total=0
- **FR-008**: Library documentation MUST include examples showing how to construct custom queries using `SearchWithCriteria`
- **FR-009**: Built-in methods MUST perform full detail fetching (not just summary data) to match existing search behavior
- **FR-010**: Query methods MUST handle special characters in search strings without causing errors
- **FR-011**: AssetService custom field methods MUST properly format custom field names with curly braces (e.g., "CustomField.{Internal Name}")
- **FR-012**: TicketService.SearchBySubject MUST use partial matching (LIKE/CONTAINS behavior) rather than exact matching
- **FR-013**: Each service's query methods MUST use the appropriate RT API syntax (JSON search for assets, TicketSQL for tickets, RT user query syntax for users)
- **FR-014**: All built-in query methods MUST return only the first page of results using RT API default pagination settings, with pagination metadata in the SearchResult structure for users to fetch subsequent pages if needed
- **FR-015**: Custom field search methods MUST return a specific error when the specified custom field does not exist in the RT instance, rather than returning empty results

### Key Entities

- **Asset Query Methods**: Built-in convenience methods for common asset lookup patterns (by name, by ID [already exists via Get], by custom field)
- **Ticket Query Methods**: Built-in convenience methods for subject-based searches and common ticket filtering patterns
- **User Query Methods**: Built-in convenience methods for user lookups by username, email, and name
- **Query Documentation**: Comprehensive examples and patterns for building custom queries when built-in methods are insufficient

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developer can query an asset by Name using a built-in method in 3 lines of code or less (including context and error handling)
- **SC-002**: Developer can query tickets by subject keyword in 3 lines of code or less
- **SC-003**: Developer can query users by username, email, or name in 3 lines of code or less
- **SC-004**: Built-in query methods cover at least 80% of common query use cases based on existing application patterns
- **SC-005**: Documentation includes at least 5 working examples of custom queries covering different complexity levels
- **SC-006**: All built-in query methods return the first page of results in under 2 seconds for typical RT instances
- **SC-007**: Error messages from failed queries provide clear indication of the problem (e.g., "Custom field 'Internal Name' not found")
- **SC-008**: Applications using the library require zero changes to maintain backward compatibility with existing Search methods
- **SC-009**: Developers successfully construct custom queries using documentation examples without requiring support assistance 90% of the time

## Scope

### In Scope

- Adding `SearchByNameExact` and `SearchByNamePartial` methods to AssetService for exact and partial name matching
- Adding `SearchByCustomFieldExact` and `SearchByCustomFieldPartial` methods to AssetService for exact and partial custom field matching
- Adding `SearchBySubject` method to TicketService
- Adding `SearchByUsernameExact`, `SearchByUsernamePartial`, `SearchByEmailExact`, `SearchByEmailPartial`, `SearchByNameExact`, and `SearchByNamePartial` methods to UserService
- Updating service interfaces in contracts package to include new methods
- Creating comprehensive documentation with custom query examples
- Adding unit tests for all new query methods
- Maintaining backward compatibility with existing Search/SearchWithCriteria methods

### Out of Scope

- Changing existing Search method behavior or signatures
- Adding query builder classes or fluent query APIs (future enhancement)
- Implementing query result caching or performance optimization beyond current implementation
- Adding query methods for other RT object types not yet supported (queues, groups, etc.)
- Real-time or streaming query results
- GraphQL or other alternative query interfaces

## Assumptions

- RT API supports the query syntax required for each method (JSON search for assets, TicketSQL for tickets, standard user queries)
- Custom field names in the RT instance follow standard naming conventions (may contain spaces and special characters)
- Developers using the library have basic familiarity with Go context handling and error checking patterns
- Developers will use the generic Search or SearchWithCriteria methods when pagination beyond the first page is required
- The existing `SearchWithCriteria` approach is sufficient for advanced query needs
- RT instances have reasonable response times (under 2 seconds for typical queries)

## Dependencies

- Existing AssetService, TicketService, and UserService implementations
- RT REST 2.0 API supporting required query operations
- Existing `SearchResult[T]` generic structure
- Service interfaces in contracts package
- Go 1.18+ for generics support

## Risks

- **Risk**: Custom field names vary widely across RT instances, making `SearchByCustomField` difficult to use consistently
  - **Mitigation**: Document the custom field name format requirements clearly; provide helper methods if needed

- **Risk**: RT query syntax limitations may prevent some built-in methods from working as expected
  - **Mitigation**: Test against real RT instances during development; document any known limitations

- **Risk**: Adding methods may bloat service interfaces
  - **Mitigation**: Keep methods focused on truly common use cases; avoid adding rarely-used convenience methods

- **Risk**: Performance issues if built-in methods always fetch full details
  - **Mitigation**: Monitor performance; consider adding lightweight variants if needed in future

## Open Questions

None - the feature scope and requirements are clear based on current usage patterns and user feedback.
