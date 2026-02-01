"""Create item use case."""

from ...domain.entities.item import Item, CreateItemData
from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import ValidationError, ConflictError


class CreateItemUseCase:
    """Use case for creating a new item."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, data: CreateItemData) -> Item:
        """
        Create a new item.

        Args:
            data: Item data for creation

        Returns:
            Created item

        Raises:
            ValidationError: If data validation fails
            ConflictError: If item violates unique constraints
        """
        # Validate title length
        if len(data.title) < 3:
            raise ValidationError("Title must be at least 3 characters", {"title": "too_short"})
        if len(data.title) > 200:
            raise ValidationError("Title must be at most 200 characters", {"title": "too_long"})

        # Validate price
        if data.price_cents < 0:
            raise ValidationError("Price cannot be negative", {"price_cents": "invalid"})

        # Create the item
        return await self._repository.create(data)
