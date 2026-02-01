"""Get item use case."""

from ...domain.entities.item import Item
from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import NotFoundError, ValidationError


class GetItemUseCase:
    """Use case for retrieving a single item."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, id: int) -> Item:
        """
        Get an item by ID.

        Args:
            id: Item ID

        Returns:
            Item

        Raises:
            ValidationError: If ID is invalid
            NotFoundError: If item not found
        """
        if id <= 0:
            raise ValidationError("Invalid item ID", {"id": "must be positive"})

        item = await self._repository.find_by_id(id)
        if item is None:
            raise NotFoundError(f"Item with ID {id} not found")

        return item
