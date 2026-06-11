package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxusecase "sandbox-api-gin/internal/application/usecase/fx"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type BarDataController struct {
	searchUseCase *fxusecase.SearchBarDataUseCase
	statusUseCase *fxusecase.StatusBarDataUseCase
}

func NewBarDataController(
	searchUseCase *fxusecase.SearchBarDataUseCase,
	statusUseCase *fxusecase.StatusBarDataUseCase,
) *BarDataController {
	return &BarDataController{
		searchUseCase: searchUseCase,
		statusUseCase: statusUseCase,
	}
}

// Search POST /v1/fx/bar-data
func (ctrl *BarDataController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.BarDataSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(req.BarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	cmd := fxcommand.SearchBarDataCommand{
		Symbol:      req.Symbol,
		BarType:     barType,
		BarDateFrom: req.BarDateFrom,
		BarDateTo:   req.BarDateTo,
		SortAsc:     req.SortAsc,
		Page:        req.Page,
		Size:        req.Size,
	}

	result, err := ctrl.searchUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.BarDataSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.BarDataList,
	})
}

// Status GET /v1/fx/bar-data/:symbolType/:barType
func (ctrl *BarDataController) Status(c *gin.Context) {
	ctx := c.Request.Context()

	symbolType := c.Param("symbolType")
	if symbolType != "Trade" && symbolType != "Analyze" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "symbolTypeは'Trade'または'Analyze'のみ有効です",
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(c.Param("barType"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	result, err := ctrl.statusUseCase.Execute(ctx, fxcommand.StatusBarDataCommand{
		SymbolType: symbolType,
		BarType:    barType,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
