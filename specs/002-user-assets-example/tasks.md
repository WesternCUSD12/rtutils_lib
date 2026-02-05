# Implementation Tasks: Example Program - User Asset Lookup

**Feature**: `002-user-assets-example`
**Spec**: [spec.md](spec.md)
**Plan**: [plan.md](plan.md)

## Phase 1: Setup
- [x] T001 Create `examples/user_assets` directory and empty `main.go` file

## Phase 2: Foundational Logic
- [x] T002 Implement `main` function structure with Environment Variable validation (`RT_BASE_URL`, `RT_TOKEN`) in `examples/user_assets/main.go`
- [x] T003 Initialize `rtutils_lib.Client` with validated environment variables in `examples/user_assets/main.go`

## Phase 3: User Story 1 - Developer runs example program
**Goal**: Developer can run the tool with flags to query assets.
**Test**: `go run examples/user_assets/main.go --username jdoe` prints table.

- [x] T004 [US1] Define CLI flags (`--username`, `--email`, `--name`) using `flag` package in `examples/user_assets/main.go`
- [x] T005 [US1] Implement validation logic to ensure exactly one flag is provided (mutually exclusive) in `examples/user_assets/main.go`
- [x] T006 [US1] Implement AssetSQL query builder: maps flag to `Owner` OR `HeldBy` logic in `examples/user_assets/main.go`
- [x] T007 [US1] Call `client.Assets.Search` with the constructed query in `examples/user_assets/main.go`
- [x] T008 [US1] Format and print the search results (ID, Name) using `text/tabwriter` in `examples/user_assets/main.go`
- [x] T009 [US1] Handle error cases (search failure, no results found) with user-friendly messages in `examples/user_assets/main.go`

## Phase 4: Polish
- [x] T010 Verify `go run examples/user_assets/main.go -h` prints helpful usage text
- [x] T011 Verify code is formatted with `gofmt`

## Dependencies
- All US1 tasks depend on T001-T003.
- T007 depends on T006.
