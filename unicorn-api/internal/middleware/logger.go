package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger returns a gin.HandlerFunc for logging HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			// "client_ip":    param.ClientIP,
			"timestamp": param.TimeStamp.Format(time.RFC3339),
			"method":    param.Method,
			"path":      param.Path,
			// "protocol":     param.Request.Proto,
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			// "user_agent":   param.Request.UserAgent(),
			"error": param.ErrorMessage,
		}).Info("HTTP Request")
		return ""
	})
}
