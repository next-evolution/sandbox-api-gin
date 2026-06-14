package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/response"
	"sandbox-api-gin/internal/domain/apperror"
	fxusecase "sandbox-api-gin/internal/application/usecase/fx"
)

type AdminMasterRefreshController struct {
	masterStatusUseCase  *fxusecase.MasterStatusUseCase
	masterRefreshUseCase *fxusecase.MasterRefreshUseCase
}

func NewAdminMasterRefreshController(
	masterStatusUseCase *fxusecase.MasterStatusUseCase,
	masterRefreshUseCase *fxusecase.MasterRefreshUseCase,
) *AdminMasterRefreshController {
	return &AdminMasterRefreshController{
		masterStatusUseCase:  masterStatusUseCase,
		masterRefreshUseCase: masterRefreshUseCase,
	}
}

// Status GET /v1/admin/master-refresh
func (ctrl *AdminMasterRefreshController) Status(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	msg, err := ctrl.masterStatusUseCase.Execute(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiResponse{
		ReturnCode: response.ReturnCodeOk,
		Message:    msg,
	})
}

// Refresh PUT /v1/admin/master-refresh
func (ctrl *AdminMasterRefreshController) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	if !authUser.IsAdmin() {
		handleError(c, apperror.NewForbiddenError("管理者用APIです"))
		return
	}

	msg, err := ctrl.masterRefreshUseCase.Execute(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ApiResponse{
		ReturnCode: response.ReturnCodeOk,
		Message:    msg,
	})
}
