package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/request"
	"sandbox-api-gin/internal/api/dto/response"
	"sandbox-api-gin/internal/api/middleware"
	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/usecase"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
)

type AuthController struct {
	loginUseCase  *usecase.LoginUseCase
	logoutUseCase *usecase.LogoutUseCase
}

func NewAuthController(loginUseCase *usecase.LoginUseCase, logoutUseCase *usecase.LogoutUseCase) *AuthController {
	return &AuthController{
		loginUseCase:  loginUseCase,
		logoutUseCase: logoutUseCase,
	}
}

// Login POST /v1/auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := &command.LoginCommand{
		AuthUser:     authUser,
		EncodedEmail: req.Email,
	}

	userDto, err := ctrl.loginUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	returnCode := response.ReturnCodeOk
	if userDto == nil {
		returnCode = response.ReturnCodeWarn
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		ApiResponse: response.ApiResponse{ReturnCode: returnCode},
		User:        userDto,
	})
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
	}
}
