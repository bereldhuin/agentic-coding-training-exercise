"""Item repository port (abstract interface)."""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from enum import Enum
from typing import Callable

from ..entities.item import Item, CreateItemData, UpdateItemData, PatchItemData


class SortField(str, Enum):
    """Fields available for sorting."""

    CREATED_AT = "created_at"
    UPDATED_AT = "updated_at"
    PUBLISHED_AT = "published_at"
    PRICE_CENTS = "price_cents"
    TITLE = "title"
    ID = "id"


class SortDirection(str, Enum):
    """Sort direction options."""

    ASC = "asc"
    DESC = "desc"


@dataclass
class FilterOptions:
    """Filter options for item listing."""

    status: str | None = None
    category: str | None = None
    city: str | None = None
    postal_code: str | None = None
    is_featured: bool | None = None
    delivery_available: bool | None = None
    garantie_months: int | None = None


@dataclass
class SortOptions:
    """Sort options for item listing."""

    field: SortField = SortField.CREATED_AT
    direction: SortDirection = SortDirection.DESC

    def __post_init__(self) -> None:
        if not isinstance(self.field, SortField):
            self.field = SortField(self.field)
        if not isinstance(self.direction, SortDirection):
            self.direction = SortDirection(self.direction)


@dataclass
class PaginationOptions:
    """Pagination options for item listing."""

    limit: int = 20
    cursor: str | None = None


@dataclass
class SearchOptions:
    """Options for full-text search."""

    query: str
    filters: FilterOptions | None = None
    pagination: PaginationOptions | None = None


@dataclass
class PaginatedResult:
    """Paginated result from repository queries."""

    items: list[Item]
    next_cursor: str | None = None


class ItemRepository(ABC):
    """
    Abstract repository port for Item persistence.

    This defines the interface that any repository implementation must follow.
    """

    @abstractmethod
    async def create(self, data: CreateItemData) -> Item:
        """
        Create a new item.

        Args:
            data: Item data for creation

        Returns:
            Created item with generated id and timestamps

        Raises:
            ConflictError: If item with same unique constraints exists
        """
        pass

    @abstractmethod
    async def find_by_id(self, id: int) -> Item | None:
        """
        Find an item by ID.

        Args:
            id: Item ID

        Returns:
            Item if found, None otherwise
        """
        pass

    @abstractmethod
    async def find_all(
        self,
        filters: FilterOptions | None = None,
        sort: SortOptions | None = None,
        pagination: PaginationOptions | None = None,
    ) -> PaginatedResult:
        """
        List all items with optional filtering, sorting, and pagination.

        Args:
            filters: Optional filters to apply
            sort: Optional sort options
            pagination: Optional pagination options

        Returns:
            Paginated result with items and next cursor
        """
        pass

    @abstractmethod
    async def update(self, id: int, data: UpdateItemData) -> Item | None:
        """
        Update (replace) an item.

        Args:
            id: Item ID
            data: New item data (full replace)

        Returns:
            Updated item if found, None if not found
        """
        pass

    @abstractmethod
    async def patch(self, id: int, data: PatchItemData) -> Item | None:
        """
        Partially update an item.

        Args:
            id: Item ID
            data: Partial item data to update

        Returns:
            Updated item if found, None if not found
        """
        pass

    @abstractmethod
    async def delete(self, id: int) -> bool:
        """
        Delete an item.

        Args:
            id: Item ID

        Returns:
            True if deleted, False if not found
        """
        pass

    @abstractmethod
    async def search(self, options: SearchOptions) -> PaginatedResult:
        """
        Full-text search for items using FTS5.

        Args:
            options: Search options including query, filters, and pagination

        Returns:
            Paginated result with matching items
        """
        pass

    @abstractmethod
    async def exists(self, id: int) -> bool:
        """
        Check if an item exists.

        Args:
            id: Item ID

        Returns:
            True if exists, False otherwise
        """
        pass
