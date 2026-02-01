"""Custom exceptions for HTTP layer."""

from typing import Any

from fastapi import HTTPException, Request, status
from fastapi.responses import JSONResponse

from ...shared.errors import LeBonPointError, ValidationError, NotFoundError


async def validation_error_handler(request: Request, exc: ValidationError) -> JSONResponse:
    """Handle ValidationError exceptions."""
    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={
            "error": {
                "code": "validation_error",
                "message": exc.message,
                "details": exc.details,
            }
        },
    )


async def not_found_error_handler(request: Request, exc: NotFoundError) -> JSONResponse:
    """Handle NotFoundError exceptions."""
    return JSONResponse(
        status_code=status.HTTP_404_NOT_FOUND,
        content={
            "error": {
                "code": "not_found",
                "message": exc.message,
                "details": None,
            }
        },
    )


async def generic_error_handler(request: Request, exc: Exception) -> JSONResponse:
    """Handle generic exceptions."""
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={
            "error": {
                "code": "internal_error",
                "message": "Internal server error",
                "details": None,
            }
        },
    )


async def pydantic_validation_error_handler(
    request: Request, exc: Any
) -> JSONResponse:
    """Handle Pydantic validation errors."""
    # Convert Pydantic ValidationError to our format
    details: dict[str, Any] = {}
    if hasattr(exc, "errors"):
        for error in exc.errors():
            field = ".".join(str(loc) for loc in error["loc"])
            details[field] = error["msg"]

    return JSONResponse(
        status_code=status.HTTP_400_BAD_REQUEST,
        content={
            "error": {
                "code": "validation_error",
                "message": "Validation failed",
                "details": details,
            }
        },
    )
