package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.ValidationException

/**
 * Use case for creating a new item
 */
class CreateItemUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param data Item creation data
     * @return Created item
     * @throws ValidationException if validation fails
     */
    suspend fun execute(data: CreateItemData): Item {
        // Validate input
        ValidationService.validateTitle(data.title)
        ValidationService.validatePriceCents(data.priceCents)
        ValidationService.validateCondition(data.condition)
        ValidationService.validateStatus(data.status)

        // Create item
        return repository.create(data)
    }
}
