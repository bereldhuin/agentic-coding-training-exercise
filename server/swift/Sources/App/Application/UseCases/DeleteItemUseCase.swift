import Foundation
import Vapor

/// Use case for deleting an item
actor DeleteItemUseCase {
    private let repository: ItemRepositoryProtocol

    init(repository: ItemRepositoryProtocol) {
        self.repository = repository
    }

    /// Execute the use case
    func execute(id: Int) async throws -> Bool {
        // Check if item exists
        guard try await repository.exists(id: id) else {
            throw DomainError.notFound("Item not found")
        }

        // Delete item via repository
        return try await repository.delete(id: id)
    }
}
