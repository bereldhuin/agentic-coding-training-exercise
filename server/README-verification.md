# Cross-Server API Verification

This directory contains a verification script that tests API compatibility across all server implementations (TypeScript, Python, Swift, Kotlin, Go).

## Purpose

The verification script ensures that all server implementations behave identically by:
- Testing all endpoints defined in the canonical OpenAPI specification
- Comparing responses across different server implementations
- Detecting behavioral differences early
- Validating that the OpenAPI specification is correctly implemented

## Prerequisites

Before running the verification script, ensure you have:

1. **Node.js 18+** installed
   ```bash
   node --version  # Should be v18 or higher
   ```

2. **PM2** installed globally
   ```bash
   npm install -g pm2
   ```

3. **Server builds** available
   - TypeScript: `npm run build` in `server/typescript`
   - Python: No build required
   - Swift: `swift build` in `server/swift`
   - Kotlin: `./gradlew installDist` in `server/kotlin`
   - Go: `go build` in `server/go`

## Installation

The verification script uses Node.js and requires PM2 for process management. Dependencies are already included in `package.json`.

To install dependencies:
```bash
npm install
```

## Usage

### Server Runner CLI

Start a specific server or verify all servers sequentially:
```bash
node scripts/server-runner.js
```

The CLI prompts you to select a server (TypeScript, Go, Python, Swift, Kotlin) or run a verify-all flow that starts each server, checks `/health`, and shuts it down.

### Basic Usage

Test all configured servers:
```bash
npm run verify
```

Or run the verification CLI directly:
```bash
node scripts/server-runner.js --verify-api
```

### Command-Line Options

```bash
node scripts/server-runner.js --verify-api [options]
```

Available options:

| Option | Description | Example |
|--------|-------------|---------|
| `--servers <list>` | Comma-separated list of servers to test | `--servers typescript,python` |
| `--tests <categories>` | Comma-separated list of test categories | `--tests health,crud` |
| `--verbose` | Enable detailed output including request/response details | `--verbose` |
| `--timeout <seconds>` | Server startup timeout in seconds (default: 30) | `--timeout 60` |
| `--output <file>` | Write JSON report to file | `--output report.json` |
| `--tap` | Output in Test Anything Protocol (TAP) format | `--tap` |
| `--help` | Show help message | `--help` |

### Examples

**Test specific servers:**
```bash
npm run verify -- --servers typescript,python
```

**Test specific categories:**
```bash
npm run verify -- --tests health,crud
```

**Verbose output with JSON report:**
```bash
npm run verify -- --verbose --output report.json
```

**TAP format for CI/CD:**
```bash
npm run verify -- --tap
```

**Custom timeout for slow-starting servers:**
```bash
npm run verify -- --timeout 60
```

## Server Configuration

Servers are configured in `scripts/lib/servers.js` with the following properties:

| Property | Description | Example |
|----------|-------------|---------|
| `name` | Server identifier | `'typescript'` |
| `startCommand` | PM2 start command | `'./server/typescript/dist/infrastructure/http/server.js'` |
| `port` | HTTP port number | `3000` |
| `healthPath` | Health check endpoint | `'/health'` |
| `startupTimeout` | Max seconds to wait for server ready | `30` |

### Configured Servers

| Server | Port | Start Command |
|--------|------|---------------|
| TypeScript | 3000 | `./server/typescript/dist/infrastructure/http/server.js` |
| Python | 8000 | `uvicorn server.python.main:app` |
| Swift | 9000 | `./server/swift/.build/release/Server` |
| Kotlin | 8080 | `./server/kotlin/build/install/server/bin/server` |
| Go | 8081 | `./server/go/bin/server` |

### Adding a New Server

To add a new server implementation, edit `scripts/lib/servers.js` and add to the `servers` array:

```javascript
{
  name: 'rust',
  startCommand: './server/rust/target/release/server',
  port: 8082,
  healthPath: '/health',
  startupTimeout: 30,
}
```

## Test Categories

The verification script includes the following test categories:

| Category | Description | Tests |
|----------|-------------|-------|
| `health` | Health check endpoints | GET /health |
| `list` | List items endpoints | GET /v1/items with various filters |
| `crud` | CRUD operations | GET, POST, PUT, PATCH, DELETE /v1/items |
| `validation` | Validation error tests | Invalid data, missing fields, enum violations |

### Test Cases

Each category includes multiple test cases:

**Health:**
- GET /health returns 200 with status and timestamp

**List:**
- GET /v1/items returns empty list
- GET /v1/items with status filter
- GET /v1/items with limit parameter

**CRUD:**
- GET /v1/items/:id returns 404 for non-existent item

**Validation:**
- POST /v1/items with title too short (< 3 chars) returns 400
- POST /v1/items with negative price returns 400
- POST /v1/items with invalid condition enum returns 400

## Output

### Console Output

The script provides color-coded console output:

- Green (✓): Test passed
- Red (✗): Test failed
- Blue: Informational messages
- Yellow: Warnings

Example output:
```
Cross-Server API Verification
============================================================

Generated 12 tests

typescript Server (3000)
==================================================
  Starting typescript server...
  ✓ typescript server started
  Waiting for typescript to be ready...
  ✓ typescript is ready
  ✓ GET /health returns 200
  ✓ GET /v1/items returns empty list
  ...

python Server (8000)
==================================================
  Starting python server...
  ✓ python server started
  ...
```

### Summary

After all tests complete, a summary is displayed:

