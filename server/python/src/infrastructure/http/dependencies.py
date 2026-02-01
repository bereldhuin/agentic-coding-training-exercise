"""Dependency injection setup for FastAPI."""

from functools import lru_cache
from typing import AsyncGenerator

from fastapi import Depends

from ...domain.repositories.item_repository import ItemRepository
from ...infrastructure.persistence.sqlite_item_repository import SQLiteItemRepository
from ...shared.config import config


@lru_cache
def get_repository() -> ItemRepository:
    """
    Get the item repository instance.

    Uses lru_cache to ensure singleton pattern.

    Returns:
        ItemRepository instance
    """
    return SQLiteItemRepository()


async def get_db_connection():
    """
    Get a database connection for the current request.

    This is a placeholder for future connection pooling if needed.
    Currently, the repository creates connections per operation.

    Yields:
        Database connection
    """
    # For now, this is a placeholder
    # The repository manages its own connections
    yield None
