# Feature Specification: RT API Wrapper

**Feature Branch**: `001-rt-api-wrapper`
**Created**: 2026-02-04
**Status**: Draft
**Input**: User description: "rtutils_lib is a wrapper for Best Practical's Request Tracker 6.0.2 REST 2.0 api. It should cover all aspects related to users, tickets, and assets."

## Clarifications

### Session 2026-02-04

- Q: How should pagination be exposed to the library user? → A: **Manual Paging**. Return a `SearchResult` struct containing the slice of `Items` and metadata (e.g., `Page`, `Total`, `NextPageURL`). The user is responsible for making the subsequent call if they want more data.
- Q: How should Custom Fields (CFs) be typed in the Go structs? → A: **Dynamic Map** (`map[string]interface{}`). This allows handling both single-value (string) and multi-value (array) Custom Fields returned by RT without schema conflicts.
- Q: how should file attachments be handled in this version? → A: **Defer**. File attachment upload and download functionality is explicitly out of scope for this version to prioritize core ticket management stability.

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Ticket Management Integration (Priority: P1)

As a developer, I want to create, search, and update tickets via the wrapper library so that I can integrate RT ticket workflows into my applications.

**Why this priority**: Ticket management is the core functionality of Request Tracker and the primary reason for this library's existence.

**Independent Test**: Can be fully tested by a script that creates a ticket, searches for it by ID, updates a field, and adds a comment, validating the API response at each step.

**Acceptance Scenarios**:

1. **Given** valid credentials, **When** I request to create a ticket with valid fields (Subject, Queue), **Then** the ticket is created in RT and the new ticket ID is returned.
2. **Given** an existing ticket ID, **When** I request to fetch the ticket, **Then** I receive a Ticket object with correct attributes (Subject, Status, Owner, etc.).
3. **Given** an existing ticket, **When** I update the status to 'resolved', **Then** the change is reflected in RT.
4. **Given** search criteria (e.g., Owner='jdoe'), **When** I search for tickets, **Then** I receive a list of matching Ticket objects with support for pagination.
5. **Given** a ticket, **When** I call the comment or correspond method, **Then** the comment/reply is added to the ticket history.

---

### User Story 2 - User and Group Management (Priority: P2)

As an administrator using the library, I want to manage user accounts and group memberships programmatically so that I can sync user permissions with our central directory.

**Why this priority**: Managing who can access the efficient workflows is critical for security and operations.

**Independent Test**: Can be tested by creating a test user, retrieving them, adding them to a group, and verifying the group membership.

**Acceptance Scenarios**:

1. **Given** admin permissions, **When** I create a new user with required fields (Name, Email), **Then** the user is created in RT.
2. **Given** a user identifier (ID or Name), **When** I fetch user details, **Then** I receive the correct user information.
3. **Given** a user and a group, **When** I add the user to the group via the library, **Then** the user is listed as a member of that group.
4. **Given** a user, **When** I update their email address, **Then** the email is updated in RT.

---

### User Story 3 - Asset Management (Priority: P3)

As an IT manager using the library, I want to manage assets and their properties so that I can track inventory alongside support tickets.

**Why this priority**: Assets are a key secondary entity in RT for IT service management contexts.

**Independent Test**: Can be tested by creating an asset, updating a custom field value, and deleting the asset.

**Acceptance Scenarios**:

1. **Given** valid permissions, **When** I create an asset with Name and Catalog, **Then** the asset is successfully created.
2. **Given** an asset ID, **When** I fetch the asset, **Then** I receive the Asset object with its details.
3. **Given** search criteria (AssetSQL), **When** I search for assets, **Then** I receive a list of matching Asset objects.
4. **Given** an asset, **When** I request to delete it, **Then** it is marked as deleted in RT.

### Edge Cases

- What happens when the API token is invalid or expired? (Should raise an AuthenticationError)
- What happens when a requested entity (Ticket/User/Asset) does not exist? (Should raise a NotFoundError)
- How does the system handle API rate limiting or temporary network failures? (Should handle connection errors gracefully)
- What happens if validation fails on the creating/updating side (e.g. missing required custom field)? (Should propagate the API error message)

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: The library MUST provide a client class to handle authentication (Token-based) and base URL configuration for RT 6.0.2+.
- **FR-002**: The library MUST implement `Ticket` resource methods: create, get, search (TicketSQL), update, delete, history, comment, correspond, take, untake, steal.
- **FR-003**: The library MUST implement `User` resource methods: create, get (by ID/Name), search, update, delete (disable), history.
- **FR-004**: The library MUST implement `Asset` resource methods: create, get, search, update, delete, history.
- **FR-005**: The library MUST support `Group` membership operations for users (add/remove/list).
- **FR-006**: The library MUST support bulk operations for Tickets (bulk create, bulk update, bulk comment/correspond).
- **FR-007**: The library MUST handle custom fields for Tickets and Assets using a dynamic map structure (`map[string]interface{}`) to support both single-value and multi-value fields during create, update, and retrieval operations.
- **FR-008**: The library MUST provide manual pagination support. Search methods MUST return a result object containing the current page's items and metadata (page number, total count, next page URL) to allow the caller to request subsequent pages.
- **FR-009**: The library MUST parse JSON responses from the RT REST 2.0 API into structured objects or native data types.
- **FR-010**: The library MUST raise structured exceptions for HTTP 4xx and 5xx error responses.

### Key Entities _(include if feature involves data)_

- **Ticket**: Represents an RT ticket (id, subject, queue, owner, status, priority, custom fields).
- **User**: Represents an RT user (id, name, email, realname, disabled status).
- **Asset**: Represents an RT asset (id, name, catalog, status, custom fields).
- **Transaction**: Represents a history entry for an object.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: Developers can perform full CRUD lifecycle for Tickets, Users, and Assets using the library without raw HTTP calls.
- **SC-002**: Authentication works correctly with a standard RT 6.0.2 REST 2.0 API token.
- **SC-003**: 100% of the public API methods specified in requirements have corresponding test cases passing against a mock or dev instance.
- **SC-004**: Library handles all defined error scenarios (404, 401) by raising appropriate exceptions.

## Assumptions

- The target RT instance is version 6.0.2 or compatible.
- The REST 2.0 API is enabled on the target instance.
- The user has a valid API token with sufficient permissions for the requested operations.
- **Out of Scope**: File attachment uploads and downloads are not supported in this version.
