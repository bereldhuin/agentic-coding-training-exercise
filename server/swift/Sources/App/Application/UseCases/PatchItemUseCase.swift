import Foundation
import Vapor

/// Use case for partially updating an item
actor PatchItemUseCase {
    private let repository: ItemRepositoryProtocol

    init(repository: ItemRepositoryProtocol) {
        self.repository = repository
    }

    /// Execute the use case
    func execute(_ data: UpdateItemData, id: Int) async throws -> Item {
        // Validate input (only validate provided fields)
        try validate(data)

        // Check if item exists
        guard try await repository.exists(id: id) else {
            throw DomainError.notFound("Item not found")
        }

        // Patch item via repository
        guard let item = try await repository.patch(data, id: id) else {
            throw DomainError.notFound("Item not found")
        }

        return item
    }

    // MARK: - Validation

    private func validate(_ data: UpdateItemData) throws {
        var errors: [String: String] = [:]

        // Title validation (only if provided)
        if let title = data.title {
            if title.count < 3 {
                errors["title"] = "Title must be at least 3 characters"
            } else if title.count > 200 {
                errors["title"] = "Title must be at most 200 characters"
            }
        }

        // Price validation (only if provided)
        if let priceCents = data.priceCents {
            if priceCents < 0 {
                errors["price_cents"] = "Price must be greater than or equal to 0"
            }
        }

        // If there are validation errors, throw
        if !errors.isEmpty {
            throw DomainError.validationError(errors)
        }
    }
}
