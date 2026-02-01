import Foundation
import Vapor

/// Domain error enum
enum DomainError: Error, Sendable {
    case validationError([String: String])
    case notFound(String)
    case internalError(String)
}

/// Validation error detail
struct ValidationErrorDetail: Codable {
    let field: String
    let message: String
}

/// Error response matching canonical format
struct ErrorResponse: Content {
    let error: ErrorContent

    struct ErrorContent: Content {
        let code: String
        let message: String
        let details: [String: String]
    }
}

extension DomainError: AbortError {
    var status: HTTPResponseStatus {
        switch self {
        case .validationError:
            return .badRequest
        case .notFound:
            return .notFound
        case .internalError:
            return .internalServerError
        }
    }

    var reason: String {
        switch self {
        case .validationError(let details):
            return "Validation failed"
        case .notFound(let message):
            return message
        case .internalError(let message):
            return message
        }
    }
}

/// Convert DomainError to ErrorResponse
extension DomainError {
    func toErrorResponse() -> ErrorResponse {
        switch self {
        case .validationError(let details):
            return ErrorResponse(error: .init(
                code: "validation_error",
                message: "Validation failed",
                details: details
            ))
        case .notFound(let message):
            return ErrorResponse(error: .init(
                code: "not_found",
                message: message,
                details: [:]
            ))
        case .internalError(let message):
            return ErrorResponse(error: .init(
                code: "internal_error",
                message: message,
                details: [:]
            ))
        }
    }
}

/// Cursor data for pagination
struct CursorData: Codable {
    let id: Int
    let createdAt: String
}
