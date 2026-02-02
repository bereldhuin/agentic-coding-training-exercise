package com.lebonpoint.infrastructure.http.routes

import com.lebonpoint.application.usecases.*
import com.lebonpoint.infrastructure.http.models.*
import io.ktor.server.application.*
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.time.Instant

/**
 * Configure API routes
 */
fun Routing.configureItemRoutes(
    createItemUseCase: CreateItemUseCase,
    getItemUseCase: GetItemUseCase,
    listItemsUseCase: ListItemsUseCase,
    updateItemUseCase: UpdateItemUseCase,
    patchItemUseCase: PatchItemUseCase,
    deleteItemUseCase: DeleteItemUseCase
) {
    // Health check endpoint
    get("/health") {
        call.respond(
            HealthResponse(
                status = "ok",
                timestamp = Instant.now()
            )
        )
    }

    // Items routes
    route("/v1/items") {
        // List items with filtering, sorting, and pagination
        get {
            val status = call.request.queryParameters["status"]
            val category = call.request.queryParameters["category"]
            val city = call.request.queryParameters["city"]
            val postalCode = call.request.queryParameters["postal_code"]
            val isFeatured = call.request.queryParameters["is_featured"]?.toBooleanStrictOrNull()
            val deliveryAvailable = call.request.queryParameters["delivery_available"]?.toBooleanStrictOrNull()
            val sort = call.request.queryParameters["sort"]
            val limit = call.request.queryParameters["limit"]?.toIntOrNull()
            val cursor = call.request.queryParameters["cursor"]
            val search = call.request.queryParameters["q"]

            val result = listItemsUseCase.execute(
                status = status,
                category = category,
                city = city,
                postalCode = postalCode,
                isFeatured = isFeatured,
                deliveryAvailable = deliveryAvailable,
                sort = sort,
                limit = limit,
                cursor = cursor,
                search = search
            )

            call.respond(
                ListItemsResponse(
                    items = result.items.map { it.toResponse() },
                    next_cursor = result.nextCursor
                )
            )
        }

        // Create a new item
        post {
            val request = call.receive<CreateItemRequest>()
            val item = createItemUseCase.execute(request.toData())
            call.respond(HttpStatusCode.Created, item.toResponse())
        }

        // Get, update, patch, delete specific item
        route("/{id}") {
            // Get item by ID
            get {
                val id = call.parameters["id"]?.toIntOrNull()
                    ?: return@get call.respond(
                        HttpStatusCode.BadRequest,
                        ErrorResponse(ErrorDetail("validation_error", "Invalid ID", mapOf("id" to "ID must be an integer")))
                    )

                val item = getItemUseCase.execute(id)
                call.respond(item.toResponse())
            }

            // Update item (full replacement)
            put {
                val id = call.parameters["id"]?.toIntOrNull()
                    ?: return@put call.respond(
                        HttpStatusCode.BadRequest,
                        ErrorResponse(ErrorDetail("validation_error", "Invalid ID", mapOf("id" to "ID must be an integer")))
                    )

                val request = call.receive<UpdateItemRequest>()
                val item = updateItemUseCase.execute(id, request.toData())
                call.respond(item.toResponse())
            }

            // Patch item (partial update)
            patch {
                val id = call.parameters["id"]?.toIntOrNull()
                    ?: return@patch call.respond(
                        HttpStatusCode.BadRequest,
                        ErrorResponse(ErrorDetail("validation_error", "Invalid ID", mapOf("id" to "ID must be an integer")))
                    )

                val request = call.receive<PatchItemRequest>()
                val item = patchItemUseCase.execute(id, request.toData())
                call.respond(item.toResponse())
            }

            // Delete item
            delete {
                val id = call.parameters["id"]?.toIntOrNull()
                    ?: return@delete call.respond(
                        HttpStatusCode.BadRequest,
                        ErrorResponse(ErrorDetail("validation_error", "Invalid ID", mapOf("id" to "ID must be an integer")))
                    )

                deleteItemUseCase.execute(id)
                call.respond(HttpStatusCode.NoContent)
            }
        }
    }
}
