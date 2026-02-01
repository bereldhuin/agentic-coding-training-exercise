package com.lebonpoint.application.services

import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.shared.validationError

/**
 * Validation service for item data
 * Provides validation logic for item creation and updates
 */
object ValidationService {

    /**
     * Validate title
     */
    fun validateTitle(title: String?) {
        when {
            title.isNullOrBlank() -> throw validationError("title", "Title is required")
            title.length < 3 -> throw validationError("title", "Title must be at least 3 characters")
            title.length > 200 -> throw validationError("title", "Title must be at most 200 characters")
        }
    }

    /**
     * Validate price in cents
     */
    fun validatePriceCents(priceCents: Int?) {
        when {
            priceCents == null -> throw validationError("priceCents", "Price is required")
            priceCents < 0 -> throw validationError("priceCents", "Price must be greater than or equal to 0")
        }
    }

    /**
     * Validate condition
     */
    fun validateCondition(condition: ItemCondition?) {
        if (condition == null) {
            throw validationError("condition", "Condition is required")
        }
    }

    /**
     * Validate status (optional, defaults to DRAFT)
     */
    fun validateStatus(status: ItemStatus?) {
        // Status is optional, defaults to DRAFT
    }

    /**
     * Validate limit for pagination
     */
    fun validateLimit(limit: Int?): Int {
        return when {
            limit == null -> 20
            limit < 1 || limit > 100 -> throw validationError("limit", "limit must be between 1 and 100")
            else -> limit
        }
    }

    /**
     * Validate ID parameter
     */
    fun validateId(id: Int?) {
        when {
            id == null -> throw validationError("id", "ID is required")
            id < 1 -> throw validationError("id", "ID must be a positive integer")
        }
    }

    /**
     * Validate sort field
     */
    fun validateSortField(field: String?): String {
        val validFields = setOf("id", "title", "price_cents", "created_at", "updated_at", "published_at")
        return if (field != null && field in validFields) field else throw validationError(
            "sort",
            "sort field must be one of: id, title, price_cents, created_at, updated_at, published_at"
        )
    }
}
