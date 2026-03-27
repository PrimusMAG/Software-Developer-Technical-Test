package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type logLine struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	TraceID   string `json:"traceId"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latencyMs"`
	ClientIP  string `json:"clientIp"`
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		line := logLine{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Level:     "info",
			TraceID:   c.GetString("request_id"),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			LatencyMs: time.Since(start).Milliseconds(),
			ClientIP:  c.ClientIP(),
		}
		raw, _ := json.Marshal(line)
		log.Println(string(raw))
	}
}
