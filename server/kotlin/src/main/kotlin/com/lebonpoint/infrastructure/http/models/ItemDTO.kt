package com.lebonpoint.infrastructure.http.models

import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.entities.ReplaceItemData
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemImage
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.shared.serializers.InstantSerializer
import kotlinx.serialization.Serializable
import java.time.Instant

/**
 * Health check response
 */
@Serializable
data class HealthResponse(
    var status: String,
    @Serializable(with = InstantSerializer::class)
    var timestamp: Instant = Instant.now()
)

/**
 * Item response for API
 */
@Serializable
data class ItemResponse(
    var id: Int,
    var title: String,
    var description: String? = null,
    var price_cents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus,
    var is_featured: Boolean,
    var city: String? = null,
    var postal_code: String? = null,
    var country: String,
    var delivery_available: Boolean,
    var created_at: String,
    var updated_at: String,
    var published_at: String? = null,
    var images: List<ItemImageResponse> = emptyList()
)

/**
 * Item image response
 */
@Serializable
data class ItemImageResponse(
    var url: String,
    var alt: String? = null,
    var sort_order: Int = 0
)

/**
 * List items response
 */
@Serializable
data class ListItemsResponse(
    var items: List<ItemResponse>,
    var next_cursor: String? = null
)

/**
 * Create item request
 */
@Serializable
data class CreateItemRequest(
    var title: String,
    var description: String? = null,
    var price_cents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus = ItemStatus.DRAFT,
    var is_featured: Boolean = false,
    var city: String? = null,
    var postal_code: String? = null,
    var country: String = "FR",
    var delivery_available: Boolean = false,
    var images: List<ItemImageRequest> = emptyList()
)

/**
 * Item image request
 */
@Serializable
data class ItemImageRequest(
    var url: String,
    var alt: String? = null,
    var sort_order: Int = 0
)

/**
 * Update item request (PUT - full replacement)
 */
@Serializable
data class UpdateItemRequest(
    var title: String,
    var description: String? = null,
    var price_cents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus,
    var is_featured: Boolean,
    var city: String? = null,
    var postal_code: String? = null,
    var country: String,
    var delivery_available: Boolean,
    var images: List<ItemImageRequest> = emptyList()
)

/**
 * Patch item request (PATCH - partial update)
 */
@Serializable
data class PatchItemRequest(
    var title: String? = null,
    var description: String? = null,
    var price_cents: Int? = null,
    var category: String? = null,
    var condition: ItemCondition? = null,
    var status: ItemStatus? = null,
    var is_featured: Boolean? = null,
    var city: String? = null,
    var postal_code: String? = null,
    var country: String? = null,
    var delivery_available: Boolean? = null,
    var images: List<ItemImageRequest>? = null
)

/**
 * Error response
 */
@Serializable
data class ErrorResponse(
    var error: ErrorDetail
)

/**
 * Error detail
 */
@Serializable
data class ErrorDetail(
    var code: String,
    var message: String,
    var details: Map<String, String?> = emptyMap()
)

/**
 * Convert Item entity to ItemResponse
 */
fun Item.toResponse(): ItemResponse {
    return ItemResponse(
        id = id,
        title = title,
        description = description,
        price_cents = priceCents,
        category = category,
        condition = condition,
        status = status,
        is_featured = isFeatured,
        city = city,
        postal_code = postalCode,
        country = country,
        delivery_available = deliveryAvailable,
        created_at = createdAt.toString(),
        updated_at = updatedAt.toString(),
        published_at = publishedAt?.toString(),
        images = images.map { it.toResponse() }
    )
}

/**
 * Convert ItemImage to ItemImageResponse
 */
fun ItemImage.toResponse(): ItemImageResponse {
    return ItemImageResponse(
        url = url,
        alt = alt,
        sort_order = sortOrder
    )
}

/**
 * Convert CreateItemRequest to CreateItemData
 */
fun CreateItemRequest.toData(): CreateItemData {
    return CreateItemData(
        title = title,
        description = description,
        priceCents = price_cents,
        category = category,
        condition = condition,
        status = status,
        isFeatured = is_featured,
        city = city,
        postalCode = postal_code,
        country = country,
        deliveryAvailable = delivery_available,
        images = images.map { it.toDomain() }
    )
}

/**
 * Convert UpdateItemRequest to ReplaceItemData
 */
fun UpdateItemRequest.toData(): ReplaceItemData {
    return ReplaceItemData(
        title = title,
        description = description,
        priceCents = price_cents,
        category = category,
        condition = condition,
        status = status,
        isFeatured = is_featured,
        city = city,
        postalCode = postal_code,
        country = country,
        deliveryAvailable = delivery_available,
        images = images.map { it.toDomain() }
    )
}

/**
 * Convert PatchItemRequest to PatchItemData
 */
fun PatchItemRequest.toData(): PatchItemData {
    return PatchItemData(
        title = title,
        description = description,
        priceCents = price_cents,
        category = category,
        condition = condition,
        status = status,
        isFeatured = is_featured,
        city = city,
        postalCode = postal_code,
        country = country,
        deliveryAvailable = delivery_available,
        images = images?.map { it.toDomain() }
    )
}

/**
 * Convert ItemImageRequest to ItemImage
 */
fun ItemImageRequest.toDomain(): ItemImage {
    return ItemImage(
        url = url,
        alt = alt,
        sortOrder = sort_order
    )
}
