# Server Version Compatibility Client

A Go client that polls a server's `/info` endpoint and checks if the server version is compatible with a specified semver range.

## Features

- Polls the server every 10 seconds (configurable)
- Checks if the server version is compatible with a specified semver range
- Configurable via environment variables

## Installation

First, install the dependencies:

```bash
cd my-client
go get github.com/Masterminds/semver/v3
```

## Running the Client

```bash
cd my-client
go run .
```

## Configuration

The client can be configured using the following environment variables:

| Environment Variable | Description | Default Value |
|---------------------|-------------|---------------|
| `SERVER_URL` | The URL of the server's `/info` endpoint | `http://localhost:8080/info` |
| `SUPPORTED_VERSIONS` | The semver range of supported server versions | `>=0.1.0` |
| `POLL_INTERVAL_SECONDS` | The interval between polls in seconds | `10` |

### Example

```bash
# Set the supported versions to 0.1.0 through 0.2.0
export SUPPORTED_VERSIONS=">=0.1.0 <0.2.0"

# Run the client
go run .
```

## Semver Range Examples

- `>=0.1.0`: Version 0.1.0 or greater
- `>=0.1.0 <0.2.0`: Version 0.1.0 or greater, but less than 0.2.0
- `~0.1.0`: Approximately equivalent to version 0.1.0 (allows patch-level changes)
- `^0.1.0`: Compatible with version 0.1.0 (allows minor-level changes)

## Version Handling

The client handles server versions in the following way:

1. If the server returns a version with a 'v' prefix (e.g., 'v0.1.1'), the prefix is removed for semver compatibility
2. If the server returns a pre-release version (e.g., 'v0.1.1-dev'), the pre-release suffix is removed for compatibility checking 