```
============================================================
Summary
============================================================
Total tests: 60
Passed: 58
Failed: 2

Results by server:
  ✓ typescript: 12/12 passed
  ✗ python: 10/12 passed
  ✓ swift: 12/12 passed
  ✗ kotlin: 11/12 passed
  ✓ go: 12/12 passed

Failed tests:
  - python: POST /v1/items with negative price returns 400
    status code mismatch: expected 400 but got 201
  - kotlin: GET /v1/items returns empty list
    schema validation failed: missing required key 'items'
```

### JSON Report

When `--output report.json` is specified, a detailed JSON report is generated:

```json
{
  "timestamp": "2026-01-31T12:00:00.000Z",
  "summary": {
    "total": 60,
    "passed": 58,
    "failed": 2
  },
  "byServer": {
    "typescript": {
      "total": 12,
      "passed": 12,
      "failed": 0,
      "tests": [...]
    },
    ...
  },
  "failures": [...]
}
```

### TAP Format

When `--tap` is specified, output follows the Test Anything Protocol:

```
TAP version 13
# Testing servers: typescript,python
# Test categories: health,crud
ok 1 - typescript: GET /health returns 200
ok 2 - typescript: GET /v1/items returns empty list
...
1..60
# 58 passed, 2 failed
```

## Response Comparison

The verification script compares responses across servers with intelligence:

### Exact Comparison
- HTTP status codes must match exactly
- Response structure (keys, nesting) must match
- Data values must match (except for special fields)

### Tolerated Differences

**Timestamps:**
- Validated as ISO 8601 format
- Exact values not compared (servers generate timestamps at different times)
- Fields: `timestamp`, `created_at`, `updated_at`, `published_at`

**IDs:**
- Validated as correct type (number/string)
- Exact values not compared (servers may generate different IDs)
- Fields: `id`

**Extra Fields:**
- Extra fields in actual response are logged as warnings (not errors)
- Helps detect undocumented response fields

## How It Works

1. **Parse Arguments**: Command-line options are parsed to configure which servers and tests to run.

2. **Check Dependencies**: Verifies PM2 is installed and available.

3. **Generate Tests**: Creates test cases from the OpenAPI specification (or hardcoded tests).

4. **Sequential Testing**: For each server:
   - Start server with PM2 (`--no-autorestart` to prevent infinite restart loops)
   - Poll health endpoint until server is ready (or timeout)
   - Run all test cases
   - Stop/delete the PM2 process
   - Move to next server

5. **Response Comparison**: Compare responses against expected status codes and schemas.

6. **Reporting**: Print summary, detailed failures, and optionally write JSON report.

7. **Cleanup**: Ensure all PM2 processes are stopped (signal handlers for graceful shutdown).

## Troubleshooting

### PM2 Not Found

**Error:** `Error: PM2 is not installed`

**Solution:** Install PM2 globally:
```bash
npm install -g pm2
```

### Server Won't Start

**Error:** `Failed to start {server} server`

**Possible causes:**
- Build artifacts don't exist (run build command)
- Start command path is incorrect
- Port is already in use

**Solution:**
```bash
# Check if port is in use
lsof -i :3000  # Replace with your port

# Kill process using the port
kill -9 <PID>

# Rebuild server
cd server/typescript && npm run build
```

### Port Conflicts

**Error:** Server starts but health check fails

**Solution:** Ensure each server is configured to use a unique port:
- TypeScript: 3000
- Python: 8000
- Swift: 9000
- Kotlin: 8080
- Go: 8081

### Server Health Check Timeout

**Error:** `{server} failed to start within 30s`

**Solution:** Increase timeout:
```bash
npm run verify -- --timeout 60
```

Or add `startupTimeout: 60` to server config in `scripts/lib/servers.js`.

### Database Issues

**Error:** Tests fail with database errors

**Solution:** Initialize and seed database:
```bash
cd server/typescript
npm run db:init
npm run db:seed  # Optional
```

### OpenAPI Spec Not Found

**Warning:** `Warning: OpenAPI spec not found at server/api/openapi.yaml`

**Solution:** This is a warning, not an error. The script uses hardcoded tests. To use OpenAPI-driven tests, ensure the spec exists at the expected path.

### Cleanup Issues

**Problem:** PM2 processes left running after script exits

**Solution:** Manually clean up:
```bash
pm2 delete all
```

Or use the cleanup function in the script (automatically triggered on SIGINT/SIGTERM).

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Verify API Compatibility

on: [push, pull_request]

jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install PM2
        run: npm install -g pm2
      - name: Install dependencies
        run: |
          cd server/typescript && npm install
      - name: Build servers
        run: |
          cd server/typescript && npm run build
          cd ../swift && swift build
          cd ../kotlin && ./gradlew installDist
          cd ../go && go build
      - name: Initialize database
        run: cd server/typescript && npm run db:init
      - name: Run verification
        run: cd server/typescript && npm run verify -- --tap
      - name: Upload report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: verification-report
          path: server/report.json
```

### Exit Codes

The script exits with:
- `0`: All tests passed
- `1`: One or more tests failed

This allows CI/CD pipelines to detect failures:

```bash
npm run verify
if [ $? -eq 0 ]; then
  echo "All tests passed"
else
  echo "Some tests failed"
  exit 1
fi
```

## Advanced Usage

### Filtering Tests

Run only specific test categories:
```bash
npm run verify -- --tests health,validation
```

### Verbose Debugging

Enable verbose output to see request/response details:
```bash
npm run verify -- --verbose
```

### Combining Options

Multiple options can be combined:
```bash
npm run verify -- --servers typescript,python --tests crud,validation --verbose --output report.json
```

## Contributing

When adding new server implementations:

1. Build the server executable
2. Add server configuration to `scripts/lib/servers.js`
3. Assign a unique port
4. Test the verification script
5. Update this documentation with any special requirements

## License

MIT
