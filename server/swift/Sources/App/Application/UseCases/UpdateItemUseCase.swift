import Foundation
import Vapor

/// Use case for updating (replacing) an item
actor UpdateItemUseCase {
    private let repository: ItemRepositoryProtocol

    init(repository: ItemRepositoryProtocol) {
        self.repository = repository
    }

    /// Execute the use case
    func execute(_ data: ReplaceItemData, id: Int) async throws -> Item {
        // Validate input
        try validate(data)

        // Check if item exists
        guard try await repository.exists(id: id) else {
            throw DomainError.notFound("Item not found")
        }

        // Update item via repository
        guard let item = try await repository.update(data, id: id) else {
            throw DomainError.notFound("Item not found")
        }

        return item
    }

    // MARK: - Validation

    private func validate(_ data: ReplaceItemData) throws {
        var errors: [String: String] = [:]

        // Title validation
        if data.title.count < 3 {
            errors["title"] = "Title must be at least 3 characters"
        } else if data.title.count > 200 {
            errors["title"] = "Title must be at most 200 characters"
        }

        // Price validation
        if data.priceCents < 0 {
            errors["price_cents"] = "Price must be greater than or equal to 0"
        }

        // If there are validation errors, throw
        if !errors.isEmpty {
            throw DomainError.validationError(errors)
        }
    }
}
