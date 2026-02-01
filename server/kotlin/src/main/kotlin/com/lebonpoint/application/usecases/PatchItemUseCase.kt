package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.PatchItemData
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.NotFoundException

/**
 * Use case for partially updating an item
 */
class PatchItemUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param id Item ID
     * @param data Item patch data
     * @return Updated item
     * @throws NotFoundException if item not found
     * @throws ValidationException if validation fails
     */
    suspend fun execute(id: Int, data: PatchItemData): Item {
        ValidationService.validateId(id)

        // Validate provided fields
        data.title?.let { ValidationService.validateTitle(it) }
        data.priceCents?.let { ValidationService.validatePriceCents(it) }
        data.condition?.let { ValidationService.validateCondition(it) }
        data.status?.let { ValidationService.validateStatus(it) }

        return repository.patch(id, data)
            ?: throw NotFoundException("Item", id)
    }
}
