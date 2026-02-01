"""Integration tests for SQLiteItemRepository."""

import pytest

from src.infrastructure.persistence.sqlite_item_repository import SQLiteItemRepository
from src.domain.entities.item import CreateItemData, UpdateItemData, PatchItemData
from src.domain.entities.enums import Condition, Status
from src.domain.value_objects.item_image import ItemImage
from src.domain.repositories.item_repository import SortOptions, PaginationOptions, FilterOptions, SearchOptions


@pytest.mark.asyncio
class TestSQLiteItemRepository:
    """Integration tests for SQLiteItemRepository."""

    async def test_create_item(self, test_repository: SQLiteItemRepository) -> None:
        """Test creating an item."""
        data = CreateItemData(
            title="Test Item",
            description="Test Description",
            price_cents=10000,
            condition=Condition.GOOD,
            status=Status.ACTIVE,
        )

        item = await test_repository.create(data)

        assert item.id > 0
        assert item.title == "Test Item"
        assert item.description == "Test Description"
        assert item.price_cents == 10000
        assert item.condition == Condition.GOOD
        assert item.status == Status.ACTIVE
        assert item.is_featured is False
        assert item.country == "FR"
        assert item.delivery_available is False
        assert item.images == []
        assert item.created_at is not None
        assert item.updated_at is not None

    async def test_find_by_id(self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> None:
        """Test finding an item by ID."""
        created = await test_repository.create(sample_item_data)

        found = await test_repository.find_by_id(created.id)

        assert found is not None
        assert found.id == created.id
        assert found.title == created.title

    async def test_find_by_id_not_found(self, test_repository: SQLiteItemRepository) -> None:
        """Test finding a non-existent item."""
        found = await test_repository.find_by_id(99999)
        assert found is None

    async def test_find_all_empty(self, test_repository: SQLiteItemRepository) -> None:
        """Test finding all items when database is empty."""
        result = await test_repository.find_all()

        assert result.items == []
        assert result.next_cursor is None

    async def test_find_all_with_items(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test finding all items."""
        await test_repository.create(sample_item_data)
        await test_repository.create(
            CreateItemData(title="Another Item", price_cents=5000, condition=Condition.FAIR)
        )

        result = await test_repository.find_all()

        assert len(result.items) == 2

    async def test_find_all_with_filters(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test finding items with filters."""
        await test_repository.create(sample_item_data)
        await test_repository.create(
            CreateItemData(
                title="Another Item",
                price_cents=5000,
                condition=Condition.FAIR,
                status=Status.DRAFT,
                category="books",
            )
        )

        # Filter by status
        result = await test_repository.find_all(filters=FilterOptions(status=Status.ACTIVE))
        assert len(result.items) == 1
        assert result.items[0].status == Status.ACTIVE

        # Filter by category
        result = await test_repository.find_all(filters=FilterOptions(category="electronics"))
        assert len(result.items) == 1

        # Filter by price range
        result = await test_repository.find_all(filters=FilterOptions(min_price_cents=10000))
        assert len(result.items) == 1

    async def test_find_all_with_sorting(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test finding items with sorting."""
        await test_repository.create(sample_item_data)
        await test_repository.create(
            CreateItemData(title="B Item", price_cents=5000, condition=Condition.FAIR)
        )
        await test_repository.create(
            CreateItemData(title="A Item", price_cents=20000, condition=Condition.NEW)
        )

        # Sort by price ascending
        result = await test_repository.find_all(sort=SortOptions(field="price_cents", direction="asc"))
        assert result.items[0].price_cents == 5000
        assert result.items[-1].price_cents == 20000

    async def test_find_all_with_pagination(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test finding items with pagination."""
        for i in range(5):
            await test_repository.create(
                CreateItemData(title=f"Item {i}", price_cents=1000 * i, condition=Condition.GOOD)
            )

        # First page
        result = await test_repository.find_all(pagination=PaginationOptions(limit=2))
        assert len(result.items) == 2
        assert result.next_cursor is not None

        # Second page
        result2 = await test_repository.find_all(pagination=PaginationOptions(limit=2, cursor=result.next_cursor))
        assert len(result2.items) == 2
        assert result2.next_cursor is not None

    async def test_search(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test full-text search."""
        await test_repository.create(sample_item_data)
        await test_repository.create(
            CreateItemData(
                title="Vintage Bicycle",
                description="An old bicycle in good condition",
                price_cents=5000,
                condition=Condition.FAIR,
            )
        )

        result = await test_repository.search(SearchOptions(query="vintage"))

        assert len(result.items) >= 1

    async def test_update(self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> None:
        """Test updating an item."""
        created = await test_repository.create(sample_item_data)

        update_data = UpdateItemData(
            title="Updated Title",
            price_cents=20000,
            condition=Condition.LIKE_NEW,
            status=Status.SOLD,
            is_featured=True,
            country="FR",
            delivery_available=True,
        )

        updated = await test_repository.update(created.id, update_data)

        assert updated is not None
        assert updated.id == created.id
        assert updated.title == "Updated Title"
        assert updated.price_cents == 20000
        assert updated.condition == Condition.LIKE_NEW
        assert updated.status == Status.SOLD

    async def test_update_not_found(self, test_repository: SQLiteItemRepository) -> None:
        """Test updating a non-existent item."""
        update_data = UpdateItemData(
            title="Title",
            price_cents=10000,
            condition=Condition.GOOD,
            status=Status.ACTIVE,
            is_featured=False,
            country="FR",
            delivery_available=False,
        )

        result = await test_repository.update(99999, update_data)
        assert result is None

    async def test_patch(self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> None:
        """Test partially updating an item."""
        created = await test_repository.create(sample_item_data)

        patch_data = PatchItemData(price_cents=25000, status=Status.SOLD)

        patched = await test_repository.patch(created.id, patch_data)

        assert patched is not None
        assert patched.id == created.id
        assert patched.title == created.title  # Unchanged
        assert patched.price_cents == 25000  # Changed
        assert patched.status == Status.SOLD  # Changed

    async def test_patch_not_found(self, test_repository: SQLiteItemRepository) -> None:
        """Test patching a non-existent item."""
        patch_data = PatchItemData(price_cents=25000)

        result = await test_repository.patch(99999, patch_data)
        assert result is None

    async def test_delete(self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> None:
        """Test deleting an item."""
        created = await test_repository.create(sample_item_data)

        deleted = await test_repository.delete(created.id)

        assert deleted is True

        # Verify item is gone
        found = await test_repository.find_by_id(created.id)
        assert found is None

    async def test_delete_not_found(self, test_repository: SQLiteItemRepository) -> None:
        """Test deleting a non-existent item."""
        deleted = await test_repository.delete(99999)
        assert deleted is False

    async def test_exists(self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> None:
        """Test checking if an item exists."""
        created = await test_repository.create(sample_item_data)

        assert await test_repository.exists(created.id) is True
        assert await test_repository.exists(99999) is False

    async def test_create_with_images(
        self, test_repository: SQLiteItemRepository, sample_item_data: CreateItemData
    ) -> None:
        """Test creating an item with images."""
        sample_item_data.images = [
            ItemImage(url="https://example.com/img1.jpg", alt="Image 1", sort_order=0),
            ItemImage(url="https://example.com/img2.jpg", alt="Image 2", sort_order=1),
        ]

        item = await test_repository.create(sample_item_data)

        assert len(item.images) == 2
        assert item.images[0].url == "https://example.com/img1.jpg"
        assert item.images[1].url == "https://example.com/img2.jpg"
