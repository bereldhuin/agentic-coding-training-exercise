package com.lebonpoint.application.usecases

import com.lebonpoint.application.services.ValidationService
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.NotFoundException

/**
 * Use case for deleting an item
 */
class DeleteItemUseCase(
    private val repository: ItemRepository
) {
    /**
     * Execute the use case
     * @param id Item ID
     * @return Unit if deleted
     * @throws NotFoundException if item not found
     */
    suspend fun execute(id: Int) {
        ValidationService.validateId(id)

        val deleted = repository.delete(id)
        if (!deleted) {
            throw NotFoundException("Item", id)
        }
    }
}
