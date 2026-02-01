package com.lebonpoint.shared

/**
 * Base domain exception
 */
sealed class DomainException(
    message: String,
    cause: Throwable? = null
) : Exception(message, cause) {
    /**
     * Error code for API responses
     */
    abstract val code: String

    /**
     * Additional error details
     */
    open val details: Map<String, Any?> = emptyMap()
}

/**
 * Validation exception
 * Thrown when input validation fails
 */
class ValidationException(
    override val details: Map<String, Any?>
) : DomainException("Validation failed") {
    override val code: String = "validation_error"
}

/**
 * Not found exception
 * Thrown when a requested resource is not found
 */
class NotFoundException(
    resource: String = "Resource",
    id: Any? = null
) : DomainException(
    message = if (id != null) "$resource with id '$id' not found" else "$resource not found"
) {
    override val code: String = "not_found"
    override val details: Map<String, Any?> = mapOf(
        "resource" to resource,
        "id" to id
    )
}

/**
 * Internal error exception
 * Thrown when an unexpected internal error occurs
 */
class InternalException(
    message: String,
    cause: Throwable? = null
) : DomainException(message, cause) {
    override val code: String = "internal_error"
}

/**
 * Create a ValidationException with a single field error
 */
fun validationError(field: String, message: String): ValidationException {
    return ValidationException(mapOf(field to message))
}

/**
 * Create a ValidationException with multiple field errors
 */
fun validationErrors(errors: Map<String, String>): ValidationException {
    return ValidationException(errors)
}
