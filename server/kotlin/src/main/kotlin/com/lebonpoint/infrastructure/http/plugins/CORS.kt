package com.lebonpoint.infrastructure.http.plugins

import io.ktor.server.application.*
import io.ktor.server.plugins.cors.routing.*

/**
 * Configure CORS
 */
fun Application.configureCORS() {
    install(CORS) {
        // Allow any host in development
        anyHost()

        // Allow common methods
        allowMethod(io.ktor.http.HttpMethod.Get)
        allowMethod(io.ktor.http.HttpMethod.Post)
        allowMethod(io.ktor.http.HttpMethod.Put)
        allowMethod(io.ktor.http.HttpMethod.Patch)
        allowMethod(io.ktor.http.HttpMethod.Delete)
        allowMethod(io.ktor.http.HttpMethod.Options)

        // Allow headers
        allowHeader(io.ktor.http.HttpHeaders.ContentType)
        allowHeader(io.ktor.http.HttpHeaders.Authorization)

        // Expose headers
        exposeHeader(io.ktor.http.HttpHeaders.ContentType)
    }
}
