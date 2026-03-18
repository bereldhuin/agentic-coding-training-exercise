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
    var id: Int,
    var title: String,
    var description: String? = null,
    var priceCents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus,
    var isFeatured: Boolean,
    var city: String? = null,
    var postalCode: String? = null,
    var country: String,
    var deliveryAvailable: Boolean,
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    var createdAt: Instant = Instant.now(),
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    var updatedAt: Instant = Instant.now(),
    @Serializable(with = com.lebonpoint.shared.serializers.InstantSerializer::class)
    var publishedAt: Instant? = null,
    var images: List<ItemImage> = emptyList()
)

/**
 * Item data for creation (without id and timestamps)
 */
@Serializable
data class CreateItemData(
    var title: String,
    var description: String? = null,
    var priceCents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus = ItemStatus.DRAFT,
    var isFeatured: Boolean = false,
    var city: String? = null,
    var postalCode: String? = null,
    var country: String = "FR",
    var deliveryAvailable: Boolean = false,
    var images: List<ItemImage> = emptyList()
)

/**
 * Item data for full replacement (PUT)
 */
@Serializable
data class ReplaceItemData(
    var title: String,
    var description: String? = null,
    var priceCents: Int,
    var category: String? = null,
    var condition: ItemCondition,
    var status: ItemStatus,
    var isFeatured: Boolean,
    var city: String? = null,
    var postalCode: String? = null,
    var country: String,
    var deliveryAvailable: Boolean,
    var images: List<ItemImage> = emptyList()
)

/**
 * Item data for partial update (PATCH)
 */
@Serializable
data class PatchItemData(
    var title: String? = null,
    var description: String? = null,
    var priceCents: Int? = null,
    var category: String? = null,
    var condition: ItemCondition? = null,
    var status: ItemStatus? = null,
    var isFeatured: Boolean? = null,
    var city: String? = null,
    var postalCode: String? = null,
    var country: String? = null,
    var deliveryAvailable: Boolean? = null,
    var images: List<ItemImage>? = null
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
