import Foundation
import SQLiteKit

/// Database configuration and connection management
final class DatabaseConfiguration {
    let databasePath: String
    let sql: SQLDatabase
    private let connection: SQLiteConnection

    init(configuration: Configuration = .shared) throws {
        // Ensure database file exists
        Self.ensureDatabaseExists(at: configuration.databasePath)

        self.databasePath = configuration.databasePath
        let storage = SQLiteConnection.Storage.file(path: configuration.databasePath)
        self.connection = try SQLiteConnection.open(storage: storage).wait()
        self.sql = connection.sql()
    }

    /// Shutdown and cleanup resources
    func shutdown() async {
        do {
            try await connection.close()
        } catch {
            // Avoid surfacing shutdown errors as hard failures.
            print("Failed to shutdown database cleanly: \(error)")
        }
    }

    /// Ensure database file exists, create directory if needed
    private static func ensureDatabaseExists(at path: String) {
        let fileManager = FileManager.default
        let filePath = (path as NSString).expandingTildeInPath

        // Check if database file exists or can be created
        if !fileManager.fileExists(atPath: filePath) {
            let directory = (filePath as NSString).deletingLastPathComponent
            if !fileManager.fileExists(atPath: directory) {
                try? fileManager.createDirectory(
                    atPath: directory,
                    withIntermediateDirectories: true,
                    attributes: nil
                )
            }
        }
    }
}
