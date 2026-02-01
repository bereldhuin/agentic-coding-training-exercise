#!/bin/bash

# Validation script for LeBonPoint OpenAPI specification
# This script validates the canonical OpenAPI specification and optionally compares server implementations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SPEC_FILE="$SCRIPT_DIR/openapi.yaml"

echo "🔍 Validating OpenAPI specification: $SPEC_FILE"
echo ""

# Check if openapi-spec-validator is installed
if ! command -v openapi-spec-validator &> /dev/null; then
    echo "⚠️  openapi-spec-validator not found. Installing..."
    pip install openapi-spec-validator
fi

# Validate the specification
echo "Validating OpenAPI 3.1 specification..."
if openapi-spec-validator "$SPEC_FILE"; then
    echo "✅ OpenAPI specification is valid!"
else
    echo "❌ OpenAPI specification has errors!"
    exit 1
fi

echo ""
echo "📊 Specification Info:"
echo "======================"

# Extract info using grep/awk
TITLE=$(grep "^title:" "$SPEC_FILE" | head -1 | sed 's/title: //')
VERSION=$(grep "^version:" "$SPEC_FILE" | head -1 | sed 's/version: //')

echo "Title: $TITLE"
echo "Version: $VERSION"

# Count endpoints
ENDPOINTS=$(grep -E "^\s+(get|post|put|patch|delete):" "$SPEC_FILE" | wc -l | tr -d ' ')
echo "Endpoints: $ENDPOINTS"

# Count schemas
SCHEMAS=$(grep -E "^  [A-Z][a-zA-Z]+:$" "$SPEC_FILE" | wc -l | tr -d ' ')
echo "Schemas: $SCHEMAS"

echo ""
echo "✅ Validation complete!"
echo ""
echo "To compare with server implementations:"
echo "  - Python (FastAPI): curl http://localhost:8000/openapi.json"
echo "  - TypeScript: Add endpoint to serve spec, then curl http://localhost:3000/openapi.json"
echo ""
echo "For more tools, see README.md"
