"""Configuration management for the application."""

import os
from pathlib import Path


class Config:
    """Application configuration."""

    def __init__(self) -> None:
        self._port = int(os.getenv("PORT", "8000"))
        self._host = os.getenv("HOST", "0.0.0.0")
        self._database_path = os.getenv("DATABASE_PATH", "../database/db.sqlite")
        self._environment = os.getenv("ENVIRONMENT", "development")

    @property
    def port(self) -> int:
        """Get the server port."""
        return self._port

    @property
    def host(self) -> str:
        """Get the server host."""
        return self._host

    @property
    def database_path(self) -> Path:
        """
        Get the database path, resolved relative to the server directory.
        """
        path = Path(self._database_path)
        if not path.is_absolute():
            # Resolve relative to the server/python directory
            server_dir = Path(__file__).parent.parent.parent
            path = server_dir / self._database_path
        return path

    @property
    def environment(self) -> str:
        """Get the environment name."""
        return self._environment

    @property
    def is_development(self) -> bool:
        """Check if running in development mode."""
        return self._environment == "development"

    @property
    def is_production(self) -> bool:
        """Check if running in production mode."""
        return self._environment == "production"


# Global configuration instance
config = Config()
