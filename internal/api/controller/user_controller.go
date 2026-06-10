package controller

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/request"
	"sandbox-api-gin/internal/api/dto/response"
	"sandbox-api-gin/internal/application/command"
	"sandbox-api-gin/internal/application/usecase"
	"sandbox-api-gin/internal/domain/apperror"
)

type UserController struct {
	getProfileUseCase   *usecase.GetProfileUseCase
	registerUserUseCase *usecase.RegisterUserUseCase
	updateUserUseCase   *usecase.UpdateUserUseCase
}

func NewUserController(
	getProfileUseCase *usecase.GetProfileUseCase,
	registerUserUseCase *usecase.RegisterUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
) *UserController {
	return &UserController{
		getProfileUseCase:   getProfileUseCase,
		registerUserUseCase: registerUserUseCase,
		updateUserUseCase:   updateUserUseCase,
	}
}

// Profile GET /v1/user
func (ctrl *UserController) Profile(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	cmd := &command.GetProfileCommand{UserID: authUser.Sub}
	userDto, err := ctrl.getProfileUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	if userDto == nil {
		c.JSON(http.StatusOK, response.ApiResponse{
			ReturnCode: response.ReturnCodeWarn,
			Message:    "利用承認待ちです",
		})
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}

// Registration POST /v1/user
func (ctrl *UserController) Registration(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	var req request.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := &command.RegisterUserCommand{
		UserID:   authUser.Sub,
		Email:    authUser.Email,
		NickName: req.NickName,
	}
	userDto, err := ctrl.registerUserUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}

// Update PUT /v1/user/:userId
func (ctrl *UserController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	userIDBase64 := c.Param("userId")
	userID := decodeBase64UserID(userIDBase64)
	if authUser.Sub != userID {
		handleError(c, apperror.NewForbiddenError("他のユーザの情報は更新できません"))
		return
	}

	var req request.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	cmd := &command.UpdateUserCommand{
		UserID:    authUser.Sub,
		NickName:  req.NickName,
		UpdatedBy: authUser.Sub,
	}
	userDto, err := ctrl.updateUserUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}

// decodeBase64UserID はパスパラメータの Base64 エンコード済み userId をデコードする。
// JavaのUserId.decodeUserIdValue()に相当。
func decodeBase64UserID(encoded string) string {
	if b, err := base64.StdEncoding.DecodeString(encoded); err == nil {
		return string(b)
	}
	if b, err := base64.RawStdEncoding.DecodeString(encoded); err == nil {
		return string(b)
	}
	return ""
}
