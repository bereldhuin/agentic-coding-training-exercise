"""Cursor encoding and decoding utilities for pagination."""

import base64
import json
from dataclasses import dataclass


@dataclass(frozen=True)
class CursorData:
    """Data encoded in a pagination cursor."""

    id: int
    created_at: str


def encode_cursor(data: CursorData) -> str:
    """
    Encode cursor data to a base64 string.

    Args:
        data: The cursor data to encode

    Returns:
        Base64-encoded JSON string
    """
    json_data = {"id": data.id, "created_at": data.created_at}
    json_str = json.dumps(json_data, separators=(",", ":"))
    return base64.b64encode(json_str.encode()).decode()


def decode_cursor(cursor: str) -> CursorData | None:
    """
    Decode a base64 cursor string to cursor data.

    Args:
        cursor: Base64-encoded cursor string

    Returns:
        CursorData if valid, None if invalid
    """
    try:
        json_bytes = base64.b64decode(cursor)
        json_str = json_bytes.decode()
        data = json.loads(json_str)
        return CursorData(id=int(data["id"]), created_at=str(data["created_at"]))
    except (ValueError, KeyError, json.JSONDecodeError):
        return None
