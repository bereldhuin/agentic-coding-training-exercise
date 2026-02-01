"""FastAPI routes for items endpoints."""

from typing import Any
from fastapi import APIRouter, HTTPException, Query, status

from ...application.use_cases import (
    CreateItemUseCase,
    GetItemUseCase,
    ListItemsUseCase,
    UpdateItemUseCase,
    PatchItemUseCase,
    DeleteItemUseCase,
)
from ...domain.repositories.item_repository import ItemRepository
from ...shared.errors import ValidationError as AppValidationError, NotFoundError
from .models import (
    ItemResponseModel,
    CreateItemRequestModel,
    UpdateItemRequestModel,
    PatchItemRequestModel,
    ListItemsResponseModel,
    to_response_model,
    to_create_data,
    to_update_data,
    to_patch_data,
)


def create_items_router(repository: ItemRepository) -> APIRouter:
    """
    Create the items router.

    Args:
        repository: Item repository instance

    Returns:
        Configured APIRouter
    """
    router = APIRouter(prefix="/v1/items", tags=["items"])

    # Initialize use cases
    create_item_use_case = CreateItemUseCase(repository)
    get_item_use_case = GetItemUseCase(repository)
    list_items_use_case = ListItemsUseCase(repository)
    update_item_use_case = UpdateItemUseCase(repository)
    patch_item_use_case = PatchItemUseCase(repository)
    delete_item_use_case = DeleteItemUseCase(repository)

    @router.get("", response_model=ListItemsResponseModel, summary="List items")
    async def list_items(
        q: str | None = Query(None, description="Search query for full-text search"),
        status: str | None = Query(None, description="Filter by status"),
        category: str | None = Query(None, description="Filter by category"),
        min_price_cents: int | None = Query(None, description="Minimum price in cents"),
        max_price_cents: int | None = Query(None, description="Maximum price in cents"),
        city: str | None = Query(None, description="Filter by city"),
        postal_code: str | None = Query(None, description="Filter by postal code"),
        is_featured: bool | None = Query(None, description="Filter by featured status"),
        delivery_available: bool | None = Query(None, description="Filter by delivery availability"),
        sort: str | None = Query(None, description="Sort by field and direction (field:direction)"),
        limit: int | None = Query(None, description="Number of items per page"),
        cursor: str | None = Query(None, description="Pagination cursor"),
    ) -> dict[str, Any]:
        """
        List items with optional filtering, sorting, and pagination.

        Supports full-text search via the `q` parameter.
        """
        params = {
            "q": q,
            "status": status,
            "category": category,
            "min_price_cents": min_price_cents,
            "max_price_cents": max_price_cents,
            "city": city,
            "postal_code": postal_code,
            "is_featured": is_featured,
            "delivery_available": delivery_available,
            "sort": sort,
            "limit": limit,
            "cursor": cursor,
        }
        result = await list_items_use_case.execute(params)
        return {
            "items": [to_response_model(item) for item in result["items"]],
            "next_cursor": result["next_cursor"],
        }

    @router.post(
        "",
        response_model=ItemResponseModel,
        status_code=status.HTTP_201_CREATED,
        summary="Create a new item",
    )
    async def create_item(data: CreateItemRequestModel) -> ItemResponseModel:
        """Create a new item."""
        try:
            create_data = to_create_data(data)
            item = await create_item_use_case.execute(create_data)
            return to_response_model(item)
        except (AppValidationError, NotFoundError) as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST if isinstance(e, AppValidationError) else status.HTTP_404_NOT_FOUND,
                detail={"error": {"code": e.__class__.__name__.lower(), "message": e.message, "details": e.details}}
            )

    @router.get("/{item_id}", response_model=ItemResponseModel, summary="Get an item by ID")
    async def get_item(item_id: int) -> ItemResponseModel:
        """Get a single item by ID."""
        try:
            item = await get_item_use_case.execute(item_id)
            return to_response_model(item)
        except (AppValidationError, NotFoundError) as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST if isinstance(e, AppValidationError) else status.HTTP_404_NOT_FOUND,
                detail={"error": {"code": e.__class__.__name__.lower(), "message": e.message, "details": e.details}}
            )

    @router.put("/{item_id}", response_model=ItemResponseModel, summary="Update an item (full replace)")
    async def update_item(item_id: int, data: UpdateItemRequestModel) -> ItemResponseModel:
        """Update (replace) an item."""
        try:
            update_data = to_update_data(data)
            item = await update_item_use_case.execute(item_id, update_data)
            return to_response_model(item)
        except (AppValidationError, NotFoundError) as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST if isinstance(e, AppValidationError) else status.HTTP_404_NOT_FOUND,
                detail={"error": {"code": e.__class__.__name__.lower(), "message": e.message, "details": e.details}}
            )

    @router.patch(
        "{item_id}",
        response_model=ItemResponseModel,
        summary="Partially update an item",
    )
    async def patch_item(item_id: int, data: PatchItemRequestModel) -> ItemResponseModel:
        """Partially update an item."""
        try:
            patch_data = to_patch_data(data)
            item = await patch_item_use_case.execute(item_id, patch_data)
            return to_response_model(item)
        except (AppValidationError, NotFoundError) as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST if isinstance(e, AppValidationError) else status.HTTP_404_NOT_FOUND,
                detail={"error": {"code": e.__class__.__name__.lower(), "message": e.message, "details": e.details}}
            )

    @router.delete(
        "{item_id}",
        status_code=status.HTTP_204_NO_CONTENT,
        summary="Delete an item",
    )
    async def delete_item(item_id: int) -> None:
        """Delete an item."""
        try:
            await delete_item_use_case.execute(item_id)
        except (AppValidationError, NotFoundError) as e:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST if isinstance(e, AppValidationError) else status.HTTP_404_NOT_FOUND,
                detail={"error": {"code": e.__class__.__name__.lower(), "message": e.message, "details": e.details}}
            )

    return router
