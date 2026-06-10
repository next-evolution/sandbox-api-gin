package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/domain/repository"
	"sandbox-api-gin/internal/security"
)

const AuthUserKey = "authUser"

// JwtMiddleware はJWTを検証してAuthUserをコンテキストにセットする。
// JwtAuthFilterに相当する処理。
func JwtMiddleware(jwtProvider *security.JwtProvider, sessionRepo repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := resolveToken(c)
		if token == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		jwtAuthUser, err := jwtProvider.Parse(token)
		if err != nil {
			slog.Error("JWT検証エラー", "path", c.Request.URL.Path, "error", err.Error())
			c.Next()
			return
		}

		// Redisからadminフラグ付きAuthUserを取得（ログイン済みの場合）
		sessionUser, err := sessionRepo.FindBySub(ctx, jwtAuthUser.Sub)
		if err != nil {
			slog.Error("Redisセッション取得エラー", "error", err.Error())
		}

		authUser := jwtAuthUser
		if sessionUser != nil {
			authUser = sessionUser
		}

		c.Set(AuthUserKey, authUser)

		// セッションTTL更新
		if err := sessionRepo.Update(ctx, authUser); err != nil {
			slog.Error("セッションTTL更新エラー", "error", err.Error())
		}

		c.Next()
	}
}

// resolveToken はAuthorizationヘッダーからBearerトークンを取り出す
func resolveToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}
