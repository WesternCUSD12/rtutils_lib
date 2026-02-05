# Tasks: Ticket Search Example

**Input**: Design documents from `/specs/001-ticket-search/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Included (Constitution requires test-first).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Phase 1: Setup (Shared Infrastructure)

- [x] T001 Create examples/ticket_search/ directory and main.go skeleton in examples/ticket_search/main.go

---

## Phase 2: Foundational (Blocking Prerequisites)

- [x] T002 [P] Update TicketService contract with Search per-page and GetByURL in contracts/interfaces.go
- [x] T003 [P] Add failing unit tests for TicketService Search per-page and GetByURL in ticket_test.go
- [x] T004 Update TicketService.Search signature and per_page query param handling in ticket.go
- [x] T005 Implement TicketService.GetByURL using client.request in ticket.go
- [x] T006 Update any Search call sites to new signature (if any) in examples/ or tests

---

## Phase 3: User Story 1 - Find Tickets by Filters (Priority: P1) 🎯 MVP

**Goal**: Search tickets by Queue, Status, Owner, Requestor, and Subject keyword filters.
**Independent Test**: Run search with filters and verify returned tickets match criteria.

### Tests for User Story 1

- [x] T007 [P] [US1] Add validation tests for filter requirements and pagination bounds in examples/ticket_search/main.go (table-driven helper tests if created)

### Implementation for User Story 1

- [x] T008 [US1] Implement loadEnv and flag parsing for queue/status/owner/requestor/subject/page/per-page in examples/ticket_search/main.go
- [x] T009 [US1] Validate at least one filter and enforce per-page default (20) and max (100) in examples/ticket_search/main.go
- [x] T010 [US1] Build TicketSQL query with AND semantics and exact/keyword matching in examples/ticket_search/main.go
- [x] T011 [US1] Execute search and render tabular results with ID/Subject/Status/Queue/Owner/Requestor/URL in examples/ticket_search/main.go
- [x] T012 [US1] Handle empty results with clear message and exit code 2 in examples/ticket_search/main.go

---

## Phase 4: User Story 2 - Navigate Search Results (Priority: P2)

**Goal**: Page through search results with predictable defaults and boundaries.
**Independent Test**: Request page 1 and page 2 for the same query and verify unique results and accurate paging metadata.

### Implementation for User Story 2

- [x] T013 [US2] Add pagination summary output (page, per-page, total, pages, next page URL) in examples/ticket_search/main.go
- [x] T014 [US2] Detect out-of-range page requests and return exit code 2 with a clear message in examples/ticket_search/main.go

---

## Phase 5: User Story 3 - View Ticket Details from Search (Priority: P3)

**Goal**: Hydrate and display full ticket details from a search result.
**Independent Test**: Provide --details-id for a ticket in results and verify full details are displayed.

### Implementation for User Story 3

- [x] T015 [US3] Add --details-id handling to select a ticket result and resolve URL when available in examples/ticket_search/main.go
- [x] T016 [US3] Implement detail view output for a single ticket (vertical key-value) in examples/ticket_search/main.go

---

## Phase 6: Polish & Cross-Cutting Concerns

- [x] T017 [P] Update quickstart usage examples if needed in specs/001-ticket-search/quickstart.md
- [x] T018 Ensure usage text and exit codes (0 success, 1 error, 2 no results/not found) are consistent in examples/ticket_search/main.go
- [x] T019 Run gofmt on examples/ticket_search/main.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on Setup completion
- **User Stories (Phase 3+)**: Depend on Foundational completion
- **Polish (Phase 6)**: Depends on all desired user stories

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational; no dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational; independent but builds on search output
- **User Story 3 (P3)**: Can start after Foundational; uses search results and detail hydration

### Parallel Opportunities

- T002 and T003 can run in parallel
- T013 and T014 can run in parallel after US1 is complete

---

## Parallel Example: User Story 1

- Task: "Add validation tests for filter requirements and pagination bounds in examples/ticket_search/main.go"
- Task: "Build TicketSQL query with AND semantics and exact/keyword matching in examples/ticket_search/main.go"

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 and Phase 2
2. Implement User Story 1 tasks (T008–T012)
3. Validate search output and exit codes

### Incremental Delivery

1. Add User Story 2 pagination summaries and bounds handling
2. Add User Story 3 detail hydration
3. Complete polish tasks and run gofmt
