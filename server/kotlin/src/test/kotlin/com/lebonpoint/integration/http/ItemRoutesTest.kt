package com.lebonpoint.integration.http

import com.lebonpoint.application.usecases.*
import com.lebonpoint.domain.entities.CreateItemData
import com.lebonpoint.domain.repositories.ItemRepository
import com.lebonpoint.domain.valueobjects.ItemCondition
import com.lebonpoint.domain.valueobjects.ItemStatus
import com.lebonpoint.infrastructure.di.appModule
import com.lebonpoint.infrastructure.http.routes.configureItemRoutes
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.server.application.Application
import io.ktor.server.routing.routing
import io.ktor.server.testing.*
import io.mockk.*
import kotlinx.coroutines.test.runTest
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.int
import kotlinx.serialization.json.jsonArray
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import org.junit.After
import org.junit.Before
import org.junit.Test
import org.koin.dsl.module
import org.koin.core.component.get
import org.koin.core.context.GlobalContext
import org.koin.core.context.startKoin
import org.koin.core.context.stopKoin
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

/**
 * Integration tests for API endpoints
 */
class ItemRoutesTest {

    private lateinit var mockRepository: ItemRepository

    @Before
    fun setup() {
        mockRepository = mockk()

        startKoin {
            modules(
                module {
                    single<ItemRepository> { mockRepository }
                    single { CreateItemUseCase(get()) }
                    single { GetItemUseCase(get()) }
                    single { ListItemsUseCase(get()) }
                    single { UpdateItemUseCase(get()) }
                    single { PatchItemUseCase(get()) }
                    single { DeleteItemUseCase(get()) }
                }
            )
        }
    }

    @After
    fun teardown() {
        stopKoin()
    }

    @Test
    fun `should return health check`() = runTest {
        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.get("/health")
            assertEquals(HttpStatusCode.OK, response.status)

            val json = Json.parseToJsonElement(response.bodyAsText())
            assertNotNull(json.jsonObject["status"])
            assertEquals("ok", json.jsonObject["status"]?.jsonPrimitive?.content)
        }
    }

    @Test
    fun `should create item`() = runTest {
        // Arrange
        coEvery { mockRepository.create(any()) } returns com.lebonpoint.domain.entities.Item(
            id = 1,
            title = "Test Item",
            description = null,
            priceCents = 10000,
            category = null,
            condition = ItemCondition.GOOD,
            status = ItemStatus.DRAFT,
            isFeatured = false,
            city = null,
            postalCode = null,
            country = "FR",
            deliveryAvailable = false,
            createdAt = java.time.Instant.EPOCH,
            updatedAt = java.time.Instant.EPOCH,
            images = emptyList()
        )

        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.post("/v1/items") {
                contentType(ContentType.Application.Json)
                setBody(
                    """
                    {
                        "title": "Test Item",
                        "price_cents": 10000,
                        "condition": "good"
                    }
                    """.trimIndent()
                )
            }

            assertEquals(HttpStatusCode.Created, response.status)

            val json = Json.parseToJsonElement(response.bodyAsText())
            assertEquals(1, json.jsonObject["id"]?.jsonPrimitive?.int)
            assertEquals("Test Item", json.jsonObject["title"]?.jsonPrimitive?.content)
        }
    }

    @Test
    fun `should get item by id`() = runTest {
        // Arrange
        coEvery { mockRepository.findById(1) } returns com.lebonpoint.domain.entities.Item(
            id = 1,
            title = "Test Item",
            description = null,
            priceCents = 10000,
            category = null,
            condition = ItemCondition.GOOD,
            status = ItemStatus.ACTIVE,
            isFeatured = false,
            city = null,
            postalCode = null,
            country = "FR",
            deliveryAvailable = false,
            createdAt = java.time.Instant.EPOCH,
            updatedAt = java.time.Instant.EPOCH,
            images = emptyList()
        )

        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.get("/v1/items/1")
            assertEquals(HttpStatusCode.OK, response.status)

            val json = Json.parseToJsonElement(response.bodyAsText())
            assertEquals(1, json.jsonObject["id"]?.jsonPrimitive?.int)
            assertEquals("Test Item", json.jsonObject["title"]?.jsonPrimitive?.content)
        }
    }

    @Test
    fun `should return 404 for non-existent item`() = runTest {
        // Arrange
        coEvery { mockRepository.findById(999) } returns null

        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.get("/v1/items/999")
            assertEquals(HttpStatusCode.NotFound, response.status)
        }
    }

    @Test
    fun `should list items`() = runTest {
        // Arrange
        coEvery {
            mockRepository.findAll(
                filters = any(),
                sort = any(),
                limit = any(),
                cursor = any()
            )
        } returns com.lebonpoint.domain.repositories.ItemPage(
            items = listOf(
                com.lebonpoint.domain.entities.Item(
                    id = 1,
                    title = "Item 1",
                    description = null,
                    priceCents = 10000,
                    category = null,
                    condition = ItemCondition.GOOD,
                    status = ItemStatus.ACTIVE,
                    isFeatured = false,
                    city = null,
                    postalCode = null,
                    country = "FR",
                    deliveryAvailable = false,
                    createdAt = java.time.Instant.EPOCH,
                    updatedAt = java.time.Instant.EPOCH,
                    images = emptyList()
                )
            ),
            nextCursor = null
        )

        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.get("/v1/items")
            assertEquals(HttpStatusCode.OK, response.status)

            val json = Json.parseToJsonElement(response.bodyAsText())
            assertNotNull(json.jsonObject["items"])
            assertTrue(json.jsonObject["items"]?.jsonArray?.isNotEmpty() == true)
        }
    }

    @Test
    fun `should delete item`() = runTest {
        // Arrange
        coEvery { mockRepository.delete(1) } returns true

        testApplication {
            application {
                configureTestRoutes()
            }

            val response = client.delete("/v1/items/1")
            assertEquals(HttpStatusCode.NoContent, response.status)
        }
    }
}

private fun Application.configureTestRoutes() {
    val koin = GlobalContext.get()
    routing {
        configureItemRoutes(
            createItemUseCase = koin.get(),
            getItemUseCase = koin.get(),
            listItemsUseCase = koin.get(),
            updateItemUseCase = koin.get(),
            patchItemUseCase = koin.get(),
            deleteItemUseCase = koin.get()
        )
    }
}
