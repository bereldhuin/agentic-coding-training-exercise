"""Custom error classes for the application."""


class LeBonPointError(Exception):
    """Base exception for all application errors."""

    def __init__(self, message: str, details: dict | None = None) -> None:
        self.message = message
        self.details = details or {}
        super().__init__(message)


class ValidationError(LeBonPointError):
    """Validation error for invalid input data."""

    def __init__(self, message: str, details: dict | None = None) -> None:
        super().__init__(message, details)


class NotFoundError(LeBonPointError):
    """Resource not found error."""

    def __init__(self, message: str = "Resource not found") -> None:
        super().__init__(message)


class ConflictError(LeBonPointError):
    """Conflict error for duplicate or conflicting resources."""

    def __init__(self, message: str, details: dict | None = None) -> None:
        super().__init__(message, details)
