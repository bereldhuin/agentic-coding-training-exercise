"""Pydantic models for request/response validation."""

from typing import Any
from pydantic import BaseModel, Field

from ...domain.entities.enums import Condition, Status


class ItemImageModel(BaseModel):
    """Model for item image."""

    url: str = Field(..., min_length=1)
    alt: str | None = None
    sort_order: int | None = None


class ItemResponseModel(BaseModel):
    """Model for item response."""

    id: int
    title: str
    description: str | None
    price_cents: int
    category: str | None
    condition: Condition
    status: Status
    is_featured: bool
    city: str | None
    postal_code: str | None
    country: str
    delivery_available: bool
    created_at: str
    updated_at: str
    published_at: str | None
    images: list[ItemImageModel] = Field(default_factory=list)


class CreateItemRequestModel(BaseModel):
    """Model for creating an item."""

    title: str = Field(..., min_length=3, max_length=200)
    description: str | None = None
    price_cents: int = Field(..., ge=0)
    category: str | None = None
    condition: Condition
    status: Status = Status.DRAFT
    is_featured: bool = False
    city: str | None = None
    postal_code: str | None = None
    country: str = "FR"
    delivery_available: bool = False
    images: list[ItemImageModel] = Field(default_factory=list)


class UpdateItemRequestModel(BaseModel):
    """Model for updating an item (PUT - full replace)."""

    title: str = Field(..., min_length=3, max_length=200)
    description: str | None = None
    price_cents: int = Field(..., ge=0)
    category: str | None = None
    condition: Condition
    status: Status
    is_featured: bool
    city: str | None = None
    postal_code: str | None = None
    country: str
    delivery_available: bool
    images: list[ItemImageModel] = Field(default_factory=list)


class PatchItemRequestModel(BaseModel):
    """Model for partially updating an item (PATCH)."""

    title: str | None = Field(None, min_length=3, max_length=200)
    description: str | None = None
    price_cents: int | None = Field(None, ge=0)
    category: str | None = None
    condition: Condition | None = None
    status: Status | None = None
    is_featured: bool | None = None
    city: str | None = None
    postal_code: str | None = None
    country: str | None = None
    delivery_available: bool | None = None
    images: list[ItemImageModel] | None = None


class ListItemsResponseModel(BaseModel):
    """Model for list items response."""

    items: list[ItemResponseModel]
    next_cursor: str | None


class ErrorDetailModel(BaseModel):
    """Model for error details."""

    code: str
    message: str
    details: dict[str, Any] | None = None


class ErrorResponseModel(BaseModel):
    """Model for error response."""

    error: ErrorDetailModel


class HealthResponseModel(BaseModel):
    """Model for health check response."""

    status: str
    timestamp: str


def to_response_model(item: Any) -> ItemResponseModel:
    """
    Convert domain entity to response model.

    Args:
        item: Item entity

    Returns:
        ItemResponseModel
    """
    from ...domain.entities.item import Item

    if isinstance(item, Item):
        images = [
            ItemImageModel(
                url=img.url,
                alt=img.alt,
                sort_order=img.sort_order,
            )
            for img in item.images
        ]
        return ItemResponseModel(
            id=item.id,
            title=item.title,
            description=item.description,
            price_cents=item.price_cents,
            category=item.category,
            condition=item.condition,
            status=item.status,
            is_featured=item.is_featured,
            city=item.city,
            postal_code=item.postal_code,
            country=item.country,
            delivery_available=item.delivery_available,
            created_at=item.created_at,
            updated_at=item.updated_at,
            published_at=item.published_at,
            images=images,
        )
    return item


def to_create_data(model: CreateItemRequestModel) -> Any:
    """
    Convert create request model to domain data.

    Args:
        model: CreateItemRequestModel

    Returns:
        CreateItemData
    """
    from ...domain.entities.item import CreateItemData
    from ...domain.value_objects.item_image import ItemImage

    images = (
        [
            ItemImage(
                url=img.url,
                alt=img.alt,
                sort_order=img.sort_order,
            )
            for img in model.images
        ]
        if model.images
        else None
    )

    return CreateItemData(
        title=model.title,
        description=model.description,
        price_cents=model.price_cents,
        category=model.category,
        condition=model.condition,
        status=model.status,
        is_featured=model.is_featured,
        city=model.city,
        postal_code=model.postal_code,
        country=model.country,
        delivery_available=model.delivery_available,
        images=images,
    )


def to_update_data(model: UpdateItemRequestModel) -> Any:
    """
    Convert update request model to domain data.

    Args:
        model: UpdateItemRequestModel

    Returns:
        UpdateItemData
    """
    from ...domain.entities.item import UpdateItemData
    from ...domain.value_objects.item_image import ItemImage

    images = (
        [
            ItemImage(
                url=img.url,
                alt=img.alt,
                sort_order=img.sort_order,
            )
            for img in model.images
        ]
        if model.images
        else None
    )

    return UpdateItemData(
        title=model.title,
        description=model.description,
        price_cents=model.price_cents,
        category=model.category,
        condition=model.condition,
        status=model.status,
        is_featured=model.is_featured,
        city=model.city,
        postal_code=model.postal_code,
        country=model.country,
        delivery_available=model.delivery_available,
        images=images,
    )


def to_patch_data(model: PatchItemRequestModel) -> Any:
    """
    Convert patch request model to domain data.

    Args:
        model: PatchItemRequestModel

    Returns:
        PatchItemData
    """
    from ...domain.entities.item import PatchItemData
    from ...domain.value_objects.item_image import ItemImage

    images = None
    if model.images is not None:
        images = [
            ItemImage(
                url=img.url,
                alt=img.alt,
                sort_order=img.sort_order,
            )
            for img in model.images
        ]

    return PatchItemData(
        title=model.title,
        description=model.description,
        price_cents=model.price_cents,
        category=model.category,
        condition=model.condition,
        status=model.status,
        is_featured=model.is_featured,
        city=model.city,
        postal_code=model.postal_code,
        country=model.country,
        delivery_available=model.delivery_available,
        images=images,
    )
