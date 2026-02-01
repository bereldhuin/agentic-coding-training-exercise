"""Update item use case (PUT - full replacement)."""

from ...domain.entities.item import Item, UpdateItemData
from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import NotFoundError, ValidationError


class UpdateItemUseCase:
    """Use case for updating (replacing) an item."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, id: int, data: UpdateItemData) -> Item:
        """
        Update an item (full replace).

        Args:
            id: Item ID
            data: New item data (full replace)

        Returns:
            Updated item

        Raises:
            ValidationError: If ID or data validation fails
            NotFoundError: If item not found
        """
        if id <= 0:
            raise ValidationError("Invalid item ID", {"id": "must be positive"})

        # Validate title length
        if len(data.title) < 3:
            raise ValidationError("Title must be at least 3 characters", {"title": "too_short"})
        if len(data.title) > 200:
            raise ValidationError("Title must be at most 200 characters", {"title": "too_long"})

        # Validate price
        if data.price_cents < 0:
            raise ValidationError("Price cannot be negative", {"price_cents": "invalid"})

        # Check if item exists
        if not await self._repository.exists(id):
            raise NotFoundError(f"Item with ID {id} not found")

        # Update the item
        updated = await self._repository.update(id, data)
        if updated is None:
            raise NotFoundError(f"Item with ID {id} not found")

        return updated
