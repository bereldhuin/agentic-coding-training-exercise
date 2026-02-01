package com.lebonpoint.shared

/**
 * Application configuration
 * Loads configuration from environment variables with defaults
 */
object Configuration {
    /**
     * Server port
     */
    val port: Int by lazy {
        System.getenv("PORT")?.toIntOrNull() ?: 8080
    }

    /**
     * Database path
     * Relative to project root or absolute path
     */
    val databasePath: String by lazy {
        System.getenv("DATABASE_PATH") ?: "../database/db.sqlite"
    }

    /**
     * Application environment
     */
    val environment: Environment by lazy {
        when (System.getenv("ENV")?.lowercase()) {
            "production" -> Environment.PRODUCTION
            "testing" -> Environment.TESTING
            else -> Environment.DEVELOPMENT
        }
    }

    /**
     * Check if in development mode
     */
    val isDevelopment: Boolean
        get() = environment == Environment.DEVELOPMENT

    /**
     * Check if in testing mode
     */
    val isTesting: Boolean
        get() = environment == Environment.TESTING

    /**
     * Check if in production mode
     */
    val isProduction: Boolean
        get() = environment == Environment.PRODUCTION

    /**
     * Connection pool settings
     */
    val poolSize: Int by lazy {
        when (environment) {
            Environment.PRODUCTION -> 10
            Environment.TESTING -> 2
            Environment.DEVELOPMENT -> 5
        }
    }
}

/**
 * Application environment enum
 */
enum class Environment {
    DEVELOPMENT,
    TESTING,
    PRODUCTION
}
