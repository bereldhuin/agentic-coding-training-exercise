import XCTest
@testable import App
import Vapor

/// Unit tests for domain entities and value objects
final class DomainTests: XCTestCase {
    // MARK: - ItemCondition Tests

    func testItemConditionAllCases() {
        XCTAssertEqual(ItemCondition.new.rawValue, "new")
        XCTAssertEqual(ItemCondition.likeNew.rawValue, "like_new")
        XCTAssertEqual(ItemCondition.good.rawValue, "good")
        XCTAssertEqual(ItemCondition.fair.rawValue, "fair")
        XCTAssertEqual(ItemCondition.parts.rawValue, "parts")
        XCTAssertEqual(ItemCondition.unknown.rawValue, "unknown")
    }

    func testItemConditionDecodable() throws {
        let json = """
        {
            "condition": "like_new"
        }
        """.data(using: .utf8)!

        let decoded = try JSONDecoder().decode(ItemCondition.self, from: json)
        XCTAssertEqual(decoded, .likeNew)
    }

    // MARK: - ItemStatus Tests

    func testItemStatusAllCases() {
        XCTAssertEqual(ItemStatus.draft.rawValue, "draft")
        XCTAssertEqual(ItemStatus.active.rawValue, "active")
        XCTAssertEqual(ItemStatus.reserved.rawValue, "reserved")
        XCTAssertEqual(ItemStatus.sold.rawValue, "sold")
        XCTAssertEqual(ItemStatus.archived.rawValue, "archived")
    }

    func testItemStatusDecodable() throws {
        let json = """
        {
            "status": "active"
        }
        """.data(using: .utf8)!

        let decoded = try JSONDecoder().decode(ItemStatus.self, from: json)
        XCTAssertEqual(decoded, .active)
    }

    // MARK: - ItemImage Tests

    func testItemImageEncoding() throws {
        let image = ItemImage(url: "https://example.com/image.jpg", alt: "Test image", sortOrder: 1)

        let encoder = JSONEncoder()
        let data = try encoder.encode(image)
        let string = String(data: data, encoding: .utf8)!

        XCTAssertTrue(string.contains("\"url\":\"https://example.com/image.jpg\""))
        XCTAssertTrue(string.contains("\"alt\":\"Test image\""))
        XCTAssertTrue(string.contains("\"sort_order\":1"))
    }

    func testItemImageDecoding() throws {
        let json = """
        {
            "url": "https://example.com/image.jpg",
            "alt": "Test image",
            "sort_order": 1
        }
        """.data(using: .utf8)!

        let decoded = try JSONDecoder().decode(ItemImage.self, from: json)

        XCTAssertEqual(decoded.url, "https://example.com/image.jpg")
        XCTAssertEqual(decoded.alt, "Test image")
        XCTAssertEqual(decoded.sortOrder, 1)
    }

    func testItemImageEquality() {
        let image1 = ItemImage(url: "https://example.com/image.jpg", alt: "Test", sortOrder: 1)
        let image2 = ItemImage(url: "https://example.com/image.jpg", alt: "Test", sortOrder: 1)
        let image3 = ItemImage(url: "https://example.com/other.jpg", alt: "Test", sortOrder: 1)

        XCTAssertEqual(image1, image2)
        XCTAssertNotEqual(image1, image3)
    }

    // MARK: - CursorCodec Tests

    func testCursorEncoding() {
        let cursor = CursorCodec.encode(id: 123, createdAt: "2026-01-31T12:00:00.000Z")

        // Cursor should be base64 encoded JSON
        XCTAssertFalse(cursor.isEmpty)

        // Should be decodable
        let decoded = CursorCodec.decode(cursor)
        XCTAssertNotNil(decoded)
        XCTAssertEqual(decoded?.id, 123)
        XCTAssertEqual(decoded?.createdAt, "2026-01-31T12:00:00.000Z")
    }

    func testCursorDecodingInvalid() {
        let invalidCursor = "invalid-base64!@#"
        let decoded = CursorCodec.decode(invalidCursor)
        XCTAssertNil(decoded)
    }

    // MARK: - ItemFilters Tests

    func testItemFiltersEmpty() {
        let filters = ItemFilters()

        XCTAssertNil(filters.status)
        XCTAssertNil(filters.category)
        XCTAssertNil(filters.minPriceCents)
        XCTAssertNil(filters.maxPriceCents)
        XCTAssertNil(filters.city)
        XCTAssertNil(filters.postalCode)
        XCTAssertNil(filters.isFeatured)
        XCTAssertNil(filters.deliveryAvailable)
    }

    func testItemFiltersWithValues() {
        let filters = ItemFilters(
            status: "active",
            category: "electronics",
            minPriceCents: 1000,
            maxPriceCents: 50000,
            isFeatured: true
        )

        XCTAssertEqual(filters.status, "active")
        XCTAssertEqual(filters.category, "electronics")
        XCTAssertEqual(filters.minPriceCents, 1000)
        XCTAssertEqual(filters.maxPriceCents, 50000)
        XCTAssertEqual(filters.isFeatured, true)
    }

    // MARK: - SortOptions Tests

    func testSortOptionsDefaults() {
        let sort = SortOptions()

        XCTAssertEqual(sort.field, .createdAt)
        XCTAssertEqual(sort.direction, .desc)
        XCTAssertEqual(sort.rawValue, "created_at:desc")
    }

    func testSortOptionsCustom() {
        let sort = SortOptions(field: .priceCents, direction: .asc)

        XCTAssertEqual(sort.field, .priceCents)
        XCTAssertEqual(sort.direction, .asc)
        XCTAssertEqual(sort.rawValue, "price_cents:asc")
    }

    func testSortOptionsFieldRawValue() {
        XCTAssertEqual(SortOptions.Field.id.rawValue, "id")
        XCTAssertEqual(SortOptions.Field.title.rawValue, "title")
        XCTAssertEqual(SortOptions.Field.priceCents.rawValue, "price_cents")
        XCTAssertEqual(SortOptions.Field.createdAt.rawValue, "created_at")
        XCTAssertEqual(SortOptions.Field.updatedAt.rawValue, "updated_at")
        XCTAssertEqual(SortOptions.Field.publishedAt.rawValue, "published_at")
    }

    // MARK: - DomainError Tests

    func testDomainErrorValidation() {
        let error = DomainError.validationError(["title": "Title too short"])

        XCTAssertEqual(error.status, .badRequest)

        let response = error.toErrorResponse()
        XCTAssertEqual(response.error.code, "validation_error")
        XCTAssertEqual(response.error.message, "Validation failed")
        XCTAssertEqual(response.error.details["title"], "Title too short")
    }

    func testDomainErrorNotFound() {
        let error = DomainError.notFound("Item not found")

        XCTAssertEqual(error.status, .notFound)

        let response = error.toErrorResponse()
        XCTAssertEqual(response.error.code, "not_found")
        XCTAssertEqual(response.error.message, "Item not found")
    }

    func testDomainErrorInternal() {
        let error = DomainError.internalError("Database connection failed")

        XCTAssertEqual(error.status, .internalServerError)

        let response = error.toErrorResponse()
        XCTAssertEqual(response.error.code, "internal_error")
        XCTAssertEqual(response.error.message, "Database connection failed")
    }
}
