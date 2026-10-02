package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/response"
)

// CSRFCookieName はCSRFトークンを保持するCookie名。
// CookieCsrfTokenRepository.DEFAULT_CSRF_COOKIE_NAME（springboot）に合わせる。
const CSRFCookieName = "XSRF-TOKEN"

// CSRFHeaderName はCSRFトークンを検証するリクエストヘッダー名。
// CookieCsrfTokenRepository.DEFAULT_CSRF_HEADER_NAME（springboot）に合わせる。
const CSRFHeaderName = "X-XSRF-TOKEN"

// CsrfMiddleware はCookie認証（sandbox-spa-react向け）のリクエストに対してCSRFトークンを検証する。
// SecurityConfig（csrf設定）+ CsrfCookieFilter（springboot）に相当。
// Bearer方式（sandbox-app-flutter、およびlogin/web・login/app自体はAuthorizationヘッダーで
// トークンを送るためBearer扱い）と、GET/HEAD/OPTIONSの安全なメソッドはCSRF検証をスキップする。
func CsrfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, cookieErr := c.Cookie(CSRFCookieName)
		if cookieErr != nil || token == "" {
			generated, err := generateCSRFToken()
			if err != nil {
				slog.Error("CSRFトークン生成エラー", "error", err.Error())
				c.JSON(http.StatusInternalServerError, response.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Error:   "INTERNAL_SERVER_ERROR",
					Message: "failed to generate CSRF token",
				})
				c.Abort()
				return
			}
			token = generated
		}

		// CsrfCookieFilter（springboot）と同様、毎リクエストでCookieを再発行する。
		// JS側（document.cookie）で読み取ってヘッダーに載せ返す必要があるためHttpOnly=falseにする。
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(CSRFCookieName, token, 0, "/", "", true, false)

		if !isCSRFProtected(c) {
			c.Next()
			return
		}

		headerToken := c.GetHeader(CSRFHeaderName)
		if headerToken == "" || subtle.ConstantTimeCompare([]byte(headerToken), []byte(token)) != 1 {
			slog.Error("CSRFトークン検証エラー", "path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, response.ErrorResponse{
				Status:  http.StatusForbidden,
				Error:   "FORBIDDEN",
				Message: "CSRF token mismatch",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isCSRFProtected はCSRF検証が必要なリクエストかどうかを判定する。
// Bearer方式のリクエスト（Flutter、およびlogin/web・login/app自体）と、
// 安全なメソッド（GET/HEAD/OPTIONS）は対象外とする。
// JwtAuthFilter.isBearerRequest()（springboot）に相当する判定をAuthorizationヘッダーの有無で行う。
func isCSRFProtected(c *gin.Context) bool {
	if ResolveBearerToken(c) != "" {
		return false
	}
	switch c.Request.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	default:
		return true
	}
}

func generateCSRFToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
