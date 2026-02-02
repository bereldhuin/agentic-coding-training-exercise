package com.lebonpoint

import com.lebonpoint.application.usecases.*
import com.lebonpoint.infrastructure.di.appModule
import com.lebonpoint.infrastructure.http.plugins.configureCORS
import com.lebonpoint.infrastructure.http.plugins.configureErrorHandling
import com.lebonpoint.infrastructure.http.plugins.configureSerialization
import com.lebonpoint.infrastructure.http.routes.configureItemRoutes
import com.lebonpoint.infrastructure.persistence.DatabaseConfig
import com.lebonpoint.shared.Configuration
import io.ktor.server.application.*
import io.ktor.server.engine.*
import io.ktor.server.netty.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.http.*
import org.koin.ktor.plugin.Koin
import org.koin.core.logger.Level
import org.koin.java.KoinJavaComponent.get
import java.nio.file.Path
import java.nio.file.Paths
import kotlin.system.exitProcess

/**
 * Main application entry point
 */
fun main() {
    val server = embeddedServer(Netty, port = Configuration.port, module = Application::module)

    // Add shutdown hook
    Runtime.getRuntime().addShutdownHook(Thread {
        println("Shutting down server...")
        DatabaseConfig.shutdown()
    })

    server.start(wait = true)
}

/**
 * Application module
 */
fun Application.module() {
    // Install Koin for dependency injection
    install(Koin) {
        Level.INFO
        allowOverride(true)
        modules(appModule)
    }

    // Configure plugins
    configureSerialization()
    configureErrorHandling()
    configureCORS()

    // Initialize database
    val dbConfig = DatabaseConfig
    // Initialize schema on startup
    // This will be called lazily on first database access

    val indexFile = resolveClientIndexPath()?.toFile()

    // Configure routes
    routing {
        get("/") {
            if (indexFile == null || !indexFile.exists()) {
                call.respond(HttpStatusCode.NotFound, "client/index.html not found")
                return@get
            }
            call.respondFile(indexFile)
        }
        configureItemRoutes(
            createItemUseCase = get(CreateItemUseCase::class.java),
            getItemUseCase = get(GetItemUseCase::class.java),
            listItemsUseCase = get(ListItemsUseCase::class.java),
            updateItemUseCase = get(UpdateItemUseCase::class.java),
            patchItemUseCase = get(PatchItemUseCase::class.java),
            deleteItemUseCase = get(DeleteItemUseCase::class.java)
        )
    }
}

private fun resolveClientIndexPath(): Path? {
    var dir = Paths.get(System.getProperty("user.dir"))
    repeat(10) {
        val candidate = dir.resolve("client").resolve("index.html")
        if (candidate.toFile().exists()) {
            return candidate
        }
        val parent = dir.parent ?: return null
        if (parent == dir) return null
        dir = parent
    }
    return null
}
