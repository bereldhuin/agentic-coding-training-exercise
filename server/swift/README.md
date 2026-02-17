# LeBonPoint Swift Server

A Swift Vapor implementation of the LeBonPoint Marketplace API. This server implements the canonical OpenAPI specification defined in `server/api/openapi.yaml` using the hexagonal architecture pattern.

## Prerequisites

- Swift 5.9+ (for async/await support)
- macOS 13+ or Linux
- SQLite3 (for database)

## Installation

1. Navigate to the server directory:
```bash
cd server/swift
```

2. Resolve dependencies:
```bash
swift package resolve
```

3. Build the project:
```bash
swift build
```

## Configuration

Create a `.env` file in the `server/swift` directory (copy from `.env.example`):

```env
PORT=9000
DATABASE_PATH=../database/db.sqlite
ENVIRONMENT=development
LOG_LEVEL=info
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `9000` |
| `DATABASE_PATH` | Path to SQLite database | `../database/db.sqlite` |
| `ENVIRONMENT` | Environment (development, testing, production) | `development` |
| `LOG_LEVEL` | Logging level (debug, info, warning, error, critical) | `info` |

## Running the Server

### Development Mode

```bash
swift run
```

### Production Mode

1. Build the release binary:
```bash
swift build -c release
```

2. Run the binary:
```bash
.build/release/App
```

## Testing

Run all tests:
```bash
swift test
```

Run tests with verbose output:
```bash
swift test --verbose
```

## Docker Support

### Build Docker Image

```bash
docker build -t lebonpoint-swift-server .
```

### Run with Docker

```bash
docker run -p 9000:9000 \
  -v $(pwd)/../database:/app/database \
  -e PORT=9000 \
  -e DATABASE_PATH=/app/database/db.sqlite \
  lebonpoint-swift-server
```

### Docker Compose

See the root `docker-compose.yml` for running all servers together.
