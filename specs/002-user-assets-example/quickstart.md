# Quickstart: User Assets Example

## Prerequisites
*   Go 1.25+ installed
*   Access to a Request Tracker (RT) instance (v4.4+)
*   An API Token for RT

## Installation
Clone the repository:
```bash
git clone <repo-url>
cd rtutils_lib
```

## Configuration
Set the required environment variables:

```bash
export RT_BASE_URL="https://rt.example.com/REST/2.0"
export RT_TOKEN="your-api-token"
```

## Usage
Run the example program using `go run` or build it first.

### Flags
*   `--username`: Search by User Login Name
*   `--email`: Search by Email Address
*   `--name`: Search by Real Name

**Note**: Flags are mutually exclusive.

### Examples

Search by Username:
```bash
go run examples/user_assets/main.go --username jdoe
```

Search by Email:
```bash
go run examples/user_assets/main.go --email jdoe@example.com
```

Search by Real Name:
```bash
go run examples/user_assets/main.go --name "John Doe"
```

## Expected Output

On success (Table Format):
```text
ID      Name
123     MacBook Pro 2023
124     Dell Monitor 27"
```

On error (e.g., missing env vars):
```text
2026/02/04 12:00:00 RT_BASE_URL and RT_TOKEN environment variables are required
exit status 1
```
