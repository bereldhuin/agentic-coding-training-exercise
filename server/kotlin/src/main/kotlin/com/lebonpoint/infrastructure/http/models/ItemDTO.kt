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
    val status: String,
    @Serializable(with = InstantSerializer::class)
    val timestamp: Instant = Instant.now()
)

/**
 * Item response for API
 */
@Serializable
data class ItemResponse(
    val id: Int,
    val title: String,
    val description: String? = null,
    val price_cents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus,
    val is_featured: Boolean,
    val city: String? = null,
    val postal_code: String? = null,
    val country: String,
    val delivery_available: Boolean,
    val created_at: String,
    val updated_at: String,
    val published_at: String? = null,
    val images: List<ItemImageResponse> = emptyList()
)

/**
 * Item image response
 */
@Serializable
data class ItemImageResponse(
    val url: String,
    val alt: String? = null,
    val sort_order: Int = 0
)

/**
 * List items response
 */
@Serializable
data class ListItemsResponse(
    val items: List<ItemResponse>,
    val next_cursor: String? = null
)

/**
 * Create item request
 */
@Serializable
data class CreateItemRequest(
    val title: String,
    val description: String? = null,
    val price_cents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus = ItemStatus.DRAFT,
    val is_featured: Boolean = false,
    val city: String? = null,
    val postal_code: String? = null,
    val country: String = "FR",
    val delivery_available: Boolean = false,
    val images: List<ItemImageRequest> = emptyList()
)

/**
 * Item image request
 */
@Serializable
data class ItemImageRequest(
    val url: String,
    val alt: String? = null,
    val sort_order: Int = 0
)

/**
 * Update item request (PUT - full replacement)
 */
@Serializable
data class UpdateItemRequest(
    val title: String,
    val description: String? = null,
    val price_cents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus,
    val is_featured: Boolean,
    val city: String? = null,
    val postal_code: String? = null,
    val country: String,
    val delivery_available: Boolean,
    val images: List<ItemImageRequest> = emptyList()
)

/**
 * Patch item request (PATCH - partial update)
 */
@Serializable
data class PatchItemRequest(
    val title: String? = null,
    val description: String? = null,
    val price_cents: Int? = null,
    val category: String? = null,
    val condition: ItemCondition? = null,
    val status: ItemStatus? = null,
    val is_featured: Boolean? = null,
    val city: String? = null,
    val postal_code: String? = null,
    val country: String? = null,
    val delivery_available: Boolean? = null,
    val images: List<ItemImageRequest>? = null
)

/**
 * Error response
 */
@Serializable
data class ErrorResponse(
    val error: ErrorDetail
)

/**
 * Error detail
 */
@Serializable
data class ErrorDetail(
    val code: String,
    val message: String,
    val details: Map<String, String?> = emptyMap()
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
