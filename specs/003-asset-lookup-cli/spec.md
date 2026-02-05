# Feature Specification: Example Tool - Asset Detail Lookup

**Feature Branch**: `003-asset-lookup-cli`
**Created**: 2026-02-05
**Status**: Draft
**Input**: User description: "An example cli tool should be created that pulls an assets detail by providing the RT ID, the name (ie, W12-1234), or the Internal Name (ie, Fluffy Rooster)"

## User Scenarios & Testing

### ## Clarifications

### Session 2026-02-05

- Q: How to handle multiple matches? → A: Print summary table and exit (no interactive selection).
- Q: Output format for single asset details? → A: Vertical Key-Value list.
- Q: Auth loading method? → A: Environment variables (dot files supported).

### User Story 1 - Lookup Asset by Request Tracker ID (Priority: P1)

A developer using the example tool can retrieve full details of a specific asset by providing its unique numeric ID. This is the most direct and precise way to find a known asset.

**Why this priority**: Essential functionality for debugging and direct object retrieval.

**Independent Test**: Can be verified by running the tool with a known ID and checking if the correct asset data is displayed.

**Acceptance Scenarios**:

1. **Given** a valid, existing Asset ID (e.g., 2091), **When** the user runs `go run examples/asset_lookup/main.go --id 2091`, **Then** the tool outputs the full details of that asset (Name, Status, Custom Fields).
2. **Given** a non-existent Asset ID, **When** the user runs the tool, **Then** it displays a clear "Asset not found" error message.

---

### User Story 2 - Lookup Asset by Asset Name (Priority: P1)

A developer can find an asset by its primary Name (e.g., "W12-1234"). The Name is the main human-readable identifier in the RT system.

**Why this priority**: Users often know the barcode or label name rather than the internal DB ID.

**Independent Test**: Run with a known asset name and verify output.

**Acceptance Scenarios**:

1. **Given** an existing asset named "HOST3", **When** the user runs `go run examples/asset_lookup/main.go --name HOST3`, **Then** the tool displays the details for that asset.
2. **Given** multiple assets potentially share a name (if allowed), **When** the tool finds multiple matches, **Then** it lists the summaries of all matching assets (ID, Name) and exits without showing details.

---

### User Story 3 - Lookup Asset by Internal Name (Priority: P1)

A developer can find an asset by its "Internal Name" (e.g., "Fluffy Rooster"), which is stored as a Custom Field in RT.

**Why this priority**: Organizations often use pet names or project names stored in custom fields.

**Independent Test**: Run with a known internal name and verify output.

**Acceptance Scenarios**:

1. **Given** an asset with Custom Field "Internal Name" set to "Fluffy Rooster", **When** the user runs `go run examples/asset_lookup/main.go --internal-name "Fluffy Rooster"`, **Then** the tool displays the details for that asset.

---

### Edge Cases

- **Missing Credentials**: If `RT_BASE_URL` or `RT_TOKEN` are missing, the tool should exit with a helpful error.
- **No Flags**: If run without arguments, print usage instructions.
- **Multiple Flags**: If multiple search flags are provided (e.g., `--id` and `--name`), the tool should either prioritize one or error out (Decision: Error, mutually exclusive).

## Requirements

### Functional Requirements

- **FR-001**: The tool MUST accept configuration via environment variables `RT_BASE_URL` and `RT_TOKEN`. It SHOULD support loading these from a `.env` file if present (using `godotenv` or similar).
- **FR-002**: The tool MUST support the following mutually exclusive CLI flags:
  - `--id`: Search by numeric Asset ID.
  - `--name`: Search by the core `Name` field.
  - `--internal-name`: Search by the custom field named "Internal Name".
- **FR-003**: When searching by `--id`, the system MUST use the direct API `GET` method.
- **FR-004**: When searching by `--name` or `--internal-name`, the system MUST use the `Search` (AssetSQL) method.
  - For `--internal-name`, the query format MUST be `'CF.{Internal Name}' = 'VALUE'`.
- **FR-005**: The tool MUST display a formatted list of asset properties in a vertical Key-Value format (e.g., using `tabwriter` with aligned columns or simple `fmt.Printf`).
  - ID
  - Name
  - Status
  - Type (Custom Field)
  - Model (Custom Field)
  - Manufacturer (Custom Field)
  - All other non-empty Custom Fields available on the asset.
- **FR-006**: If no assets are found, the tool MUST print a user-friendly "No results found" message to `stderr` and exit with a non-zero code.
- **FR-007**: If exactly one asset is found, output its full details.
- **FR-008**: If multiple assets are found (via name search), list a summary table (ID, Name, URL) to stdout and exit with a non-zero exit code to indicate the specific target was ambiguous.

### Success Criteria

- **Efficiency**: Returns asset details in under 2 seconds for typical network conditions.
- **Usability**: Users can retrieve an asset knowing only one of its identifiers (ID, Name, or Internal Name).
- **Clarity**: Output is human-readable (not raw JSON), utilizing a vertical key-value format for detail views.

### Assumptions

- The `rtutils_lib` library is available and functions as expected.
- The "Internal Name" is indeed a Custom Field on the assets in the target RT instance.
- The environment has the necessary permissions to read assets.
