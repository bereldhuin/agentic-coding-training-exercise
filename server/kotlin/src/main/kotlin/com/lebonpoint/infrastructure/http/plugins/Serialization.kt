package com.lebonpoint.infrastructure.http.plugins

import io.ktor.serialization.kotlinx.json.*
import io.ktor.server.application.*
import io.ktor.server.plugins.contentnegotiation.*
import kotlinx.serialization.json.Json
import java.time.Instant

/**
 * Configure JSON content negotiation
 */
fun Application.configureSerialization() {
    install(ContentNegotiation) {
        json(Json {
            ignoreUnknownKeys = true
            encodeDefaults = true
            isLenient = true
            prettyPrint = false
            // Custom serializers for Instant are handled by kotlinx.serialization
        })
    }
}
