package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.entities.ReplaceItemData
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.NotFoundException

/**
 * Use case for updating (replacing) an item
 */
class UpdateItemUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param id Item ID
     * @param data Item replacement data
     * @return Updated item
     * @throws NotFoundException if item not found
     * @throws ValidationException if validation fails
     */
    suspend fun execute(id: Int, data: ReplaceItemData): Item {
        ValidationService.validateId(id)
        ValidationService.validateTitle(data.title)
        ValidationService.validatePriceCents(data.priceCents)
        ValidationService.validateCondition(data.condition)
        ValidationService.validateStatus(data.status)

        return repository.update(id, data)
            ?: throw NotFoundException("Item", id)
    }
}
