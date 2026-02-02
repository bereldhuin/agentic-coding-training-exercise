package com.lebonpoint.domain.repositories

import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.entities.ReplaceItemData

/**
 * Filter options for listing items
 */
data class ItemFilters(
    val status: String? = null,
    val category: String? = null,
    val minPriceCents: Int? = null,
    val maxPriceCents: Int? = null,
    val city: String? = null,
    val postalCode: String? = null,
    val isFeatured: Boolean? = null,
    val deliveryAvailable: Boolean? = null
)

/**
 * Sort options for listing items
 */
data class SortOptions(
    val field: String = "created_at",
    val direction: SortDirection = SortDirection.ASC
)

/**
 * Sort direction enum
 */
enum class SortDirection {
    ASC, DESC
}

/**
 * Paginated result page
 */
data class ItemPage(
    val items: List<Item>,
    val nextCursor: String? = null
)

/**
 * Item repository interface (port)
 * Defines the contract for item persistence operations
 */
interface ItemRepository {
    /**
     * Create a new item
     * @param data Item creation data
     * @return Created item with generated ID
     */
    suspend fun create(data: CreateItemData): Item

    /**
     * Find item by ID
     * @param id Item ID
     * @return Item if found, null otherwise
     */
    suspend fun findById(id: Int): Item?

    /**
     * List items with optional filters, sorting, and pagination
     * @param filters Filter options
     * @param sort Sort options
     * @param limit Maximum number of items to return (1-100)
     * @param cursor Pagination cursor (base64-encoded JSON)
     * @return Paginated result page
     */
    suspend fun findAll(
        filters: ItemFilters = ItemFilters(),
        sort: SortOptions = SortOptions(),
        limit: Int = 20,
        cursor: String? = null
    ): ItemPage

    /**
     * Full-text search for items
     * @param query Search query string
     * @param filters Filter options
     * @param sort Sort options
     * @param limit Maximum number of items to return (1-100)
     * @param cursor Pagination cursor (base64-encoded JSON)
     * @return Paginated result page
     */
    suspend fun search(
        query: String,
        filters: ItemFilters = ItemFilters(),
        sort: SortOptions = SortOptions(),
        limit: Int = 20,
        cursor: String? = null
    ): ItemPage

    /**
     * Update an item (full replacement)
     * @param id Item ID
     * @param data Item replacement data
     * @return Updated item if found, null otherwise
     */
    suspend fun update(id: Int, data: ReplaceItemData): Item?

    /**
     * Patch an item (partial update)
     * @param id Item ID
     * @param data Item patch data
     * @return Updated item if found, null otherwise
     */
    suspend fun patch(id: Int, data: PatchItemData): Item?

    /**
     * Delete an item
     * @param id Item ID
     * @return true if deleted, false if not found
     */
    suspend fun delete(id: Int): Boolean
}
