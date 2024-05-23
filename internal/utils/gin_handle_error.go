package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinHandleError handles errors by logging them and sending a JSON response.
func GinHandleError(c *gin.Context, log logrus.FieldLogger, status int, err error, msg string) {
	// Log the error with different levels based on the status code
	switch {
	case status >= http.StatusInternalServerError:
		log.WithError(err).Error(msg) // Internal server errors are logged as errors
	case status >= http.StatusBadRequest:
		log.WithError(err).Warn(msg) // Client errors are logged as warnings
	default:
		log.WithError(err).Info(msg) // Other cases (e.g., redirects) are logged as info
	}

	c.JSON(status, gin.H{"error": msg})
}
