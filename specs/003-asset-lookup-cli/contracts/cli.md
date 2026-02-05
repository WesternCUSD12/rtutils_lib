# CLI Interface Contract

**Tool**: `examples/asset_lookup`

## Arguments

The tool accepts **exactly one** of the following search flags.

| Flag | Type | Description | Requirement |
|------|------|-------------|-------------|
| `--id` | `string` | The numeric ID of the asset. | Mutually exclusive |
| `--name` | `string` | The asset `Name`. | Mutually exclusive |
| `--internal-name` | `string` | The `Internal Name` custom field. | Mutually exclusive |

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `RT_BASE_URL` | Base URL of the RT API (e.g. `https://rt.example.com/REST/2.0`) | Yes |
| `RT_TOKEN` | API Token for authentication | Yes |

## Exit Codes

| Code | Description |
|------|-------------|
| `0` | Success (found 1 asset, or found >1 and printed summary) |
| `1` | Configuration error (missing flags/env) or search failed |
| `2` | No assets found |
