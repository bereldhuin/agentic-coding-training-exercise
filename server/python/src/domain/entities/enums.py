"""Domain enums for Item entity."""

from enum import Enum


class Condition(str, Enum):
    """Item condition enum."""

    NEW = "new"
    LIKE_NEW = "like_new"
    GOOD = "good"
    FAIR = "fair"
    PARTS = "parts"
    UNKNOWN = "unknown"

    @classmethod
    def values(cls) -> list[str]:
        """Get all condition values."""
        return [c.value for c in cls]


class Status(str, Enum):
    """Item status enum."""

    DRAFT = "draft"
    ACTIVE = "active"
    RESERVED = "reserved"
    SOLD = "sold"
    ARCHIVED = "archived"

    @classmethod
    def values(cls) -> list[str]:
        """Get all status values."""
        return [s.value for s in cls]
