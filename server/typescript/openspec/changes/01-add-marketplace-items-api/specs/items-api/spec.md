## ADDED Requirements

### Requirement: Item Entity Model

The system SHALL define an Item entity with the following fields:
- `id`: Auto-incrementing integer primary key
- `title`: Required string, 3-200 characters
- `description`: Optional string
- `price_cents`: Required non-negative integer
- `category`: Optional string (e.g., "Ordinateurs")
- `condition`: Required enum value from `new|like_new|good|fair|parts|unknown`
- `status`: Required enum value from `draft|active|reserved|sold|archived`
- `is_featured`: Required boolean
- `city`: Optional string
- `postal_code`: Optional string (preserve leading zeros)
- `country`: Required string, default "FR"
- `delivery_available`: Required boolean
- `created_at`: Required ISO 8601 timestamp
- `updated_at`: Required ISO 8601 timestamp
- `published_at`: Optional nullable ISO 8601 timestamp
- `images`: JSON array of ItemImage objects with `url` (required), `alt` (optional), `sort_order` (optional)

#### Scenario: Valid item creation
- **GIVEN** all required fields are provided with valid values
- **WHEN** creating an item
- **THEN** the item is created with the provided values

#### Scenario: Default value application
- **GIVEN** an item is created without optional fields
- **WHEN** the item is persisted
- **THEN** default values are applied (country="FR", images=[])

### Requirement: List Items with Filtering

The system SHALL provide a `GET /v1/items` endpoint that returns a paginated list of items with optional filtering.

#### Scenario: List all items with default sorting
- **GIVEN** multiple items exist in the database
- **WHEN** requesting `GET /v1/items` without query parameters
- **THEN** return items sorted by `created_at:desc`, limited to 20 items

#### Scenario: Filter by status
- **GIVEN** items with various statuses exist
- **WHEN** requesting `GET /v1/items?status=active`
- **THEN** return only items with status="active"

#### Scenario: Filter by category
- **GIVEN** items from various categories exist
- **WHEN** requesting `GET /v1/items?category=Ordinateurs`
- **THEN** return only items where category="Ordinateurs"

#### Scenario: Filter by price range
- **GIVEN** items with various prices exist
- **WHEN** requesting `GET /v1/items?min_price_cents=10000&max_price_cents=50000`
- **THEN** return only items with price_cents between 10000 and 50000

#### Scenario: Filter by location
- **GIVEN** items from various cities exist
- **WHEN** requesting `GET /v1/items?city=Strasbourg&postal_code=67000`
- **THEN** return only items matching both city and postal_code

#### Scenario: Filter by boolean flags
- **GIVEN** items with various is_featured and delivery_available values exist
- **WHEN** requesting `GET /v1/items?is_featured=true&delivery_available=false`
- **THEN** return only items where is_featured=true AND delivery_available=false

#### Scenario: Custom sorting
- **GIVEN** multiple items exist
- **WHEN** requesting `GET /v1/items?sort=price_cents:asc`
- **THEN** return items sorted by price_cents ascending

#### Scenario: Pagination with limit
- **GIVEN** more than 50 items exist
- **WHEN** requesting `GET /v1/items?limit=50`
- **THEN** return at most 50 items

#### Scenario: Pagination validation
- **GIVEN** items exist
- **WHEN** requesting `GET /v1/items?limit=200`
- **THEN** return validation error as limit exceeds maximum of 100

### Requirement: Cursor-based Pagination

The system SHALL support cursor-based pagination using an opaque `cursor` parameter and return `next_cursor` in the response.

#### Scenario: Initial page request
- **GIVEN** 100 items exist in the database
- **WHEN** requesting `GET /v1/items?limit=20`
- **THEN** return 20 items with `next_cursor` in the response metadata

#### Scenario: Next page using cursor
- **GIVEN** a previous request returned `next_cursor=eyJpZCI6MjB9`
- **WHEN** requesting `GET /v1/items?limit=20&cursor=eyJpZCI6MjB9`
- **THEN** return the next 20 items after id=20

