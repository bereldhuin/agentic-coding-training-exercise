"""Delete item use case."""

from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import NotFoundError, ValidationError


class DeleteItemUseCase:
    """Use case for deleting an item."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, id: int) -> None:
        """
        Delete an item.

        Args:
            id: Item ID

        Raises:
            ValidationError: If ID is invalid
            NotFoundError: If item not found
        """
        if id <= 0:
            raise ValidationError("Invalid item ID", {"id": "must be positive"})

        # Check if item exists
        if not await self._repository.exists(id):
            raise NotFoundError(f"Item with ID {id} not found")

        # Delete the item
        deleted = await self._repository.delete(id)
        if not deleted:
            raise NotFoundError(f"Item with ID {id} not found")
