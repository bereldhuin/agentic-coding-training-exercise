package com.lebonpoint.infrastructure.di

import com.lebonpoint.application.usecases.*
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.infrastructure.persistence.DatabaseConfig
import com.lebonpoint.infrastructure.persistence.SQLiteItemRepository
import org.koin.dsl.module

/**
 * Koin dependency injection module
 */
val appModule = module {
    // Database configuration
    single { DatabaseConfig }

    // Repository
    single<ItemRepository> { SQLiteItemRepository(get()) }

    // Use cases
    single { CreateItemUseCase(get()) }
    single { GetItemUseCase(get()) }
    single { ListItemsUseCase(get()) }
    single { UpdateItemUseCase(get()) }
    single { PatchItemUseCase(get()) }
    single { DeleteItemUseCase(get()) }
}
