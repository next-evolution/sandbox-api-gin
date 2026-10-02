package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/request"
	"sandbox-api-gin/internal/api/dto/response"
	"sandbox-api-gin/internal/api/middleware"
	"sandbox-api-gin/internal/application/command"
	userusecase "sandbox-api-gin/internal/application/usecase/user"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
)

type AuthController struct {
	loginUseCase  *userusecase.LoginUseCase
	logoutUseCase *userusecase.LogoutUseCase
	sessionTTL    int
}

func NewAuthController(loginUseCase *userusecase.LoginUseCase, logoutUseCase *userusecase.LogoutUseCase, sessionTTL int) *AuthController {
	return &AuthController{
		loginUseCase:  loginUseCase,
		logoutUseCase: logoutUseCase,
		sessionTTL:    sessionTTL,
	}
}

// LoginWeb POST /v1/auth/login/web
// sandbox-spa-react（Web）向け。ログイン成立前はCookieが無いためAuthorizationヘッダーで受け取り、
// 成功時はJWTをHttpOnly CookieとしてSet-Cookieする（レスポンスボディにトークンは含めない）。
func (ctrl *AuthController) LoginWeb(c *gin.Context) {
	body, ok := ctrl.login(c)
	if !ok {
		return
	}

	if token := middleware.ResolveBearerToken(c); token != "" {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(middleware.JWTCookieName, token, ctrl.sessionTTL, "/", "", true, true)
	}

	c.JSON(http.StatusOK, body)
}

// LoginApp POST /v1/auth/login/app
// sandbox-app-flutter向け。従来どおりBearer方式のみ・Cookie発行は行わない。
func (ctrl *AuthController) LoginApp(c *gin.Context) {
	body, ok := ctrl.login(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, body)
}

func (ctrl *AuthController) login(c *gin.Context) (response.LoginResponse, bool) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return response.LoginResponse{}, false
	}

	cmd := &command.LoginCommand{
		AuthUser:     authUser,
		EncodedEmail: req.Email,
	}

	userDto, err := ctrl.loginUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return response.LoginResponse{}, false
	}

	returnCode := response.ReturnCodeOk
	if userDto == nil {
		returnCode = response.ReturnCodeWarn
	}

	return response.LoginResponse{
		ApiResponse: response.ApiResponse{ReturnCode: returnCode},
		User:        userDto,
	}, true
}

// Logout POST /v1/auth/logout-api
func (ctrl *AuthController) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	var req request.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := &command.LogoutCommand{
		AuthUser:      authUser,
		EncodedUserID: req.UserID,
	}
	ctrl.logoutUseCase.Execute(ctx, cmd)

	// sandbox-spa-react（Web）向けCookieを失効させる。sandbox-app-flutter（App）はCookie未使用のため無害
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(middleware.JWTCookieName, "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, response.ApiResponse{ReturnCode: response.ReturnCodeOk})
}

func getAuthUser(c *gin.Context) *model.AuthUser {
	val, exists := c.Get(middleware.AuthUserKey)
	if !exists {
		return nil
	}
	authUser, _ := val.(*model.AuthUser)
	return authUser
}

func handleError(c *gin.Context, err error) {
	switch {
	case apperror.IsAuthenticationError(err):
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Error:   "UNAUTHORIZED",
			Message: err.Error(),
		})
	case apperror.IsForbiddenError(err):
		c.JSON(http.StatusForbidden, response.ErrorResponse{
			Status:  http.StatusForbidden,
			Error:   "FORBIDDEN",
			Message: err.Error(),
		})
	case apperror.IsNotFoundError(err):
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Status:  http.StatusNotFound,
			Error:   "NOT_FOUND",
			Message: err.Error(),
		})
	case apperror.IsDuplicateError(err), apperror.IsInsertError(err), apperror.IsUpdateError(err):
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
	default:
		slog.Error("internal server error", "error", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
}
