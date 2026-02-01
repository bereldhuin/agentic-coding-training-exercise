import Foundation
import Vapor

/// Use case for getting a single item by ID
actor GetItemUseCase {
    private let repository: ItemRepositoryProtocol

    init(repository: ItemRepositoryProtocol) {
        self.repository = repository
    }

    /// Execute the use case
    func execute(id: Int) async throws -> Item {
        guard let item = try await repository.findById(id) else {
            throw DomainError.notFound("Item not found")
        }
        return item
    }
}
