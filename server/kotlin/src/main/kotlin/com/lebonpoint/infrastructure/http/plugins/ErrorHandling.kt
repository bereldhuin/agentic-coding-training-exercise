package com.lebonpoint.infrastructure.http.plugins

import com.lebonpoint.shared.DomainException
import com.lebonpoint.shared.InternalException
import com.lebonpoint.shared.NotFoundException
import com.lebonpoint.shared.ValidationException
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.plugins.statuspages.*
import io.ktor.server.response.*
import com.lebonpoint.infrastructure.http.models.ErrorDetail
import com.lebonpoint.infrastructure.http.models.ErrorResponse

/**
 * Configure error handling with StatusPages
 */
fun Application.configureErrorHandling() {
    install(StatusPages) {
        // Validation exceptions
        exception<ValidationException> { call, cause ->
            call.respond(
                HttpStatusCode.BadRequest,
                ErrorResponse(
                    ErrorDetail(
                        code = cause.code,
                        message = cause.message ?: "Validation failed",
                        details = cause.details.mapValues { it.value?.toString() }
                    )
                )
            )
        }

        // Not found exceptions
        exception<NotFoundException> { call, cause ->
            call.respond(
                HttpStatusCode.NotFound,
                ErrorResponse(
                    ErrorDetail(
                        code = cause.code,
                        message = cause.message ?: "Resource not found",
                        details = cause.details.mapValues { it.value?.toString() }
                    )
                )
            )
        }

        // Internal exceptions
        exception<InternalException> { call, cause ->
            call.application.environment.log.error("Internal error", cause)
            call.respond(
                HttpStatusCode.InternalServerError,
                ErrorResponse(
                    ErrorDetail(
                        code = cause.code,
                        message = if (call.application.developmentMode) {
                            cause.message ?: "Internal server error"
                        } else {
                            "Internal server error"
                        },
                        details = emptyMap()
                    )
                )
            )
        }

        // Generic domain exceptions
        exception<DomainException> { call, cause ->
            call.respond(
                HttpStatusCode.InternalServerError,
                ErrorResponse(
                    ErrorDetail(
                        code = cause.code,
                        message = cause.message ?: "Domain error",
                        details = cause.details.mapValues { it.value?.toString() }
                    )
                )
            )
        }

        // Catch-all for unexpected exceptions
        exception<Throwable> { call, cause ->
            call.application.environment.log.error("Unexpected error", cause)
            call.respond(
                HttpStatusCode.InternalServerError,
                ErrorResponse(
                    ErrorDetail(
                        code = "internal_error",
                        message = if (call.application.developmentMode) {
                            cause.message ?: "Unexpected error occurred"
                        } else {
                            "Internal server error"
                        },
                        details = emptyMap()
                    )
                )
            )
        }
    }
}
