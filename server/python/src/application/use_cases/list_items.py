"""List items use case."""

from typing import Any

from ...domain.entities.item import Item
from ...domain.entities.enums import Status
from ...domain.repositories.item_repository import (
    ItemRepository,
    FilterOptions,
    SortOptions,
    PaginationOptions,
    SearchOptions,
)
from ...shared.errors import ValidationError


class ListItemsUseCase:
    """Use case for listing items with filtering and search."""

    def __init__(self, repository: ItemRepository) -> None:
        """
        Initialize the use case.

        Args:
            repository: Item repository port
        """
        self._repository = repository

    async def execute(self, query_params: dict[str, Any]) -> dict[str, Any]:
        """
        List items with optional filtering, sorting, and search.

        Args:
            query_params: Query parameters from request

        Returns:
            Dict with items list and next cursor
        """
        # Check if this is a search request
        search_query = query_params.get("q")

        # Parse filters
        status = query_params.get("status")
        if status is not None and status not in Status.values():
            raise ValidationError("Validation failed", {"status": "Invalid status value"})

        filters = FilterOptions(
            status=status,
            category=query_params.get("category"),
            min_price_cents=self._parse_int(query_params.get("min_price_cents")),
            max_price_cents=self._parse_int(query_params.get("max_price_cents")),
            city=query_params.get("city"),
            postal_code=query_params.get("postal_code"),
            is_featured=self._parse_bool(query_params.get("is_featured")),
            delivery_available=self._parse_bool(query_params.get("delivery_available")),
        )

        # Parse sort options
        sort = self._parse_sort(query_params.get("sort"))

        # Parse pagination
        limit = self._parse_int(query_params.get("limit"))
        if limit is None:
            limit = 20
        if limit < 1 or limit > 100:
            raise ValidationError("Validation failed", {"limit": "limit must be between 1 and 100"})
        cursor = query_params.get("cursor")
        pagination = PaginationOptions(limit=min(limit, 100), cursor=cursor)

        # Execute query
        if search_query:
            result = await self._repository.search(
                SearchOptions(query=search_query, filters=filters, pagination=pagination)
            )
        else:
            result = await self._repository.find_all(filters=filters, sort=sort, pagination=pagination)

        return {
            "items": result.items,
            "next_cursor": result.next_cursor,
        }

    def _parse_int(self, value: Any) -> int | None:
        """Parse an integer value."""
        if value is None:
            return None
        try:
            return int(value)
        except (ValueError, TypeError):
            return None

    def _parse_bool(self, value: Any) -> bool | None:
        """Parse a boolean value."""
        if value is None:
            return None
        if isinstance(value, bool):
            return value
        if isinstance(value, str):
            return value.lower() in ("true", "1", "yes")
        return None

    def _parse_sort(self, value: Any) -> SortOptions:
        """Parse sort options from string."""
        from ...domain.repositories.item_repository import SortField, SortDirection

        if value is None or value == "":
            return SortOptions()

        if not isinstance(value, str) or ":" not in value:
            raise ValidationError("Validation failed", {"sort": "sort must be in the form field:direction"})

        field_value, direction_value = value.split(":", 1)

        try:
            field = SortField(field_value)
        except ValueError as exc:
            raise ValidationError(
                "Validation failed",
                {
                    "sort": "sort field must be one of: id, title, price_cents, created_at, updated_at, published_at"
                },
            ) from exc

        try:
            direction = SortDirection(direction_value)
        except ValueError as exc:
            raise ValidationError("Validation failed", {"sort": "sort direction must be asc or desc"}) from exc

        return SortOptions(field=field, direction=direction)
