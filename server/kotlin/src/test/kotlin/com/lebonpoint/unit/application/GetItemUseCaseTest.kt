package com.lebonpoint.unit.application

import com.lebonpoint.application.usecases.GetItemUseCase
import com.lebonpoint.domain.entities.Item
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.shared.NotFoundException
import com.lebonpoint.shared.ValidationException
import io.mockk.*
import kotlinx.coroutines.test.runTest
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

/**
 * Unit tests for GetItemUseCase
 */
class GetItemUseCaseTest {

    @Test
    fun `should get item by id`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = GetItemUseCase(repository)

        val expectedItem = Item(
            id = 1,
            title = "Test Item",
            description = null,
            priceCents = 10000,
            category = null,
            condition = com.lebonpoint.domain.valueobjects.ItemCondition.GOOD,
            status = com.lebonpoint.domain.valueobjects.ItemStatus.ACTIVE,
            isFeatured = false,
            city = null,
            postalCode = null,
            country = "FR",
            deliveryAvailable = false,
            createdAt = java.time.Instant.EPOCH,
            updatedAt = java.time.Instant.EPOCH,
            images = emptyList()
        )

        coEvery { repository.findById(1) } returns expectedItem

        // Act
        val result = useCase.execute(1)

        // Assert
        assertEquals(expectedItem, result)
        coVerify { repository.findById(1) }
    }

    @Test
    fun `should throw not found for non-existent item`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = GetItemUseCase(repository)

        coEvery { repository.findById(999) } returns null

        // Act & Assert
        assertFailsWith<NotFoundException> {
            useCase.execute(999)
        }
    }

    @Test
    fun `should throw validation error for invalid id`() = runTest {
        // Arrange
        val repository = mockk<ItemRepository>()
        val useCase = GetItemUseCase(repository)

        // Act & Assert
        assertFailsWith<ValidationException> {
            useCase.execute(0)
        }
    }
}
