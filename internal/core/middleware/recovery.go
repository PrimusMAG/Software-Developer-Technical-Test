package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"football-api/internal/shared/response"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("panic recovered trace_id=%s err=%v", c.GetString("request_id"), recovered)
		response.Error(c, http.StatusInternalServerError, "internal server error", "INTERNAL_ERROR")
		c.Abort()
	})
}