#### Scenario: Last page returns no cursor
- **GIVEN** all items have been retrieved
- **WHEN** requesting the final page
- **THEN** return remaining items without `next_cursor` in response

### Requirement: Full-Text Search

The system SHALL provide full-text search capability using SQLite FTS5 with accent-insensitive matching.

#### Scenario: Simple term search
- **GIVEN** items with titles "Ordinateur portable HP" and "Appareil photo Canon" exist
- **WHEN** requesting `GET /v1/items?q=ordinateur`
- **THEN** return items matching "ordinateur" in title or description

#### Scenario: Multiple terms search
- **GIVEN** items exist with various titles
- **WHEN** requesting `GET /v1/items?q=hp+elitebook`
- **THEN** return items containing both "hp" AND "elitebook"

#### Scenario: Prefix search
- **GIVEN** items exist
- **WHEN** requesting `GET /v1/items?q=ordi*`
- **THEN** return items with words starting with "ordi" (e.g., "ordinateur", "ordin")

#### Scenario: Phrase search
- **GIVEN** items exist
- **WHEN** requesting `GET /v1/items?q="ordinateur+portable"`
- **THEN** return items with the exact phrase "ordinateur portable"

#### Scenario: FTS results ranked by relevance
- **GIVEN** items with varying match quality exist
- **WHEN** performing a search query
- **THEN** return results ordered by BM25 relevance score, then by created_at desc

#### Scenario: Combined FTS with filters
- **GIVEN** items exist across various categories
- **WHEN** requesting `GET /v1/items?q=ordinateur&category=Ordinateurs&status=active`
- **THEN** return only active items in Ordinateurs category matching "ordinateur"

### Requirement: Create Item

The system SHALL provide a `POST /v1/items` endpoint to create a new item.

#### Scenario: Successful item creation
- **GIVEN** the user provides valid item data with title and price_cents
- **WHEN** sending `POST /v1/items` with valid JSON body
- **THEN** return `201 Created` with the created item including generated id and timestamps

#### Scenario: Validation error - title too short
- **GIVEN** the user provides a title with 2 characters
- **WHEN** sending `POST /v1/items` with title="HP"
- **THEN** return `400 Bad Request` with validation error about minimum title length

#### Scenario: Validation error - invalid condition
- **GIVEN** the user provides condition="unknown_value"
- **WHEN** sending `POST /v1/items` with invalid condition
- **THEN** return `400 Bad Request` with validation error about allowed condition values

#### Scenario: Validation error - negative price
- **GIVEN** the user provides price_cents=-100
- **WHEN** sending `POST /v1/items` with negative price
- **THEN** return `400 Bad Request` with validation error about non-negative price

#### Scenario: Default values applied
- **GIVEN** the user omits optional fields in the request
- **WHEN** sending `POST /v1/items` with only title and price_cents
- **THEN** create item with defaults: status="draft", country="FR", images=[], is_featured=false

### Requirement: Get Single Item

The system SHALL provide a `GET /v1/items/{id}` endpoint to retrieve a single item by ID.

#### Scenario: Existing item
- **GIVEN** an item with id=123 exists
- **WHEN** requesting `GET /v1/items/123`
- **THEN** return `200 OK` with the complete item data

#### Scenario: Non-existent item
- **GIVEN** no item with id=999 exists
- **WHEN** requesting `GET /v1/items/999`
- **THEN** return `404 Not Found` with error details

### Requirement: Replace Item (PUT)

The system SHALL provide a `PUT /v1/items/{id}` endpoint to completely replace an existing item.

#### Scenario: Successful replacement
- **GIVEN** an item with id=123 exists
- **WHEN** sending `PUT /v1/items/123` with complete item data including all required fields
- **THEN** return `200 OK` with the updated item and updated_at timestamp

#### Scenario: Replace non-existent item
- **GIVEN** no item with id=999 exists
- **WHEN** sending `PUT /v1/items/999` with valid data
- **THEN** return `404 Not Found`

#### Scenario: Missing required field on PUT
- **GIVEN** an item exists
- **WHEN** sending `PUT /v1/items/{id}` without required field `title`
- **THEN** return `400 Bad Request` with validation error

