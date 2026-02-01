package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.NotFoundException

/**
 * Use case for getting an item by ID
 */
class GetItemUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param id Item ID
     * @return Item if found
     * @throws NotFoundException if item not found
     */
    suspend fun execute(id: Int): Item {
        ValidationService.validateId(id)

        return repository.findById(id)
            ?: throw NotFoundException("Item", id)
    }
}
