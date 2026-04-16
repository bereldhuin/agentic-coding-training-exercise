"""Item entity and related data classes."""

from dataclasses import dataclass, field
from datetime import datetime

from .enums import Condition, Status
from ..value_objects.item_image import ItemImage


@dataclass(frozen=True)
class Item:
    """Item entity."""

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
    garantie_months: int | None
    created_at: str  # ISO 8601 datetime string
    updated_at: str  # ISO 8601 datetime string
    published_at: str | None
    images: list[ItemImage] = field(default_factory=list)


@dataclass
class CreateItemData:
    """Item data for creation (without id and timestamps)."""

    title: str
    price_cents: int
    condition: Condition
    description: str | None = None
    category: str | None = None
    status: Status = Status.DRAFT
    is_featured: bool = False
    city: str | None = None
    postal_code: str | None = None
    country: str = "FR"
    delivery_available: bool = False
    garantie_months: int | None = None
    images: list[ItemImage] | None = None


@dataclass
class UpdateItemData:
    """Item data for PUT (full replace)."""

    title: str
    price_cents: int
    condition: Condition
    status: Status
    is_featured: bool
    country: str
    delivery_available: bool
    description: str | None = None
    category: str | None = None
    city: str | None = None
    postal_code: str | None = None
    garantie_months: int | None = None
    images: list[ItemImage] | None = None


@dataclass
class PatchItemData:
    """Item data for PATCH (partial update). All fields optional."""

    title: str | None = None
    description: str | None = None
    price_cents: int | None = None
    category: str | None = None
    condition: Condition | None = None
    status: Status | None = None
    is_featured: bool | None = None
    city: str | None = None
    postal_code: str | None = None
    country: str | None = None
    delivery_available: bool | None = None
    garantie_months: int | None = None
    images: list[ItemImage] | None = None


def item_from_row(row: dict) -> Item:
    """
    Convert a database row to an Item entity.

    Args:
        row: Database row as dict

    Returns:
        Item entity
    """
    return Item(
        id=row["id"],
        title=row["title"],
        description=row.get("description"),
        price_cents=row["price_cents"],
        category=row.get("category"),
        condition=Condition(row["condition"]),
        status=Status(row["status"]),
        is_featured=bool(row["is_featured"]),
        city=row.get("city"),
        postal_code=row.get("postal_code"),
        country=row["country"],
        delivery_available=bool(row["delivery_available"]),
        garantie_months=row.get("garantie_months"),
        created_at=row["created_at"],
        updated_at=row["updated_at"],
        published_at=row.get("published_at"),
        images=parse_images_column(row.get("images")),
    )


def parse_images_column(value: str | None) -> list[ItemImage]:
    """
    Parse images from JSON column.

    Args:
        value: JSON string or None

    Returns:
        List of ItemImage objects
    """
    if not value:
        return []

    try:
        import json

        parsed = json.loads(value)
        if isinstance(parsed, list):
            return [
                ItemImage(
                    url=img.get("url", ""),
                    alt=img.get("alt"),
                    sort_order=img.get("sort_order"),
                )
                for img in parsed
                if isinstance(img, dict) and img.get("url")
            ]
    except (json.JSONDecodeError, AttributeError):
        pass

    return []


def serialize_images_column(images: list[ItemImage]) -> str:
    """
    Serialize images to JSON column.

    Args:
        images: List of ItemImage objects

    Returns:
        JSON string
    """
    import json

    data = [
        {"url": img.url, "alt": img.alt, "sort_order": img.sort_order}
        for img in images
    ]
    return json.dumps(data)
