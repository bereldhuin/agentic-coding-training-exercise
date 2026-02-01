# LeBonPoint API Specification

This directory contains the canonical OpenAPI 3.1 specification for the LeBonPoint Marketplace API.

## Overview

The `openapi.yaml` file is the **single source of truth** for the API contract. All server implementations (TypeScript, Python, and future implementations) must conform to this specification.

## File Structure

```
server/api/
├── openapi.yaml      # Canonical OpenAPI 3.1 specification
├── README.md         # This file
└── validate.sh       # Validation script (optional)
```

## API-First Development Workflow

When making API changes, follow this workflow:

### 1. Design (Update Specification)

Start by updating `openapi.yaml` with your proposed changes:

- Add new endpoints
- Modify existing endpoints
- Add or modify schemas
- Update validation rules
- Add examples

### 2. Review

Create a pull request with the specification changes. Review should focus on:

- API design and consistency
- Endpoint naming and structure
- Schema completeness
- Validation rules
- Example accuracy

### 3. Implement

After the specification is approved, implement the changes in one or more servers:

- **TypeScript (Express)**: Update routes, handlers, validation schemas
- **Python (FastAPI)**: Update routes, Pydantic models, handlers

### 4. Validate

Ensure implementations match the specification:

- Compare server-generated specs to canonical spec
- Test endpoints with various inputs
- Verify error responses match specification
- Check validation rules are enforced

### 5. Deploy

Deploy updated servers to appropriate environments.

## Specification Contents

The OpenAPI specification includes:

### Endpoints

- **GET /health** - Health check
- **GET /v1/items** - List items with filtering, sorting, pagination, search
- **POST /v1/items** - Create a new item
- **GET /v1/items/{id}** - Get a single item
- **PUT /v1/items/{id}** - Update (replace) an item
- **PATCH /v1/items/{id}** - Partially update an item
- **DELETE /v1/items/{id}** - Delete an item

### Schemas

- **Item** - Full item model with all fields
- **ItemImage** - Image object with url, alt, sort_order
- **CreateItemRequest** - Request schema for creating items
- **UpdateItemRequest** - Request schema for updating items (all fields required)
- **PatchItemRequest** - Request schema for patching items (all fields optional)
- **ListItemsResponse** - Response schema for list endpoint
- **HealthResponse** - Health check response
- **ErrorResponse** - Standard error response format

### Enums

- **Condition**: new, like_new, good, fair, parts, unknown
- **Status**: draft, active, reserved, sold, archived
- **Error Codes**: validation_error, not_found, internal_error

## Validation

### Manual Validation

Validate the OpenAPI specification syntax:

```bash
# Using swagger-cli (recommended)
npm install -g @apidevtools/swagger-cli
swagger-cli validate server/api/openapi.yaml

# Using openapi-spec-validator (Python)
pip install openapi-spec-validator
openapi-spec-validator server/api/openapi.yaml
```

### Comparing Server Implementations

Compare server-generated specs to the canonical spec:

```bash
# Python FastAPI server
curl http://localhost:8000/openapi.json > /tmp/python-openapi.json

# Convert YAML to JSON for comparison
yq eval -o=json server/api/openapi.yaml > /tmp/canonical-openapi.json

# Compare (requires jq and diff)
diff <(jq -S . /tmp/canonical-openapi.json) <(jq -S . /tmp/python-openapi.json)
```

## Using the Specification

### Client Code Generation

Generate client SDKs from the OpenAPI specification:

```bash
# Using openapi-generator
docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
  -i /local/server/api/openapi.yaml \
  -g typescript-axios \
  -o /local/generated/client

# Using orval (for TypeScript)
npm install -g orval
orval --config ./orval.config.js
```

### Mock Server

Run a mock server for client development:

```bash
# Using Prism
npm install -g @stoplight/prism-cli
prism mock server/api/openapi.yaml

# Mock server will be available at http://localhost:4010
```

### Contract Testing

Test implementations against the specification:

```bash
# Using Schemathesis (Python)
pip install schemathesis
schemathesis run http://localhost:8000/openapi.json

# Using Dredd (Node.js)
npm install -g dredd
dredd server/api/openapi.yaml http://localhost:3000
```

## Server Integration

### TypeScript (Express)

The TypeScript server at `../typescript/` implements this specification. Key points:

- Uses Zod for validation that matches OpenAPI schemas
- Returns error responses in the standard ErrorResponse format
- References this specification as the source of truth
- Can serve the canonical spec at `/openapi.json` (optional)

### Python (FastAPI)

The Python server at `../python/` implements this specification. Key points:

- Uses Pydantic models that match OpenAPI schemas
- FastAPI auto-generates OpenAPI spec for `/docs` and `/redoc`
- The auto-generated spec should match this canonical spec
- References this specification as the authoritative source

## Versioning

- **Current API Version**: v1.0.0
- **API Path Prefix**: `/v1/`
- **Spec File**: `openapi.yaml` (no version suffix for v1)

When v2 is needed with breaking changes:

1. Create `openapi-v2.yaml`
2. Add v2 endpoints with `/v2/` path prefix
3. Maintain v1 spec and implementation for backward compatibility

## Related Documentation

- [OpenAPI 3.1 Specification](https://spec.openapis.org/oas/v3.1.0)
- [OpenAPI Tools](https://openapi.tools/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [TypeScript Server](../typescript/README.md)
- [Python Server](../python/README.md)

## Tools Reference

### Validation Tools

- **@apidevtools/swagger-cli** - Official OpenAPI/Swagger CLI tool
- **openapi-spec-validator** - Python OpenAPI validator
- **spectral** - Linter for OpenAPI specs

### Generation Tools

- **openapi-generator** - Generate clients, servers, documentation
- **orval** - TypeScript client generator with React Query integration
- **openapi-typescript** - Generate TypeScript types from OpenAPI

### Testing Tools

- **Schemathesis** - Property-based testing for OpenAPI
- **Dredd** - API testing against OpenAPI specification
- **Prism** - Mock server from OpenAPI spec

### Documentation Tools

- **Swagger UI** - Interactive API documentation
- **Redoc** - Beautiful API reference documentation
- **Stoplight Elements** - API documentation components

## Contributing

When contributing API changes:

1. **Always start with the spec** - Update `openapi.yaml` first
2. **Validate** - Run `swagger-cli validate` before committing
3. **Document** - Add clear descriptions and examples
4. **Test** - Verify implementations match the spec
5. **Review** - Get team approval on spec changes before implementation

## Troubleshooting

### Common Issues

**Issue**: Server-generated spec doesn't match canonical spec
- **Solution**: Update server implementation to match canonical spec. The canonical spec is authoritative.

**Issue**: Validation fails with `$ref` errors
- **Solution**: Check that all referenced components exist in `components/schemas`

**Issue**: Examples don't match validation rules
- **Solution**: Ensure all example values conform to the schema constraints

**Issue**: Type mismatch between TypeScript/Python and OpenAPI
- **Solution**: OpenAPI types should map correctly: integer -> number/Int32, string -> str, etc.

## Questions?

For questions about the API specification or development workflow, please refer to:

- Main project documentation: `../../README.md`
- TypeScript server: `../typescript/README.md`
- Python server: `../python/README.md`
- OpenSpec change proposals: `../../openspec/changes/`
