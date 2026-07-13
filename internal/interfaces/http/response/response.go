package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func SuccessWithMeta(c *gin.Context, status int, data any, meta any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
		"meta":    meta,
	})
}

func Message(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": true,
		"message": message,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   message,
	})
}
