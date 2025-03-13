# Simple Go HTTP Server

A minimalist HTTP server written in Go using only the standard library.

## Features

- `/info` endpoint that returns version information in JSON format
- Comprehensive test suite

## Running the Server

```bash
cd my-server
go run .
```

The server will start on port 8080 by default.

## Testing

Run the tests with:

```bash
cd my-server
go test -v
```

## API Endpoints

### GET /info

Returns version information in JSON format.

**Example Response:**

```json
{
  "version": "v0.1.1-dev"
}
``` 