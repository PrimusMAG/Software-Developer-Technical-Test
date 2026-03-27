package response

import "github.com/gin-gonic/gin"

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"totalPages,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta Meta) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *gin.Context, status int, message, code string) {
	traceID := c.GetString("request_id")
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"code":    code,
		"traceId": traceID,
	})
}
