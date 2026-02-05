# Tasks: Built-in Query Methods for RT Data Types

**Input**: Design documents from `/specs/004-built-in-queries/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Required (per constitution). Tests must be written before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and shared test scaffolding

- [x] T001 Create shared httpmock test helpers in test_helpers_test.go
- [x] T002 [P] Verify test dependencies in go.mod/go.sum include testify and httpmock

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core changes that MUST be complete before ANY user story can be implemented

- [x] T003 Update service contracts with new query methods in contracts/interfaces.go
- [x] T004 Add custom field not found error type and input validation helpers in types.go

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Asset Queries by Standard Fields (Priority: P1) 

**Goal**: Provide exact and partial asset queries by name and custom field.

**Independent Test**: Call `SearchByNameExact`, `SearchByNamePartial`, `SearchByCustomFieldExact`, and `SearchByCustomFieldPartial` and validate results and error behavior.

### Tests for User Story 1 (TDD)

- [x] T005 [US1] Add contract tests for asset query methods in asset_test.go

### Implementation for User Story 1

- [x] T006 [US1] Implement asset query methods and custom field error mapping in asset.go

**Checkpoint**: Asset queries functional and independently testable

---

## Phase 4: User Story 2 - Ticket Queries by Subject Content (Priority: P2)

**Goal**: Provide subject-based ticket search using partial matching.

**Independent Test**: Call `SearchBySubject` with sample keywords and validate returned tickets.

### Tests for User Story 2 (TDD)

- [x] T007 [US2] Add contract tests for subject search in ticket_test.go

### Implementation for User Story 2

- [x] T008 [US2] Implement subject search method in ticket.go

**Checkpoint**: Ticket subject queries functional and independently testable

---

## Phase 5: User Story 3 - User Queries by Identifier Fields (Priority: P3)

**Goal**: Provide exact and partial user queries for username, email, and name.

**Independent Test**: Call exact and partial user query methods and validate returned users.

### Tests for User Story 3 (TDD)

- [x] T009 [US3] Add contract tests for user query methods in user_test.go

### Implementation for User Story 3

- [x] T010 [US3] Implement user query methods in user.go

**Checkpoint**: User queries functional and independently testable

---

## Phase 6: User Story 4 - Custom Query Support Documentation (Priority: P4)

**Goal**: Provide clear documentation for built-in and custom query usage.

**Independent Test**: Follow README and quickstart examples to construct built-in and custom queries without additional context.

### Documentation for User Story 4

- [x] T011 [P] [US4] Update built-in query method docs and examples in README.md
- [x] T012 [P] [US4] Validate and refine quickstart examples in specs/004-built-in-queries/quickstart.md

**Checkpoint**: Documentation supports built-in and custom query usage

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup, formatting, and final verification

- [x] T013 [P] Run gofmt on asset.go, ticket.go, user.go, types.go
- [x] T014 Run go test ./... and fix any failures in *_test.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories
- **User Stories (Phases 3-6)**: Depend on Foundational phase completion
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational - no dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational - no dependencies on other stories
- **User Story 3 (P3)**: Can start after Foundational - no dependencies on other stories
- **User Story 4 (P4)**: Can start after Foundational - depends on API names stabilized by US1-US3

### Parallel Opportunities

- T002 can run in parallel with T001
- T011 and T012 can run in parallel
- T013 can run in parallel with documentation tasks

---

## Parallel Example: User Story 4

```bash
Task: "Update built-in query method docs and examples in README.md"
Task: "Validate and refine quickstart examples in specs/004-built-in-queries/quickstart.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. Stop and validate asset query methods independently

### Incremental Delivery

1. Setup + Foundational
2. User Story 1 (Assets)
3. User Story 2 (Tickets)
4. User Story 3 (Users)
5. User Story 4 (Documentation)
6. Polish & verification
