import Vapor

/// Custom error middleware to catch and format errors consistently
struct ErrorMiddleware: Middleware {
    func respond(to request: Request, chainingTo next: Responder) -> EventLoopFuture<Response> {
        return next.respond(to: request).flatMapErrorThrowing { error in
            return self.handleError(error, on: request)
        }
    }

    private func handleError(_ error: Error, on request: Request) -> Response {
        let response: ErrorResponse
        let status: HTTPResponseStatus

        switch error {
        case let domainError as DomainError:
            response = domainError.toErrorResponse()
            status = domainError.status

        case let abortError as AbortError:
            response = ErrorResponse(error: .init(
                code: errorCodeFromStatus(abortError.status),
                message: abortError.reason,
                details: [:]
            ))
            status = abortError.status

        default:
            response = ErrorResponse(error: .init(
                code: "internal_error",
                message: "An unexpected error occurred",
                details: [:]
            ))
            status = .internalServerError
        }

        // Log error
        request.logger.error("\(status): \(response.error.message)")

        // Return error response
        return Response(
            status: status,
            headers: ["Content-Type": "application/json"],
            body: .init(string: String(data: try! JSONEncoder().encode(response), encoding: .utf8)!)
        )
    }

    private func errorCodeFromStatus(_ status: HTTPResponseStatus) -> String {
        switch status.code {
        case 400: return "validation_error"
        case 404: return "not_found"
        default: return "internal_error"
        }
    }
}
