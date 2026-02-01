import Foundation
import Vapor

/// Use case for listing items with filters, sorting, and pagination
actor ListItemsUseCase {
    private let repository: ItemRepositoryProtocol

    init(repository: ItemRepositoryProtocol) {
        self.repository = repository
    }

    /// Execute the use case for listing items
    func execute(
        filters: ItemFilters,
        sort: SortOptions,
        limit: Int,
        cursor: String?
    ) async throws -> ItemPage {
        // Validate limit
        let validatedLimit = min(max(1, limit), 100)

        // Get items from repository
        return try await repository.findAll(
            filters: filters,
            sort: sort,
            limit: validatedLimit,
            cursor: cursor
        )
    }

    /// Execute the use case for searching items
    func executeSearch(
        query: String,
        filters: ItemFilters,
        sort: SortOptions,
        limit: Int,
        cursor: String?
    ) async throws -> ItemPage {
        // Validate query
        guard !query.trimmingCharacters(in: .whitespaces).isEmpty else {
            throw DomainError.validationError(["q": "Search query cannot be empty"])
        }

        // Validate limit
        let validatedLimit = min(max(1, limit), 100)

        // Search items via repository
        return try await repository.search(
            query: query,
            filters: filters,
            sort: sort,
            limit: validatedLimit,
            cursor: cursor
        )
    }
}
