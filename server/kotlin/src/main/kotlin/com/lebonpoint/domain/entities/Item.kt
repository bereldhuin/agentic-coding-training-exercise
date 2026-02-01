package com.lebonpoint.domain.entities

import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemImage
import com.lebonpoint.domain.valueobjects.ItemStatus
import kotlinx.serialization.Serializable
import java.time.Instant

/**
 * Item entity
 * Represents a marketplace item with all its properties
 */
@Serializable
data class Item(
    val id: Int,
    val title: String,
    val description: String? = null,
    val priceCents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus,
    val isFeatured: Boolean,
    val city: String? = null,
    val postalCode: String? = null,
    val country: String,
    val deliveryAvailable: Boolean,
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    val createdAt: Instant = Instant.now(),
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    val updatedAt: Instant = Instant.now(),
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    val publishedAt: Instant? = null,
    val images: List<ItemImage> = emptyList()
)

/**
 * Item data for creation (without id and timestamps)
 */
@Serializable
data class CreateItemData(
    val title: String,
    val description: String? = null,
    val priceCents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus = ItemStatus.DRAFT,
    val isFeatured: Boolean = false,
    val city: String? = null,
    val postalCode: String? = null,
    val country: String = "FR",
    val deliveryAvailable: Boolean = false,
    val images: List<ItemImage> = emptyList()
)

/**
 * Item data for full replacement (PUT)
 */
@Serializable
data class ReplaceItemData(
    val title: String,
    val description: String? = null,
    val priceCents: Int,
    val category: String? = null,
    val condition: ItemCondition,
    val status: ItemStatus,
    val isFeatured: Boolean,
    val city: String? = null,
    val postalCode: String? = null,
    val country: String,
    val deliveryAvailable: Boolean,
    val images: List<ItemImage> = emptyList()
)

/**
 * Item data for partial update (PATCH)
 */
@Serializable
data class PatchItemData(
    val title: String? = null,
    val description: String? = null,
    val priceCents: Int? = null,
    val category: String? = null,
    val condition: ItemCondition? = null,
    val status: ItemStatus? = null,
    val isFeatured: Boolean? = null,
    val city: String? = null,
    val postalCode: String? = null,
    val country: String? = null,
    val deliveryAvailable: Boolean? = null,
    val images: List<ItemImage>? = null
)

/**
 * Apply patch data to an existing item
 */
fun Item.applyPatch(patch: PatchItemData): Item {
    return this.copy(
        title = patch.title ?: this.title,
        description = patch.description ?: this.description,
        priceCents = patch.priceCents ?: this.priceCents,
        category = patch.category ?: this.category,
        condition = patch.condition ?: this.condition,
        status = patch.status ?: this.status,
        isFeatured = patch.isFeatured ?: this.isFeatured,
        city = patch.city ?: this.city,
        postalCode = patch.postalCode ?: this.postalCode,
        country = patch.country ?: this.country,
        deliveryAvailable = patch.deliveryAvailable ?: this.deliveryAvailable,
        images = patch.images ?: this.images,
        updatedAt = Instant.now()
    )
}
