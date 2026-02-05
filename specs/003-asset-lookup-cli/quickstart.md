# Quickstart: Asset Lookup CLI

**Branch**: `003-asset-lookup-cli`

## Prerequisites

- Go 1.25+ installed
- Access to a Request Tracker instance (RT 4.4/5.0 with REST 2.0)
- `RT_BASE_URL` and `RT_TOKEN` available (env vars or `.env` file)

## Installation

This is an example tool included in the library. No separate installation is required if you have the source.

## Configuration

Create a `.env` file in the project root or the `examples/asset_lookup/` directory:

```bash
RT_BASE_URL=https://rt.example.com/REST/2.0
RT_TOKEN=your-api-token
```

## Usage

Run the tool using `go run`:

```bash
# 1. Lookup by ID
go run examples/asset_lookup/main.go --id 2091

# 2. Lookup by Name
go run examples/asset_lookup/main.go --name "HOST3"

# 3. Lookup by Internal Name
go run examples/asset_lookup/main.go --internal-name "Fluffy Rooster"
```

## Expected Output

**Success (Single Match):**
```text
ID:             2091
Name:           HOST3
Status:         production
Type:           Server
Model:          PowerEdge R740
Manufacturer:   Dell
Serial Number:  ABC1234
Internal Name:  Fluffy Rooster
```

**Ambiguous (Multiple Matches):**
```text
Multiple assets found. Please refine your search.

ID    Name       URL
2091  HOST3      ...
2092  HOST3-OLD  ...
```

**Failure:**
```text
Asset not found.
```
