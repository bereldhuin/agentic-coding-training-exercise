package middleware

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/hymaia/lebonpoint-app/server/go/internal/shared/errors"
)

// ErrorHandler handles errors and returns consistent error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) == 0 {
			return
		}

		// Get the last error
		err := c.Errors.Last().Err

		// Convert to AppError
		appErr := apperrors.GetAppError(err)

		// Build response
		response := appErr.ToErrorResponse()

		// Set status code
		c.JSON(appErr.StatusCode, response)
	}
}

// RecoveryMiddleware handles panics and converts them to errors
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.Error(apperrors.NewInternalError(err))
		} else {
			c.Error(apperrors.NewInternalError("An unexpected error occurred"))
		}
		c.Abort()
	})
}
