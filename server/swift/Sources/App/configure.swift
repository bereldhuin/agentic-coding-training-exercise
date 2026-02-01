import Vapor

/// Configure the Vapor application
public func configure(_ app: Application) throws {
    // Get configuration
    let config = Configuration.shared

    // Configure logger
    app.logger.logLevel = config.logLevel

    // Set up error middleware
    app.middleware.use(ErrorMiddleware())

    // Initialize database
    let databaseConfig = try DatabaseConfiguration(configuration: config)

    // Create repository
    let repository = SQLiteItemRepository(database: databaseConfig.sql)

    // Create use cases
    let createItemUseCase = CreateItemUseCase(repository: repository)
    let getItemUseCase = GetItemUseCase(repository: repository)
    let listItemsUseCase = ListItemsUseCase(repository: repository)
    let updateItemUseCase = UpdateItemUseCase(repository: repository)
    let patchItemUseCase = PatchItemUseCase(repository: repository)
    let deleteItemUseCase = DeleteItemUseCase(repository: repository)

    // Create controller
    let controller = ItemController(
        createItemUseCase: createItemUseCase,
        getItemUseCase: getItemUseCase,
        listItemsUseCase: listItemsUseCase,
        updateItemUseCase: updateItemUseCase,
        patchItemUseCase: patchItemUseCase,
        deleteItemUseCase: deleteItemUseCase
    )

    // Store controller in application
    app.itemController = controller

    // Configure JSON encoder/decoder
    app.routes.defaultMaxBodySize = "10mb"

    // Store database for cleanup
    app.storage[DatabaseKey.self] = databaseConfig

    // Register routes
    try routes(app)
}

/// Storage key for database configuration
struct DatabaseKey: StorageKey {
    typealias Value = DatabaseConfiguration
}

/// Extension to make database storable
extension Application {
    var databaseConfig: DatabaseConfiguration? {
        get {
            storage[DatabaseKey.self]
        }
        set {
            storage[DatabaseKey.self] = newValue
        }
    }
}
