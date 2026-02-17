# LeBonPoint Python Server (FastAPI)

Python FastAPI implementation of the LeBonPoint marketplace items API.

## Requirements

- Python 3.11+
- SQLite database (shared with TypeScript implementation)

## Installation

1. **Install Poetry (if not already available):**

```bash
python -m pip install --user poetry
```

2. **Install project dependencies:**

```bash
cd server/python
poetry install --with dev
```

Poetry will create a reproducible lockfile (`poetry.lock`). You can activate the virtualenv with `poetry shell` or run commands via `poetry run <cmd>`.

## Configuration

The server can be configured via environment variables or a `.env` file:

```bash
# Server configuration
PORT=8000

# Database path (relative to server/python directory)
DATABASE_PATH=../database/db.sqlite

# Environment
ENVIRONMENT=development
```

## Running the Server

### Development mode (with auto-reload):

```bash
poetry run uvicorn src.main:app --reload --port 8000
```

### Production mode:

```bash
poetry run uvicorn src.main:app --port 8000 --workers 4
```

Or using the built-in main:

```bash
poetry run python -m src.main
```

## Testing

### Run all tests:

```bash
pytest
```

### Run with coverage:

```bash
pytest --cov=src --cov-report=html
```

### Run specific test file:

```bash
pytest tests/unit/domain/test_enums.py
```

### Run with verbose output:

```bash
pytest -v
```
