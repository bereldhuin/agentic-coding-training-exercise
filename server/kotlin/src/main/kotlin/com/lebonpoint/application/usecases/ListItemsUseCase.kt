package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.repositories.ItemFilters
import com.lebonpoint.domain.repositories.ItemPage
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.domain.repositories.SortDirection
import com.lebonpoint.domain.repositories.SortOptions
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.domain.valueobjects.toDatabaseString
import com.lebonpoint.shared.validationError

/**
 * Use case for listing items with filtering, sorting, and pagination
 */
class ListItemsUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param status Filter by status
     * @param category Filter by category
     * @param minPriceCents Filter by minimum price
     * @param maxPriceCents Filter by maximum price
     * @param city Filter by city
     * @param postalCode Filter by postal code
     * @param isFeatured Filter by featured status
     * @param deliveryAvailable Filter by delivery availability
     * @param sort Sort string (field:direction)
     * @param limit Maximum items per page
     * @param cursor Pagination cursor
     * @param search Search query (optional)
     * @return Paginated result page
     */
    suspend fun execute(
        status: String? = null,
        category: String? = null,
        minPriceCents: Int? = null,
        maxPriceCents: Int? = null,
        city: String? = null,
        postalCode: String? = null,
        isFeatured: Boolean? = null,
        deliveryAvailable: Boolean? = null,
        sort: String? = null,
        limit: Int? = null,
        cursor: String? = null,
        search: String? = null
    ): ItemPage {
        // Validate and parse limit
        val validatedLimit = ValidationService.validateLimit(limit)

        // Parse sort options
        val sortOptions = parseSortOptions(sort)

        // Build filters
        val normalizedStatus = status?.let {
            ItemStatus.fromString(it)?.toDatabaseString()
                ?: throw validationError("status", "Invalid status value")
        }
        val filters = ItemFilters(
            status = normalizedStatus,
            category = category,
            minPriceCents = minPriceCents,
            maxPriceCents = maxPriceCents,
            city = city,
            postalCode = postalCode,
            isFeatured = isFeatured,
            deliveryAvailable = deliveryAvailable
        )

        // Execute search or list
        return if (!search.isNullOrBlank()) {
            repository.search(search, filters, sortOptions, validatedLimit, cursor)
        } else {
            repository.findAll(filters, sortOptions, validatedLimit, cursor)
        }
    }

    /**
     * Parse sort string into SortOptions
     * Format: "field:direction" (e.g., "created_at:desc")
     */
    private fun parseSortOptions(sort: String?): SortOptions {
        if (sort.isNullOrBlank()) {
            return SortOptions(field = "created_at", direction = SortDirection.DESC)
        }

        val parts = sort.split(":")
        if (parts.size != 2) {
            throw validationError("sort", "sort must be in the form field:direction")
        }

        val field = ValidationService.validateSortField(parts.getOrNull(0))
        val direction = when (parts.getOrNull(1)?.lowercase()) {
            "asc" -> SortDirection.ASC
            "desc" -> SortDirection.DESC
            else -> throw validationError("sort", "sort direction must be asc or desc")
        }

        return SortOptions(field = field, direction = direction)
    }
}
