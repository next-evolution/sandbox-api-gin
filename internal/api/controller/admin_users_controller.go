package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	adminrequest "sandbox-api-gin/internal/api/dto/request/admin"
	"sandbox-api-gin/internal/api/dto/response"
	admincommand "sandbox-api-gin/internal/application/command/admin"
	userusecase "sandbox-api-gin/internal/application/usecase/user"
	"sandbox-api-gin/internal/domain/apperror"
)

type AdminUsersController struct {
	searchUsersUseCase *userusecase.SearchUsersUseCase
	approveUserUseCase *userusecase.ApproveUserUseCase
	blockUserUseCase   *userusecase.BlockUserUseCase
	grantAdminUseCase  *userusecase.GrantAdminUseCase
}

func NewAdminUsersController(
	searchUsersUseCase *userusecase.SearchUsersUseCase,
	approveUserUseCase *userusecase.ApproveUserUseCase,
	blockUserUseCase *userusecase.BlockUserUseCase,
	grantAdminUseCase *userusecase.GrantAdminUseCase,
) *AdminUsersController {
	return &AdminUsersController{
		searchUsersUseCase: searchUsersUseCase,
		approveUserUseCase: approveUserUseCase,
		blockUserUseCase:   blockUserUseCase,
		grantAdminUseCase:  grantAdminUseCase,
	}
}

// Search POST /v1/admin/users
func (ctrl *AdminUsersController) Search(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	var req adminrequest.UserSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	result, err := ctrl.searchUsersUseCase.Execute(ctx, &admincommand.SearchUsersCommand{
		EmailAddress: req.EmailAddress,
		Approved:     req.Approved,
		Page:         req.Page,
		Size:         req.Size,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserSearchResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage,
		List:        result.List,
	})
}

// Approved PUT /v1/admin/users/approved/:userId
func (ctrl *AdminUsersController) Approved(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	userID := decodeBase64UserID(c.Param("userId"))

	userDto, err := ctrl.approveUserUseCase.Execute(ctx, &admincommand.ApproveUserCommand{
		UserID:    userID,
		UpdatedBy: authUser.Sub,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}

// Block PUT /v1/admin/users/block/:userId
func (ctrl *AdminUsersController) Block(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	userID := decodeBase64UserID(c.Param("userId"))

	var req adminrequest.UserBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	userDto, err := ctrl.blockUserUseCase.Execute(ctx, &admincommand.BlockUserCommand{
		UserID:    userID,
		Blocked:   *req.Blocked,
		UpdatedBy: authUser.Sub,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}

// GrantAdmin PUT /v1/admin/users/admin/:userId
func (ctrl *AdminUsersController) GrantAdmin(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	userID := decodeBase64UserID(c.Param("userId"))

	var req adminrequest.UserAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	userDto, err := ctrl.grantAdminUseCase.Execute(ctx, &admincommand.GrantAdminCommand{
		UserID:    userID,
		Admin:     *req.Admin,
		UpdatedBy: authUser.Sub,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.UserResponse{
		ApiResponse: response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		User:        userDto,
	})
}
