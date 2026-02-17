# LeBonPoint Kotlin Server

Kotlin implementation of the LeBonPoint Marketplace API using Ktor framework.

## Prerequisites

- Kotlin 1.9+
- JDK 17+
- Gradle 8+

## Installation

```bash
# Install dependencies
./gradlew build
```

## Running the Server

```bash
# Run with default configuration (port 8080)
./gradlew run

# Run with custom port
PORT=9000 ./gradlew run

# Run with custom database path
DATABASE_PATH=/path/to/db.sqlite ./gradlew run
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_PATH` - Path to SQLite database (default: ../database/db.sqlite)
- `ENV` - Environment: development, testing, production (default: development)

## Running Tests

```bash
# Run all tests
./gradlew test

# Run with coverage
./gradlew test jacocoTestReport
```
