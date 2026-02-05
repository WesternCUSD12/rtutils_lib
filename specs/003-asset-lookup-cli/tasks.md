# Implementation Tasks: Asset Lookup CLI

**Feature**: `003-asset-lookup-cli`
**Spec**: [spec.md](spec.md)
**Plan**: [plan.md](plan.md)

## Phase 1: Setup
- [x] T001 Create `examples/asset_lookup` directory and empty `main.go` file

## Phase 2: Foundational Logic
- [x] T002 Implement `loadEnv` function to parse `.env` file (if present) and set environment variables in `examples/asset_lookup/main.go`
- [x] T003 Define `main` function with `flag` parsing for `--id`, `--name`, `--internal-name` in `examples/asset_lookup/main.go`
- [x] T004 Implement validation: Ensure exactly one flag is provided and `RT_BASE_URL`/`RT_TOKEN` are set in `examples/asset_lookup/main.go`
- [x] T005 Initialize `rtutils_lib.Client` using the validated configuration in `examples/asset_lookup/main.go`

## Phase 3: User Story 1 - Lookup by ID
**Goal**: Retrieve a single asset by its explicit ID.
**Test**: `go run examples/asset_lookup/main.go --id 2091` -> Prints Details.

- [x] T006 [US1] Implement `printAssetDetails(asset *rtutils_lib.Asset)` using `text/tabwriter` for Vertical Key-Value output in `examples/asset_lookup/main.go`
- [x] T007 [US1] Implement `runIDSearch(client, id)`: Call `client.Assets.Get`, handle 404, call `printAssetDetails` in `examples/asset_lookup/main.go`
- [x] T008 [US1] Wire up `--id` flag to execute `runIDSearch` and exit with code 0 (success) or 1 (error) in `examples/asset_lookup/main.go`

## Phase 4: User Story 2 - Lookup by Name
**Goal**: Find asset(s) by Name, handling ambiguity.
**Test**: `go run ... --name "HOST3"` -> Details (if 1) or List (if >1).

- [X] T009 [US2] Implement `printAssetSummary(assets []rtutils_lib.Asset)` using `text/tabwriter` for Table output (ID, Name, URL) in `examples/asset_lookup/main.go`
- [X] T010 [US2] Implement `handleSearchResults(results)`: Logic to route 0 (Err), 1 (Details), or >1 (Summary) matches in `examples/asset_lookup/main.go`
- [X] T011 [US2] Wirup up `--name` flag: Construct AssetSQL `Name = '...'`, call `client.Assets.Search`, use handle logic in `examples/asset_lookup/main.go`

## Phase 5: User Story 3 - Lookup by Internal Name
**Goal**: Find asset(s) by Custom Field "Internal Name".
**Test**: `go run ... --internal-name "Fluffy"` -> Details/List.

- [X] T012 [US3] Wire up `--internal-name` flag: Construct AssetSQL `'CF.{Internal Name}' = '...'`, call `client.Assets.Search`, use reuse handle logic in `examples/asset_lookup/main.go`

## Phase 6: Polish
- [X] T013 Verify standard usage help text matches CLI conventions using `flag.Usage` in `examples/asset_lookup/main.go`
- [X] T014 Ensure correct Exit Codes (0=Success, 1=Error/Ambiguous, 2=NotFound) are returned in `examples/asset_lookup/main.go`

## Dependencies
- US2/US3 display logic (T009, T010) depends on US1 detail printer (T006)
- Setup tasks (T001-T005) block all user stories.
