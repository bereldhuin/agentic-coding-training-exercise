"""Integration tests for HTTP endpoints."""

import pytest
from fastapi.testclient import TestClient

from src.domain.entities.enums import Condition, Status


@pytest.mark.integration
class TestHTTPEndpoints:
    """Integration tests for FastAPI endpoints."""

    def test_health_check(self, client: TestClient) -> None:
        """Test health check endpoint."""
        response = client.get("/health")

        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "ok"
        assert "timestamp" in data

    def test_root_endpoint(self, client: TestClient) -> None:
        """Test root endpoint."""
        response = client.get("/")

        assert response.status_code == 200
        data = response.json()
        assert data["name"] == "LeBonPoint API (Python)"
        assert "docs" in data
        assert "health" in data

    def test_create_item(self, client: TestClient) -> None:
        """Test creating an item via POST /v1/items."""
        item_data = {
            "title": "Vintage Camera",
            "description": "A beautiful vintage camera",
            "price_cents": 15000,
            "category": "electronics",
            "condition": "good",
            "status": "active",
            "is_featured": True,
            "city": "Paris",
            "postal_code": "75001",
            "country": "FR",
            "delivery_available": True,
            "images": [
                {"url": "https://example.com/camera.jpg", "alt": "Front view", "sort_order": 0}
            ],
        }

        response = client.post("/v1/items", json=item_data)

        assert response.status_code == 201
        data = response.json()
        assert data["id"] > 0
        assert data["title"] == "Vintage Camera"
        assert data["price_cents"] == 15000
        assert data["condition"] == "good"
        assert data["status"] == "active"
        assert data["is_featured"] is True
        assert len(data["images"]) == 1

    def test_create_item_validation_error(self, client: TestClient) -> None:
        """Test creating an item with validation errors."""
        item_data = {
            "title": "AB",  # Too short
            "price_cents": -100,  # Negative
            "condition": "good",
        }

        response = client.post("/v1/items", json=item_data)

        assert response.status_code == 400
        data = response.json()
        assert "error" in data

    def test_get_item(self, client: TestClient) -> None:
        """Test getting an item via GET /v1/items/{id}."""
        # First create an item
        create_data = {
            "title": "Test Item",
            "price_cents": 10000,
            "condition": "good",
        }
        create_response = client.post("/v1/items", json=create_data)
        item_id = create_response.json()["id"]

        # Get the item
        response = client.get(f"/v1/items/{item_id}")

        assert response.status_code == 200
        data = response.json()
        assert data["id"] == item_id
        assert data["title"] == "Test Item"

    def test_get_item_not_found(self, client: TestClient) -> None:
        """Test getting a non-existent item."""
        response = client.get("/v1/items/99999")

        assert response.status_code == 404
        data = response.json()
        assert "error" in data

    def test_list_items_empty(self, client: TestClient) -> None:
        """Test listing items when database is empty."""
        response = client.get("/v1/items")

        assert response.status_code == 200
        data = response.json()
        assert "items" in data
        assert data["items"] == []
        assert data["next_cursor"] is None

    def test_list_items_with_data(self, client: TestClient) -> None:
        """Test listing items with data."""
        # Create a few items
        for i in range(3):
            client.post(
                "/v1/items",
                json={
                    "title": f"Item {i}",
                    "price_cents": 10000 * (i + 1),
                    "condition": "good",
                },
            )

        response = client.get("/v1/items")

        assert response.status_code == 200
        data = response.json()
        assert len(data["items"]) == 3

    def test_list_items_with_filters(self, client: TestClient) -> None:
        """Test listing items with filters."""
        # Create items with different statuses
        client.post(
            "/v1/items",
            json={"title": "Active Item", "price_cents": 10000, "condition": "good", "status": "active"},
        )
        client.post(
            "/v1/items",
            json={"title": "Draft Item", "price_cents": 10000, "condition": "good", "status": "draft"},
        )

        # Filter by status
        response = client.get("/v1/items?status=active")

        assert response.status_code == 200
        data = response.json()
        assert all(item["status"] == "active" for item in data["items"])

    def test_list_items_with_search(self, client: TestClient) -> None:
        """Test listing items with search query."""
        # Create items
        client.post(
            "/v1/items",
            json={"title": "Vintage Camera", "description": "An old camera", "price_cents": 10000, "condition": "good"},
        )
        client.post(
            "/v1/items",
            json={"title": "Modern Laptop", "description": "A new laptop", "price_cents": 50000, "condition": "new"},
        )

        # Search for "vintage"
        response = client.get("/v1/items?q=vintage")

        assert response.status_code == 200
        data = response.json()
        assert len(data["items"]) >= 1
        assert any("vintage" in item["title"].lower() or "vintage" in (item.get("description") or "").lower()
                   for item in data["items"])

    def test_list_items_with_sorting(self, client: TestClient) -> None:
        """Test listing items with sorting."""
        # Create items
        client.post("/v1/items", json={"title": "B Item", "price_cents": 20000, "condition": "good"})
        client.post("/v1/items", json={"title": "A Item", "price_cents": 10000, "condition": "good"})

        # Sort by price ascending
        response = client.get("/v1/items?sort_by=price_cents&sort_order=asc")

        assert response.status_code == 200
        data = response.json()
        prices = [item["price_cents"] for item in data["items"]]
        assert prices == sorted(prices)

    def test_list_items_with_pagination(self, client: TestClient) -> None:
        """Test listing items with pagination."""
        # Create multiple items
        for i in range(5):
            client.post(
                "/v1/items",
                json={"title": f"Item {i}", "price_cents": 10000, "condition": "good"},
            )

        # Get first page
        response = client.get("/v1/items?limit=2")

        assert response.status_code == 200
        data = response.json()
        assert len(data["items"]) == 2
        assert data["next_cursor"] is not None

        # Get second page
        cursor = data["next_cursor"]
        response2 = client.get(f"/v1/items?limit=2&cursor={cursor}")

        assert response2.status_code == 200
        data2 = response2.json()
        assert len(data2["items"]) == 2

    def test_update_item_put(self, client: TestClient) -> None:
        """Test updating an item via PUT /v1/items/{id}."""
        # Create an item
        create_response = client.post(
            "/v1/items",
            json={"title": "Original Title", "price_cents": 10000, "condition": "good", "status": "active",
                  "is_featured": False, "country": "FR", "delivery_available": False},
        )
        item_id = create_response.json()["id"]

        # Update the item
        update_data = {
            "title": "Updated Title",
            "price_cents": 20000,
            "condition": "like_new",
            "status": "sold",
            "is_featured": True,
            "country": "FR",
            "delivery_available": True,
        }

        response = client.put(f"/v1/items/{item_id}", json=update_data)

        assert response.status_code == 200
        data = response.json()
        assert data["title"] == "Updated Title"
        assert data["price_cents"] == 20000
        assert data["condition"] == "like_new"
        assert data["status"] == "sold"

    def test_patch_item(self, client: TestClient) -> None:
        """Test partially updating an item via PATCH /v1/items/{id}."""
        # Create an item
        create_response = client.post(
            "/v1/items",
            json={"title": "Original Title", "price_cents": 10000, "condition": "good"},
        )
        item_id = create_response.json()["id"]

        # Patch the item
        patch_data = {"price_cents": 15000}

        response = client.patch(f"/v1/items/{item_id}", json=patch_data)

        assert response.status_code == 200
        data = response.json()
        assert data["title"] == "Original Title"  # Unchanged
        assert data["price_cents"] == 15000  # Changed

    def test_delete_item(self, client: TestClient) -> None:
        """Test deleting an item via DELETE /v1/items/{id}."""
        # Create an item
        create_response = client.post(
            "/v1/items",
            json={"title": "To Delete", "price_cents": 10000, "condition": "good"},
        )
        item_id = create_response.json()["id"]

        # Delete the item
        response = client.delete(f"/v1/items/{item_id}")

        assert response.status_code == 204

        # Verify item is gone
        get_response = client.get(f"/v1/items/{item_id}")
        assert get_response.status_code == 404

    def test_openapi_spec(self, client: TestClient) -> None:
        """Test that OpenAPI spec is available."""
        response = client.get("/openapi.json")

        assert response.status_code == 200
        data = response.json()
        assert "openapi" in data
        assert "info" in data
        assert "paths" in data

    def test_docs_endpoint(self, client: TestClient) -> None:
        """Test that Swagger UI docs endpoint works."""
        response = client.get("/docs")

        assert response.status_code == 200
        assert "text/html" in response.headers.get("content-type", "")
