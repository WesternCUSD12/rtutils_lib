# Feature Specification: Example Program - User Asset Lookup

**Feature Branch**: `002-user-assets-example`  
**Created**: 2026-02-04
**Status**: Draft  
**Input**: An example program should be created that pulls all assets assigned or held by a user. a new examples directory should be created.

## Clarifications

### Session 2026-02-04

- Q: Which relationship should be used to find user assets (HeldBy vs Owner)? → A: Option B - Support both (Search `Owner.Name` OR `HeldBy.Name`).
- Q: What format should the output take? → A: Option A - Human-Readable Table (using `text/tabwriter`).

## User Scenarios & Testing _(mandatory)_

### User Story 1 - Developer runs example program (Priority: P1)

**Why this priority**: It is the core request and provides immediate documentation-as-code usage for the library.

**Independent Test**: Can be tested by running `go run examples/user_assets/main.go --username [username]` and verifying it outputs the expected asset list or error message.

**Acceptance Scenarios**:

1.  **Scenario: Successful Lookup by Username**

    - **Given**: Valid `RT_BASE_URL` and `RT_TOKEN` environment variables are set.
    - **Given**: A user "jdoe" exists in RT and holds 2 assets.
    - **When**: I run the program with argument `--username jdoe`.
    - **Then**: The program prints the details (ID, Name) of the 2 assets held by "jdoe".

2.  **Scenario: Successful Lookup by Email**

    - **Given**: A user with email "jdoe@example.com" exists in RT.
    - **When**: I run the program with argument `--email jdoe@example.com`.
    - **Then**: The program prints the assets held by that user.

3.  **Scenario: Successful Lookup by Name**

    - **Given**: A user with Real Name "John Doe" exists in RT.
    - **When**: I run the program with argument `--name "John Doe"`.
    - **Then**: The program prints the assets held by that user.

4.  **Scenario: No Arguments Provided**

    - **When**: I run the program without arguments.
    - **Then**: The program prints usage instructions and exits with a non-zero status.

5.  **Scenario: Multiple Flags Provided**

    - **When**: I run the program with `--username jdoe --email jdoe@example.com`.
    - **Then**: The program prints an error stating that flags are mutually exclusive.
    - **Given**: `RT_TOKEN` is unset.
    - **When**: I run the program.
    - **Then**: The program prints an error about missing environment variables and exits.

6.  **Scenario: No Assets Found**
    _ **Given**: User "newhire" has no assets.
    _ **When**: I run the program with argument "newhire". \* **Then**: The program indicates that 0 assets were found.
    search criteria via flags: `--username`, `--email`, or `--name` (mutually exclusive).

- **FR-004**: System MUST use the `rtutils_lib` library to authenticate with the Request Tracker API.
- **FR-005**: System MUST execute an asset search query filtered by the provided specific criteria:
  - User: `Owner.Name = 'VAL' OR HeldBy.Name = 'VAL'`
  - Email: `Owner.EmailAddress = 'VAL' OR HeldBy.EmailAddress = 'VAL'`
  - Name: `Owner.RealName = 'VAL' OR HeldBy.RealName = 'VAL'`
- **FR-006**: System MUST output the ID and Name of each found asset to the console in a table format using `text/tabwriter`.
- **FR-007**: System MUST validate that required environment variables are present before attempting API calls.

### Key Entities

- **Asset**: The item being retrieved (ID, Name).
- **User**: The holder of the asset (Username, Email, or RealN
- **FR-001**: System MUST provide a main entry point at `examples/user_assets/main.go`.
- **FR-002**: System MUST read `RT_BASE_URL` and `RT_TOKEN` from environment variables.
- **FR-003**: System MUST accept a username as the first command-line argument.
- **FR-004**: System MUST use the `rtutils_lib` library to authenticate with the Request Tracker API.
- **FR-005**: System MUST execute an asset search query filtered by the provided username where `Owner.Name = 'USERNAME' OR HeldBy.Name = 'USERNAME'`.
- **FR-006**: System MUST output the ID and Name of each found asset to the console in a table format using `text/tabwriter`.
- **FR-007**: System MUST validate that required environment variables are present before attempting API calls.

### Key Entities

- **Asset**: The item being retrieved (ID, Name).
- **User**: The holder of the asset (Username).

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: `go run examples/user_assets/main.go` compiles successfully without errors.
- **SC-002**: When provided with valid inputs, the program outputs a list of assets within 5 seconds (assuming standard network latency).
- **SC-003**: The code serves as a valid copy-paste template for developers (clean, commented, idiomatic Go).
