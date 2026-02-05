# Data Model: User Asset Lookup Example

**Feature**: `002-user-assets-example`

## Entities

### Asset (Read-Only)
This entity is retrieved from the `rtutils_lib`.

| Field | Type | Description | Source |
| :--- | :--- | :--- | :--- |
| `id` | `string` | Unique numeric identifier | `Asset.ID` |
| `Name` | `string` | Display name of the asset | `Asset.Name` |
| `HeldBy/Owner` | `Relationship` | The user holding/owning the asset | Query Parameter |

### User (Input)
This entity represents the target of the search.

| Field | Type | Description | Source |
| :--- | :--- | :--- | :--- |
| `Username` | `string` | Login name of the user | CLI Flag `--username` |
| `Email` | `string` | Email address of the user | CLI Flag `--email` |
| `Name` | `string` | Real Name of the user | CLI Flag `--name` |

## API Contracts (Usage)

The example program acts as a **Client** to the `rtutils_lib`.

### `AssetService.Search`
*   **Input**: `ctx context.Context`, `query string`
*   **Query Construction**:
    *   `Username`: `Owner.Name = 'val' OR HeldBy.Name = 'val'`
    *   `Email`: `Owner.EmailAddress = 'val' OR HeldBy.EmailAddress = 'val'`
    *   `Name`: `Owner.RealName = 'val' OR HeldBy.RealName = 'val'`
*   **Output**: `*SearchResult[Asset]`, `error`

## Data Flow

1.  **Input**: User launches program with flags (e.g., `-u jdoe`).
2.  **Validation**: Verify exactly one search flag is provided.
3.  **Auth**: User initializes `Client` with Env Vars.
4.  **Request**: `Client` sends `GET /assets?query=...` to RT API.
5.  **Response**: RT API returns JSON list of Assets.
6.  **Output**: Program formats list as an aligned table using `text/tabwriter`.
