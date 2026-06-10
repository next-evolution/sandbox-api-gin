package middleware

import (
	"log/slog"
	"net/http"
	"sandbox-api-gin/internal/api/dto/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware はAuthUserがコンテキストに存在しない場合401を返す。
// AuthInterceptorに相当する処理。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get(AuthUserKey); !exists {
			slog.Error("Unauthorized", "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Error:   "UNAUTHORIZED",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
