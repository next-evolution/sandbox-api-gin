package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	fxcommand "sandbox-api-gin/internal/application/command/fx"
	"sandbox-api-gin/internal/application/usecase/fx/bardata"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
)

type BarDataController struct {
	searchUseCase    *bardata.SearchBarDataUseCase
	statusUseCase    *bardata.StatusBarDataUseCase
	importCsvUseCase *bardata.ImportCsvBarDataUseCase
}

func NewBarDataController(
	searchUseCase *bardata.SearchBarDataUseCase,
	statusUseCase *bardata.StatusBarDataUseCase,
	importCsvUseCase *bardata.ImportCsvBarDataUseCase,
) *BarDataController {
	return &BarDataController{
		searchUseCase:    searchUseCase,
		statusUseCase:    statusUseCase,
		importCsvUseCase: importCsvUseCase,
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

	c.JSON(http.StatusOK, fxresponse.BarDataSearchResponse{
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

// ImportCsv POST /v1/fx/bar-data/import-csv/:symbol/:barType/:skipLatest
func (ctrl *BarDataController) ImportCsv(c *gin.Context) {
	ctx := c.Request.Context()

	symbol := c.Param("symbol")
	barTypeStr := c.Param("barType")
	skipLatestStr := c.Param("skipLatest")

	barType, err := fxmodel.BarTypeOf(barTypeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	file, err := c.FormFile("uploadFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "uploadFileが見つかりません",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "ファイルオープンエラー",
		})
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			slog.Error("アップロードファイルクローズエラー", "error", err)
		}
	}()

	authUser := getAuthUser(c)

	result, err := ctrl.importCsvUseCase.Execute(ctx, fxcommand.ImportCsvBarDataCommand{
		Symbol:           symbol,
		BarType:          barType,
		SkipLatest:       skipLatestStr == "true",
		FileReader:       src,
		OriginalFileName: file.Filename,
		FileSize:         file.Size,
		UserSub:          authUser.Sub,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
