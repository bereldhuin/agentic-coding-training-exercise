package com.lebonpoint.integration.http

import com.lebonpoint.application.usecases.*
import com.lebonpoint.infrastructure.http.routes.configureItemRoutes
import com.lebonpoint.infrastructure.persistence.DatabaseConfig
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.Application
import io.ktor.server.application.install
import io.ktor.server.application.pluginOrNull
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.routing.routing
import io.ktor.server.testing.*
import kotlinx.coroutines.test.runTest
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.int
import kotlinx.serialization.json.jsonArray
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import org.junit.After
import org.junit.Before
import org.junit.Test
import org.koin.core.context.GlobalContext
import org.koin.core.context.loadKoinModules
import org.koin.core.context.startKoin
import org.koin.core.context.stopKoin
import org.koin.core.context.unloadKoinModules
import org.koin.core.module.Module
import kotlin.test.assertEquals
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

/**
 * Integration tests for API endpoints
 */
class ItemRoutesTest {

    private var loadedTestModule: Module? = null

    @Before
    fun setup() {
        loadedTestModule = null

        // Use the same database setup as repository tests
        kotlinx.coroutines.runBlocking {
            DatabaseConfig.shutdown()
            System.setProperty("DATABASE_PATH", ":memory:?cache=shared")
            DatabaseConfig.initializeSchema()

            // Clear any existing data
            DatabaseConfig.withConnection { conn ->
                conn.createStatement().use { stmt ->
                    stmt.execute("DELETE FROM items_fts")
                    stmt.execute("DELETE FROM items")
                }
            }
        }
    }

    @After
    fun teardown() {
        val hasContext = GlobalContext.getOrNull() != null
        if (hasContext) {
            loadedTestModule?.let { unloadKoinModules(it) }
            stopKoin()
        }
        loadedTestModule = null
        DatabaseConfig.shutdown()
    }

    @Test
    fun `should create item`() = runTest {
        testApplication {
            application {
                installTestModule()
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
            assertNotNull(json.jsonObject["id"]?.jsonPrimitive?.int)
            assertEquals("Test Item", json.jsonObject["title"]?.jsonPrimitive?.content)
        }
    }

    @Test
    fun `should get item by id`() = runTest {
        testApplication {
            application {
                installTestModule()
                configureTestRoutes()
            }

            // First create an item
            val createResponse = client.post("/v1/items") {
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
            val createdJson = Json.parseToJsonElement(createResponse.bodyAsText())
            val itemId = createdJson.jsonObject["id"]?.jsonPrimitive?.int

            // Then get it by ID
            val response = client.get("/v1/items/$itemId")
            assertEquals(HttpStatusCode.OK, response.status)

            val json = Json.parseToJsonElement(response.bodyAsText())
            assertEquals(itemId, json.jsonObject["id"]?.jsonPrimitive?.int)
            assertEquals("Test Item", json.jsonObject["title"]?.jsonPrimitive?.content)
        }
    }

    @Test
    fun `should return 404 for non-existent item`() = runTest {
        testApplication {
            application {
                installTestModule()
                configureTestRoutes()
            }

            val response = client.get("/v1/items/999")
            assertEquals(HttpStatusCode.NotFound, response.status)
        }
    }

    @Test
    fun `should list items`() = runTest {
        testApplication {
            application {
                installTestModule()
                configureTestRoutes()
            }

            // Create an item first
            client.post("/v1/items") {
                contentType(ContentType.Application.Json)
                setBody(
                    """
                    {
                        "title": "Item 1",
                        "price_cents": 10000,
                        "condition": "good"
                    }
                    """.trimIndent()
                )
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
        testApplication {
            application {
                installTestModule()
                configureTestRoutes()
            }

            // Create an item first
            val createResponse = client.post("/v1/items") {
                contentType(ContentType.Application.Json)
                setBody(
                    """
                    {
                        "title": "Item to Delete",
                        "price_cents": 10000,
                        "condition": "good"
                    }
                    """.trimIndent()
                )
            }
            val createdJson = Json.parseToJsonElement(createResponse.bodyAsText())
            val itemId = createdJson.jsonObject["id"]?.jsonPrimitive?.int

            // Delete it
            val response = client.delete("/v1/items/$itemId")
            assertEquals(HttpStatusCode.NoContent, response.status)
        }
    }

    private fun Application.installTestModule() {
        // Install content negotiation for JSON if not already installed
        if (pluginOrNull(ContentNegotiation) == null) {
            install(ContentNegotiation) {
                json()
            }
        }

        // Start Koin with app module (only if not already started)
        // Database schema is already initialized in @Before setup
        if (GlobalContext.getOrNull() == null) {
            startKoin {
                modules(com.lebonpoint.infrastructure.di.appModule)
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
}
