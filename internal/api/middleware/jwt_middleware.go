package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/domain/model"
	"sandbox-api-gin/internal/domain/repository"
	"sandbox-api-gin/internal/security"
)

const AuthUserKey = "authUser"

// JWTCookieName はReact（Web）向けにJWTを格納するHttpOnly Cookie名。
// JwtCookieProvider.COOKIE_NAME（springboot）に合わせる。
const JWTCookieName = "sandbox_jwt"

// JwtMiddleware はJWTを検証してAuthUserをコンテキストにセットする。
// JwtAuthFilterに相当する処理。
func JwtMiddleware(jwtProvider *security.JwtProvider, sessionRepo repository.SessionRepository, userRepo repository.UserRepository) gin.HandlerFunc {
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

		var authUser *model.AuthUser

		// Redisからadmin/approvedフラグ付きAuthUserを取得（ログイン済みの場合）
		sessionUser, err := sessionRepo.FindBySub(ctx, jwtAuthUser.Sub)
		if err != nil {
			slog.Error("Redisセッション取得エラー", "error", err.Error())
		}

		if sessionUser != nil {
			authUser = sessionUser
			if err := sessionRepo.Update(ctx, authUser); err != nil {
				slog.Error("セッションTTL更新エラー", "error", err.Error())
			}
		} else {
			// セッションなし → DBから admin/approved を復元して Redis に保存（silent login）
			user, err := userRepo.FindByUserID(ctx, jwtAuthUser.Sub)
			if err != nil {
				slog.Error("silent loginエラー", "path", c.Request.URL.Path, "error", err.Error())
			} else if user != nil {
				authUser = &model.AuthUser{
					Sub:           jwtAuthUser.Sub,
					Email:         jwtAuthUser.Email,
					EmailVerified: jwtAuthUser.EmailVerified,
					Admin:         user.Admin,
					Approved:      user.Approved,
				}
				if err := sessionRepo.Save(ctx, authUser); err != nil {
					slog.Error("セッション保存エラー", "error", err.Error())
				}
			} else {
				// sandbox_user未登録（初回ログイン）でもJWT由来のauthUserにフォールバックする
				authUser = jwtAuthUser
			}
		}

		if authUser == nil {
			c.Next()
			return
		}

		c.Set(AuthUserKey, authUser)
		c.Next()
	}
}

// resolveToken はAuthorizationヘッダー優先、無ければCookieからトークンを取り出す。
// JwtAuthFilter.resolveToken()（springboot）に相当。
func resolveToken(c *gin.Context) string {
	if token := ResolveBearerToken(c); token != "" {
		return token
	}
	if cookie, err := c.Cookie(JWTCookieName); err == nil {
		return cookie
	}
	return ""
}

// ResolveBearerToken はAuthorizationヘッダーからBearerトークンを取り出す。
// BearerTokenResolver.resolve()（springboot）に相当。ログイン時のCookie発行に使う。
func ResolveBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}
