import Vapor
import Logging

@main
struct App {
    static func main() async throws {
        var env = try Environment.detect()
        try LoggingSystem.bootstrap(from: &env)

        let app = try await Application.make(env)
        var exitCode = 0

        do {
            let config = Configuration.shared

            app.http.server.configuration.port = config.port

            try configure(app)

            app.logger.info("Server starting on port \(config.port)")
            app.logger.info("Database path: \(config.databasePath)")

            try await app.execute()
        } catch {
            app.logger.error("Failed to start server: \(error)")
            exitCode = 1
        }

        if let dbConfig = app.databaseConfig {
            await dbConfig.shutdown()
        }

        do {
            try await app.asyncShutdown()
        } catch {
            app.logger.warning("Failed to shutdown cleanly: \(error)")
        }

        if exitCode != 0 {
            exit(Int32(exitCode))
        }
    }
}
