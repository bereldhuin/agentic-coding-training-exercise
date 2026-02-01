"""Tests for Item entity and related functions."""

from src.domain.entities.item import (
    Item,
    CreateItemData,
    UpdateItemData,
    PatchItemData,
    item_from_row,
    parse_images_column,
    serialize_images_column,
)
from src.domain.entities.enums import Condition, Status
from src.domain.value_objects.item_image import ItemImage


class TestItemEntity:
    """Tests for Item entity."""

    def test_item_creation(self) -> None:
        """Test creating an Item entity."""
        images = [
            ItemImage(url="https://example.com/img1.jpg", alt="Image 1"),
            ItemImage(url="https://example.com/img2.jpg", alt="Image 2"),
        ]

        item = Item(
            id=1,
            title="Test Item",
            description="Test Description",
            price_cents=10000,
            category="electronics",
            condition=Condition.GOOD,
            status=Status.ACTIVE,
            is_featured=True,
            city="Paris",
            postal_code="75001",
            country="FR",
            delivery_available=True,
            created_at="2026-01-30T17:36:03.653Z",
            updated_at="2026-01-30T17:36:03.653Z",
            published_at="2026-01-30T17:36:03.653Z",
            images=images,
        )

        assert item.id == 1
        assert item.title == "Test Item"
        assert item.price_cents == 10000
        assert item.condition == Condition.GOOD
        assert len(item.images) == 2


class TestCreateItemData:
    """Tests for CreateItemData."""

    def test_defaults(self) -> None:
        """Test default values for CreateItemData."""
        data = CreateItemData(
            title="Test Item",
            price_cents=10000,
            condition=Condition.GOOD,
        )

        assert data.title == "Test Item"
        assert data.description is None
        assert data.price_cents == 10000
        assert data.category is None
        assert data.condition == Condition.GOOD
        assert data.status == Status.DRAFT
        assert data.is_featured is False
        assert data.city is None
        assert data.postal_code is None
        assert data.country == "FR"
        assert data.delivery_available is False
        assert data.images is None


class TestPatchItemData:
    """Tests for PatchItemData."""

    def test_all_optional(self) -> None:
        """Test that all fields in PatchItemData are optional."""
        data = PatchItemData()
        assert data.title is None
        assert data.description is None
        assert data.price_cents is None
        assert data.condition is None
        assert data.status is None


class TestItemImageFunctions:
    """Tests for image-related functions."""

    def test_parse_images_column_empty(self) -> None:
        """Test parsing empty images column."""
        assert parse_images_column(None) == []
        assert parse_images_column("") == []

    def test_parse_images_column_valid(self) -> None:
        """Test parsing valid images column."""
        json_str = '[{"url": "https://example.com/img.jpg", "alt": "Test", "sort_order": 0}]'
        images = parse_images_column(json_str)

        assert len(images) == 1
        assert images[0].url == "https://example.com/img.jpg"
        assert images[0].alt == "Test"
        assert images[0].sort_order == 0

    def test_parse_images_column_invalid(self) -> None:
        """Test parsing invalid images column."""
        assert parse_images_column("invalid json") == []

    def test_serialize_images_column(self) -> None:
        """Test serializing images to JSON."""
        images = [
            ItemImage(url="https://example.com/img.jpg", alt="Test", sort_order=0)
        ]
        json_str = serialize_images_column(images)

        assert json_str == '[{"url": "https://example.com/img.jpg", "alt": "Test", "sort_order": 0}]'


class TestItemFromRow:
    """Tests for item_from_row function."""

    def test_from_row(self) -> None:
        """Test converting database row to Item entity."""
        row = {
            "id": 1,
            "title": "Test Item",
            "description": "Test Description",
            "price_cents": 10000,
            "category": "electronics",
            "condition": "good",
            "status": "active",
            "is_featured": 1,
            "city": "Paris",
            "postal_code": "75001",
            "country": "FR",
            "delivery_available": 1,
            "created_at": "2026-01-30T17:36:03.653Z",
            "updated_at": "2026-01-30T17:36:03.653Z",
            "published_at": "2026-01-30T17:36:03.653Z",
            "images": '[{"url": "https://example.com/img.jpg"}]',
        }

        item = item_from_row(row)

        assert item.id == 1
        assert item.title == "Test Item"
        assert item.price_cents == 10000
        assert item.condition == Condition.GOOD
        assert item.status == Status.ACTIVE
        assert item.is_featured is True
        assert item.delivery_available is True
