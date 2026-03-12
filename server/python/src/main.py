"""Main FastAPI application entry point."""

from contextlib import asynccontextmanager
from datetime import datetime, timezone
from pathlib import Path

from fastapi import FastAPI, status, HTTPException
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import FileResponse

from .shared.config import config
from .shared.errors import ValidationError, NotFoundError
from .infrastructure.http.dependencies import get_repository
from .infrastructure.http.routes import create_items_router
from .infrastructure.http.exceptions import (
    validation_error_handler,
    not_found_error_handler,
    generic_error_handler,
    pydantic_validation_error_handler,
)
from .infrastructure.http.models import HealthResponseModel

from pydantic import ValidationError as PydanticValidationError


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    # Startup
    print(f"Starting LeBonPoint Python server on port {config.port}")
    print(f"Database path: {config.database_path}")
    print(f"Environment: {config.environment}")
    yield
    # Shutdown
    print("Shutting down LeBonPoint Python server")


# Create FastAPI application
app = FastAPI(
    title="LeBonPoint API",
    description="Python FastAPI implementation of LeBonPoint marketplace items API",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    lifespan=lifespan,
)

# CORS middleware (allow all for development)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Register exception handlers
app.add_exception_handler(ValidationError, validation_error_handler)
app.add_exception_handler(NotFoundError, not_found_error_handler)
app.add_exception_handler(RequestValidationError, pydantic_validation_error_handler)
app.add_exception_handler(PydanticValidationError, pydantic_validation_error_handler)
app.add_exception_handler(Exception, generic_error_handler)

# Include routers - register immediately for tests
repository = get_repository()
items_router = create_items_router(repository)
app.include_router(items_router)


# Health check endpoint
@app.get(
    "/health",
    response_model=HealthResponseModel,
    summary="Health check",
    tags=["health"],
)
async def health_check() -> dict[str, str]:
    """
    Health check endpoint.

    Returns the current status and timestamp.
    """
    return {
        "status": "ok",
        "timestamp": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
    }


# Root endpoint
@app.get("/", summary="Root endpoint", tags=["root"])
async def root() -> FileResponse:
    """
    Root endpoint.

    Returns the shared static HTML page.
    """
    repo_root = Path(__file__).resolve().parents[3]
    index_path = repo_root / "client" / "index.html"
    if not index_path.exists():
        raise HTTPException(status_code=500, detail="client/index.html not found")
    return FileResponse(index_path, media_type="text/html")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("main:app", port=config.port, reload=config.is_development)
