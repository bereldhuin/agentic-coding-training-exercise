"""Pytest configuration and fixtures."""

import sqlite3
import pytest
from typing import AsyncGenerator, Generator
from pathlib import Path

from fastapi.testclient import TestClient

from src.main import app
from src.domain.entities.item import CreateItemData, Item
from src.domain.entities.enums import Condition, Status
from src.domain.value_objects.item_image import ItemImage
from src.infrastructure.persistence.sqlite_item_repository import SQLiteItemRepository


@pytest.fixture
def in_memory_db() -> Generator[sqlite3.Connection, None, None]:
    """
    Create an in-memory SQLite database for testing.

    Yields:
        SQLite connection with initialized schema
    """
    conn = sqlite3.connect(":memory:")
    conn.row_factory = sqlite3.Row

    # Initialize schema
    conn.execute("""
        CREATE TABLE IF NOT EXISTS items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL CHECK(length(title) >= 3 AND length(title) <= 200),
            description TEXT,
            price_cents INTEGER NOT NULL CHECK(price_cents >= 0),
            category TEXT,
            condition TEXT NOT NULL CHECK(condition IN ('new', 'like_new', 'good', 'fair', 'parts', 'unknown')),
            status TEXT NOT NULL DEFAULT 'draft' CHECK(status IN ('draft', 'active', 'reserved', 'sold', 'archived')),
            is_featured INTEGER NOT NULL DEFAULT 0,
            city TEXT,
            postal_code TEXT,
            country TEXT NOT NULL DEFAULT 'FR',
            delivery_available INTEGER NOT NULL DEFAULT 0,
            created_at TEXT NOT NULL DEFAULT (datetime('now')),
            updated_at TEXT NOT NULL DEFAULT (datetime('now')),
            published_at TEXT,
            images TEXT NOT NULL DEFAULT '[]'
        )
    """)

    # Create indexes
    conn.execute("CREATE INDEX IF NOT EXISTS idx_items_status ON items(status)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_items_category ON items(category)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_items_condition ON items(condition)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_items_city ON items(city)")
    conn.execute("CREATE INDEX IF NOT EXISTS idx_items_price_cents ON items(price_cents)")

    # Create FTS5 virtual table
    conn.execute("""
        CREATE VIRTUAL TABLE IF NOT EXISTS items_fts USING fts5(
            title,
            description,
            content=items,
            content_rowid=rowid,
            tokenize='unicode61 remove_diacritics 1'
        )
    """)

    # Create triggers to keep FTS index in sync
    conn.execute("""
        CREATE TRIGGER IF NOT EXISTS items_fts_insert AFTER INSERT ON items BEGIN
            INSERT INTO items_fts(rowid, title, description) VALUES (new.id, new.title, new.description);
        END
    """)
    conn.execute("""
        CREATE TRIGGER IF NOT EXISTS items_fts_update AFTER UPDATE ON items BEGIN
            INSERT INTO items_fts(items_fts, rowid, title, description) VALUES('delete', old.id, old.title, old.description);
            INSERT INTO items_fts(rowid, title, description) VALUES (new.id, new.title, new.description);
        END
    """)
    conn.execute("""
        CREATE TRIGGER IF NOT EXISTS items_fts_delete AFTER DELETE ON items BEGIN
            INSERT INTO items_fts(items_fts, rowid, title, description) VALUES('delete', old.id, old.title, old.description);
        END
    """)

    yield conn

    conn.close()


@pytest.fixture
def test_repository(in_memory_db: sqlite3.Connection) -> SQLiteItemRepository:
    """
    Create a repository with in-memory database.

    Args:
        in_memory_db: In-memory database connection

    Returns:
        SQLiteItemRepository instance
    """
    # We need to patch the database path
    import tempfile
    import shutil

    # Create a temporary file-based database for the repository
    temp_dir = tempfile.mkdtemp()
    temp_db = Path(temp_dir) / "test.db"

    # Initialize the temp database with the same schema
    temp_conn = sqlite3.connect(str(temp_db))
    temp_conn.row_factory = sqlite3.Row

    # Copy schema, excluding auto-created internal objects:
    # - sqlite_sequence: internal AUTOINCREMENT tracking table
    # - FTS shadow tables (type='table'): auto-created when the FTS virtual table is created
    for row in in_memory_db.execute(
        "SELECT sql FROM sqlite_master WHERE sql IS NOT NULL "
        "AND name != 'sqlite_sequence' "
        "AND NOT (type = 'table' AND name LIKE 'items_fts_%')"
    ):
        temp_conn.execute(row[0])

    temp_conn.commit()
    temp_conn.close()

    repo = SQLiteItemRepository(db_path=temp_db)

    yield repo

    # Cleanup
    shutil.rmtree(temp_dir, ignore_errors=True)


@pytest.fixture
def sample_item_data() -> CreateItemData:
    """
    Create sample item data for testing.

    Returns:
        CreateItemData instance
    """
    return CreateItemData(
        title="Vintage Camera",
        description="A beautiful vintage camera from the 1970s",
        price_cents=15000,
        category="electronics",
        condition=Condition.GOOD,
        status=Status.ACTIVE,
        is_featured=True,
        city="Paris",
        postal_code="75001",
        country="FR",
        delivery_available=True,
        images=[
            ItemImage(url="https://example.com/camera1.jpg", alt="Front view", sort_order=0),
            ItemImage(url="https://example.com/camera2.jpg", alt="Side view", sort_order=1),
        ],
    )


@pytest.fixture
async def sample_item(test_repository: SQLiteItemRepository, sample_item_data: CreateItemData) -> Item:
    """
    Create a sample item in the database.

    Args:
        test_repository: Repository instance
        sample_item_data: Item data

    Returns:
        Created item
    """
    return await test_repository.create(sample_item_data)


@pytest.fixture
def client() -> Generator[TestClient, None, None]:
    """
    Create a FastAPI test client.

    Yields:
        TestClient instance
    """
    with TestClient(app) as test_client:
        yield test_client


# Async test support
@pytest.fixture
def anyio_backend() -> str:
    """Select the async backend for pytest-anyio."""
    return "asyncio"
