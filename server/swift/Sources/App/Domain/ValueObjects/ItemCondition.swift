import Foundation

/// Item condition enum
enum ItemCondition: String, Codable, CaseIterable, Sendable {
    case new = "new"
    case likeNew = "like_new"
    case good = "good"
    case fair = "fair"
    case parts = "parts"
    case unknown = "unknown"
}

/// Item status enum
enum ItemStatus: String, Codable, CaseIterable, Sendable {
    case draft = "draft"
    case active = "active"
    case reserved = "reserved"
    case sold = "sold"
    case archived = "archived"
}
