"""Tests for CreateItemUseCase."""

from unittest.mock import AsyncMock, Mock

import pytest

from src.application.use_cases.create_item import CreateItemUseCase
from src.domain.entities.item import Item, CreateItemData
from src.domain.entities.enums import Condition, Status
from src.domain.value_objects.item_image import ItemImage
from src.shared.errors import ValidationError


@pytest.fixture
def mock_repository():
    """Create a mock repository."""
    repo = Mock()
    repo.create = AsyncMock()
    return repo


@pytest.fixture
def create_item_data():
    """Create sample item data."""
    return CreateItemData(
        title="Vintage Camera",
        description="A beautiful vintage camera",
        price_cents=15000,
        condition=Condition.GOOD,
        status=Status.ACTIVE,
    )


@pytest.mark.asyncio
class TestCreateItemUseCase:
    """Tests for CreateItemUseCase."""

    async def test_execute_success(self, mock_repository, create_item_data):
        """Test successful item creation."""
        # Setup
        expected_item = Item(
            id=1,
            title="Vintage Camera",
            description="A beautiful vintage camera",
            price_cents=15000,
            category=None,
            condition=Condition.GOOD,
            status=Status.ACTIVE,
            is_featured=False,
            city=None,
            postal_code=None,
            country="FR",
            delivery_available=False,
            created_at="2026-01-30T17:36:03.653Z",
            updated_at="2026-01-30T17:36:03.653Z",
            published_at=None,
            images=[],
        )
        mock_repository.create.return_value = expected_item

        # Execute
        use_case = CreateItemUseCase(mock_repository)
        result = await use_case.execute(create_item_data)

        # Verify
        assert result == expected_item
        mock_repository.create.assert_called_once_with(create_item_data)

    async def test_execute_title_too_short(self, mock_repository):
        """Test validation error for title too short."""
        use_case = CreateItemUseCase(mock_repository)

        data = CreateItemData(
            title="AB",  # Too short
            price_cents=10000,
            condition=Condition.GOOD,
        )

        with pytest.raises(ValidationError) as exc_info:
            await use_case.execute(data)

        assert exc_info.value.details == {"title": "too_short"}
        mock_repository.create.assert_not_called()

    async def test_execute_title_too_long(self, mock_repository):
        """Test validation error for title too long."""
        use_case = CreateItemUseCase(mock_repository)

        data = CreateItemData(
            title="A" * 201,  # Too long
            price_cents=10000,
            condition=Condition.GOOD,
        )

        with pytest.raises(ValidationError) as exc_info:
            await use_case.execute(data)

        assert exc_info.value.details == {"title": "too_long"}
        mock_repository.create.assert_not_called()

    async def test_execute_negative_price(self, mock_repository):
        """Test validation error for negative price."""
        use_case = CreateItemUseCase(mock_repository)

        data = CreateItemData(
            title="Valid Title",
            price_cents=-100,  # Negative
            condition=Condition.GOOD,
        )

        with pytest.raises(ValidationError) as exc_info:
            await use_case.execute(data)

        assert exc_info.value.details == {"price_cents": "invalid"}
        mock_repository.create.assert_not_called()
