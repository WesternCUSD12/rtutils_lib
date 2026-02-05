# Tasks: RT API Wrapper

**Feature**: RT API Wrapper
**Branch**: `001-rt-api-wrapper`
**Spec**: [spec.md](./spec.md)
**Plan**: [plan.md](./plan.md)

## Phase 1: Setup

- [x] T001 Initialize Go module `rtutils_lib` if not present
- [x] T002 Create initial project structure (`contracts`, `types`, etc.) per plan
- [x] T003 [P] Add `testify` dependency for testing
- [x] T004 Define `Client` struct with base URL and token fields in `client.go`
- [x] T005 Implement `NewClient` factory function in `client.go`

## Phase 2: Foundational & Contracts

- [x] T006 Define shared error types (`APIError` struct) in `types.go`
- [x] T007 Implement shared HTTP request helper method (accepting `context.Context`, handling JSON/Auth) in `client.go`
- [x] T008 Define `SearchResult` generic struct in `types.go` with simple pagination logic
- [x] T009 Define `ActionResult` struct for operation responses in `types.go`
- [x] T010 Define Go interfaces for `TicketService`, `UserService`, and `AssetService` in `contracts/` package (Interface-First)

## Phase 3: Ticket Management (User Story 1 - P1)

- [x] T011 [US1] Define `Ticket` struct with core fields and `map[string]interface{}` custom fields in `ticket.go`
- [x] T012 [P] [US1] Create unit test suite for Ticket operations in `ticket_test.go`
- [x] T013 [US1] Implement `Create` method for Tickets (taking `context.Context`) in `ticket.go`
- [x] T014 [US1] Implement `Get` ticket by ID method (taking `context.Context`) in `ticket.go`
- [x] T015 [US1] Implement `Search` tickets method (taking `context.Context`) with TicketSQL and pagination support
- [x] T016 [P] [US1] Implement `Update` ticket method (taking `context.Context`) in `ticket.go`
- [x] T017 [US1] Implement `Delete` ticket method (taking `context.Context`) in `ticket.go`
- [x] T018 [US1] Implement `GetHistory` method for Tickets (taking `context.Context`) to retrieve transactions
- [ ] T019 [US1] Implement `Comment` and `Correspond` methods (taking `context.Context`) in `ticket.go`
- [ ] T020 [US1] Implement ticket actions (`Take`, `Untake`, `Steal`) (taking `context.Context`) in `ticket.go`
- [ ] T021 [US1] Implement `BulkCreate` and `BulkUpdate` methods (taking `context.Context`) in `ticket.go`

## Phase 4: User & Group Management (User Story 2 - P2)

- [x] T022 [US2] Define `User` struct in `user.go`
- [x] T023 [P] [US2] Create unit test suite for User operations in `user_test.go`
- [x] T024 [US2] Implement `Create` user method (taking `context.Context`) in `user.go`
- [x] T025 [US2] Implement `Get` user (by ID/Name) method (taking `context.Context`) in `user.go`
- [x] T026 [US2] Implement `Search` users method (taking `context.Context`) in `user.go`
- [x] T027 [P] [US2] Implement `Update` user method (taking `context.Context`) in `user.go`
- [x] T028 [US2] Implement `Disable` user method (taking `context.Context`) in `user.go`
- [x] T029 [US2] Implement `GetHistory` method for Users (taking `context.Context`)
- [x] T030 [US2] Implement Group membership methods (`GetGroupMemberships`, `AddToGroup`, `RemoveFromGroup`) (taking `context.Context`) in `user.go`

## Phase 5: Asset Management (User Story 3 - P3)

- [x] T031 [US3] Define `Asset` struct in `asset.go`
- [x] T032 [P] [US3] Create unit test suite for Asset operations in `asset_test.go`
- [x] T033 [US3] Implement `Create` asset method (taking `context.Context`) in `asset.go`
- [x] T034 [US3] Implement `Get` asset method (taking `context.Context`) in `asset.go`
- [x] T035 [US3] Implement `Search` assets method (taking `context.Context`) with AssetSQL in `asset.go`
- [x] T036 [P] [US3] Implement `Update` asset method (taking `context.Context`) in `asset.go`
- [x] T037 [US3] Implement `Delete` asset method (taking `context.Context`) in `asset.go`

## Phase 6: Polish

- [x] T038 Verify all `_url` fields are populated correctly in responses
- [x] T039 Review and clean up `APIError` messages for better debugging
- [x] T040 Ensure all public methods are documented with comments
- [x] T041 Run full test suite and ensure no regressions

## Dependencies

- Phase 1 & 2 blocks all User Stories
- User Story 1, 2, 3 are largely independent, but share the common `Client` infrastructure.
- Interface definitions (T010) must complete before specific implementations.

## Implementation Strategy

- MVP: Client Setup + Contracts + Ticket Create/Get (T001-T014)
- Incremental: Add Search/Update for Tickets, then move to Users and Assets.
