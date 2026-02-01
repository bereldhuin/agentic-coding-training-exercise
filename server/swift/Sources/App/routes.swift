import Foundation
import Vapor

/// Register all routes
func routes(_ app: Application) throws {
    let sourcePath = URL(fileURLWithPath: #file)
    let repoRoot = sourcePath
        .deletingLastPathComponent() // App
        .deletingLastPathComponent() // Sources
        .deletingLastPathComponent() // swift
        .deletingLastPathComponent() // server
        .deletingLastPathComponent() // repo
    let indexPath = repoRoot.appendingPathComponent("client/index.html").path

    // Get controller from application storage
    guard let controller = app.storage[ItemControllerKey.self] else {
        fatalError("ItemController not registered in application storage")
    }

    // Root page
    app.get { req -> Response in
        guard FileManager.default.fileExists(atPath: indexPath) else {
            return Response(status: .notFound)
        }
        let response = req.fileio.streamFile(at: indexPath)
        response.headers.contentType = .html
        return response
    }

    // Health check
    app.get("health", use: controller.health)

    // API v1 routes
    let v1 = app.grouped("v1")

    // Items routes
    let items = v1.grouped("items")
    items.get(use: controller.list)
    items.post(use: controller.create)

    // Item by ID routes
    items.get(":id", use: controller.get)
    items.put(":id", use: controller.update)
    items.patch(":id", use: controller.patch)
    items.delete(":id", use: controller.delete)
}

/// Storage key for ItemController
struct ItemControllerKey: StorageKey {
    typealias Value = ItemController
}

/// Extension to make ItemController storable
extension Application {
    var itemController: ItemController? {
        get {
            storage[ItemControllerKey.self]
        }
        set {
            storage[ItemControllerKey.self] = newValue
        }
    }
}
