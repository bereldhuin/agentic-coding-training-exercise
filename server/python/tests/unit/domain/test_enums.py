"""Tests for domain enums."""

from src.domain.entities.enums import Condition, Status


class TestCondition:
    """Tests for Condition enum."""

    def test_values(self) -> None:
        """Test that Condition.values() returns all values."""
        values = Condition.values()
        assert values == ["new", "like_new", "good", "fair", "parts", "unknown"]

    def test_enum_members(self) -> None:
        """Test that all enum members are accessible."""
        assert Condition.NEW == "new"
        assert Condition.LIKE_NEW == "like_new"
        assert Condition.GOOD == "good"
        assert Condition.FAIR == "fair"
        assert Condition.PARTS == "parts"
        assert Condition.UNKNOWN == "unknown"


class TestStatus:
    """Tests for Status enum."""

    def test_values(self) -> None:
        """Test that Status.values() returns all values."""
        values = Status.values()
        assert values == ["draft", "active", "reserved", "sold", "archived"]

    def test_enum_members(self) -> None:
        """Test that all enum members are accessible."""
        assert Status.DRAFT == "draft"
        assert Status.ACTIVE == "active"
        assert Status.RESERVED == "reserved"
        assert Status.SOLD == "sold"
        assert Status.ARCHIVED == "archived"