### Requirement: Patch Item

The system SHALL provide a `PATCH /v1/items/{id}` endpoint to partially update an existing item.

#### Scenario: Partial update
- **GIVEN** an item with id=123 exists with title="Old Title"
- **WHEN** sending `PATCH /v1/items/123` with `{ "title": "New Title" }`
- **THEN** return `200 OK` with updated title and updated_at timestamp, other fields unchanged

#### Scenario: Patch non-existent item
- **GIVEN** no item with id=999 exists
- **WHEN** sending `PATCH /v1/items/999` with valid patch data
- **THEN** return `404 Not Found`

#### Scenario: Patch with invalid value
- **GIVEN** an item exists
- **WHEN** sending `PATCH /v1/items/{id}` with `{ "price_cents": -10 }`
- **THEN** return `400 Bad Request` with validation error

### Requirement: Delete Item

The system SHALL provide a `DELETE /v1/items/{id}` endpoint to delete an item.

#### Scenario: Successful deletion
- **GIVEN** an item with id=123 exists
- **WHEN** sending `DELETE /v1/items/123`
- **THEN** return `204 No Content` and the item is removed from database

#### Scenario: Delete non-existent item
- **GIVEN** no item with id=999 exists
- **WHEN** sending `DELETE /v1/items/999`
- **THEN** return `404 Not Found`

### Requirement: Standard Error Format

All error responses SHALL follow a consistent JSON format.

#### Scenario: Validation error response
- **GIVEN** a request contains invalid data
- **WHEN** the validation fails
- **THEN** return `400 Bad Request` with `{ "error": { "code": "validation_error", "message": "...", "details": {} } }`

#### Scenario: Not found error response
- **GIVEN** a requested resource does not exist
- **WHEN** attempting to access it
- **THEN** return `404 Not Found` with `{ "error": { "code": "not_found", "message": "..." } }`

#### Scenario: Internal server error response
- **GIVEN** an unexpected error occurs
- **WHEN** processing the request
- **THEN** return `500 Internal Server Error` with `{ "error": { "code": "internal_error", "message": "..." } }`

### Requirement: Hexagonal Architecture

The system SHALL be implemented using hexagonal architecture (ports and adapters pattern).

#### Scenario: Domain layer independence
- **GIVEN** the domain layer contains entities and repository interfaces
- **WHEN** examining the domain code
- **THEN** no external dependencies (Express, SQLite) are imported in domain files

#### Scenario: Repository port interface
- **GIVEN** the application needs item persistence
- **WHEN** defining the repository port
- **THEN** create an interface with methods: create, findById, findAll, update, delete, search

#### Scenario: SQLite repository adapter
- **GIVEN** the repository port interface exists
- **WHEN** implementing the SQLite adapter
- **THEN** the adapter implements the port interface using better-sqlite3

#### Scenario: Use case layer orchestrates
- **GIVEN** a use case like CreateItem
- **WHEN** executing the use case
- **THEN** it coordinates between validation and repository without knowing HTTP details

### Requirement: Test Coverage

The system SHALL have comprehensive test coverage including unit and integration tests.

#### Scenario: Unit tests for use cases
- **GIVEN** use cases are implemented
- **WHEN** running unit tests
- **THEN** all use cases have tests covering success and failure paths with mocked repositories

#### Scenario: Integration tests for API endpoints
- **GIVEN** the Express server is running
- **WHEN** running integration tests
- **THEN** all endpoints have tests covering various scenarios

#### Scenario: Repository tests with in-memory database
- **GIVEN** the repository implementation exists
- **WHEN** running repository tests
- **THEN** tests use an in-memory SQLite database for isolation

### Requirement: Server Execution and Verification

The system SHALL be runnable using skills and all APIs shall be verifiable.

#### Scenario: Server startup
- **GIVEN** the project is built
- **WHEN** running the server using the appropriate skill
- **THEN** the Express server starts on the configured port

#### Scenario: API health check
- **GIVEN** the server is running
- **WHEN** all API endpoints are tested
- **THEN** each endpoint responds correctly according to its specification
