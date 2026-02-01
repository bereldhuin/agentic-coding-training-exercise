"""Tests for GetItemUseCase."""

from unittest.mock import AsyncMock, Mock

import pytest

from src.application.use_cases.get_item import GetItemUseCase
from src.domain.entities.item import Item
from src.domain.entities.enums import Condition, Status
from src.shared.errors import NotFoundError, ValidationError


@pytest.fixture
def mock_repository():
    """Create a mock repository."""
    repo = Mock()
    repo.find_by_id = AsyncMock()
    return repo


@pytest.fixture
def sample_item():
    """Create a sample item."""
    return Item(
        id=1,
        title="Test Item",
        description="Test Description",
        price_cents=10000,
        category="electronics",
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


@pytest.mark.asyncio
class TestGetItemUseCase:
    """Tests for GetItemUseCase."""

    async def test_execute_success(self, mock_repository, sample_item):
        """Test successful item retrieval."""
        mock_repository.find_by_id.return_value = sample_item

        use_case = GetItemUseCase(mock_repository)
        result = await use_case.execute(1)

        assert result == sample_item
        mock_repository.find_by_id.assert_called_once_with(1)

    async def test_execute_not_found(self, mock_repository):
        """Test item not found."""
        mock_repository.find_by_id.return_value = None

        use_case = GetItemUseCase(mock_repository)

        with pytest.raises(NotFoundError) as exc_info:
            await use_case.execute(999)

        assert "Item with ID 999 not found" in str(exc_info.value)

    async def test_execute_invalid_id(self, mock_repository):
        """Test invalid item ID."""
        use_case = GetItemUseCase(mock_repository)

        with pytest.raises(ValidationError) as exc_info:
            await use_case.execute(0)

        assert exc_info.value.details == {"id": "must be positive"}

        with pytest.raises(ValidationError):
            await use_case.execute(-1)
