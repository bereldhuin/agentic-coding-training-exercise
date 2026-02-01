import Foundation
import Vapor

/// Application configuration
struct Configuration: Sendable {
    let port: Int
    let databasePath: String
    let environment: Environment
    let logLevel: Logger.Level

    static let shared = Configuration()

    private init() {
        // Load from environment variables with defaults
        self.port = Int(Environment.get("PORT") ?? "9000") ?? 9000

        // Database path - relative to project or absolute
        let dbPathFromEnv = Environment.get("DATABASE_PATH")
        if let dbPath = dbPathFromEnv, !dbPath.isEmpty {
            // Check if it's a relative or absolute path
            if dbPath.hasPrefix("/") {
                self.databasePath = dbPath
            } else {
                // Resolve relative to the Swift server directory
                let currentDir = FileManager.default.currentDirectoryPath
                self.databasePath = "\(currentDir)/\(dbPath)"
            }
        } else {
            // Default to ../database/db.sqlite
            let currentDir = FileManager.default.currentDirectoryPath
            self.databasePath = "\(currentDir)/../database/db.sqlite"
        }

        // Environment detection
        let envString = Environment.get("ENVIRONMENT") ?? "development"
        switch envString.lowercased() {
        case "production":
            self.environment = .production
        case "testing":
            self.environment = .testing
        default:
            self.environment = .development
        }

        // Log level
        let logLevelString = Environment.get("LOG_LEVEL") ?? "info"
        switch logLevelString.lowercased() {
        case "debug":
            self.logLevel = .debug
        case "info":
            self.logLevel = .info
        case "warning":
            self.logLevel = .warning
        case "error":
            self.logLevel = .error
        case "critical":
            self.logLevel = .critical
        default:
            self.logLevel = .info
        }
    }
}
