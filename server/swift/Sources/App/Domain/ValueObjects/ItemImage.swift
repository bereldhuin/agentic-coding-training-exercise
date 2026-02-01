import Foundation

/// Item image value object
struct ItemImage: Codable, Equatable, Sendable {
    let url: String
    let alt: String?
    let sortOrder: Int?

    enum CodingKeys: String, CodingKey {
        case url
        case alt
        case sortOrder = "sort_order"
    }

    init(url: String, alt: String? = nil, sortOrder: Int? = nil) {
        self.url = url
        self.alt = alt
        self.sortOrder = sortOrder
    }
}
