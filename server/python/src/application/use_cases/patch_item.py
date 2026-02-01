"""Patch item use case (PATCH - partial update)."""

from ...domain.entities.item import Item, PatchItemData
from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import NotFoundError, ValidationError


class PatchItemUseCase:
    """Use case for partially updating an item."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, id: int, data: PatchItemData) -> Item:
        """
        Partially update an item.

        Args:
            id: Item ID
            data: Partial item data to update

        Returns:
            Updated item

        Raises:
            ValidationError: If ID or data validation fails
            NotFoundError: If item not found
        """
        if id <= 0:
            raise ValidationError("Invalid item ID", {"id": "must be positive"})

        # Validate title if provided
        if data.title is not None:
            if len(data.title) < 3:
                raise ValidationError("Title must be at least 3 characters", {"title": "too_short"})
            if len(data.title) > 200:
                raise ValidationError("Title must be at most 200 characters", {"title": "too_long"})

        # Validate price if provided
        if data.price_cents is not None and data.price_cents < 0:
            raise ValidationError("Price cannot be negative", {"price_cents": "invalid"})

        # Check if item exists
        if not await self._repository.exists(id):
            raise NotFoundError(f"Item with ID {id} not found")

        # Update the item
        updated = await self._repository.patch(id, data)
        if updated is None:
            raise NotFoundError(f"Item with ID {id} not found")

        return updated
