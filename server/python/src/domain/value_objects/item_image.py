"""Item image value object."""

from dataclasses import dataclass


@dataclass(frozen=True)
class ItemImage:
    """Item image value object."""

    url: str
    alt: str | None = None
    sort_order: int | None = None
