package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	economicindicatordata "sandbox-api-gin/internal/application/usecase/fx/economicindicatordata"
)

type EconomicIndicatorDataController struct {
	searchUseCase *economicindicatordata.SearchEconomicIndicatorDataUseCase
	getUseCase    *economicindicatordata.GetEconomicIndicatorDataUseCase
	addUseCase    *economicindicatordata.AddEconomicIndicatorDataUseCase
	updateUseCase *economicindicatordata.UpdateEconomicIndicatorDataUseCase
	importUseCase *economicindicatordata.ImportEconomicIndicatorDataUseCase
}

func NewEconomicIndicatorDataController(
	searchUseCase *economicindicatordata.SearchEconomicIndicatorDataUseCase,
	getUseCase *economicindicatordata.GetEconomicIndicatorDataUseCase,
	addUseCase *economicindicatordata.AddEconomicIndicatorDataUseCase,
	updateUseCase *economicindicatordata.UpdateEconomicIndicatorDataUseCase,
	importUseCase *economicindicatordata.ImportEconomicIndicatorDataUseCase,
) *EconomicIndicatorDataController {
	return &EconomicIndicatorDataController{
		searchUseCase: searchUseCase,
		getUseCase:    getUseCase,
		addUseCase:    addUseCase,
		updateUseCase: updateUseCase,
		importUseCase: importUseCase,
	}
}

// Search POST /v1/fx/economic-indicator-data/search
func (ctrl *EconomicIndicatorDataController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.EconomicIndicatorDataSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	result, err := ctrl.searchUseCase.Execute(ctx, req.ID, req.Importance, req.CountryCode, req.PublicationBaseDate, req.Page, req.Size, req.SortAsc)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, fxresponse.EconomicIndicatorDataSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.List,
	})
}

// Get GET /v1/fx/economic-indicator-data/:economicIndicatorId/:publication
func (ctrl *EconomicIndicatorDataController) Get(c *gin.Context) {
	ctx := c.Request.Context()

	id, publication, ok := parseIDAndPublication(c)
	if !ok {
		return
	}

	dto, err := ctrl.getUseCase.Execute(ctx, id, publication)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Add POST /v1/fx/economic-indicator-data
func (ctrl *EconomicIndicatorDataController) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.EconomicIndicatorDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.addUseCase.Execute(ctx, req.Data); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Update PUT /v1/fx/economic-indicator-data/:economicIndicatorId/:publication
func (ctrl *EconomicIndicatorDataController) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, publication, ok := parseIDAndPublication(c)
	if !ok {
		return
	}

	var req fxrequest.EconomicIndicatorDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.updateUseCase.Execute(ctx, id, publication, req.Data); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// ImportText POST /v1/fx/economic-indicator-data/import-text
func (ctrl *EconomicIndicatorDataController) ImportText(c *gin.Context) {
	ctx := c.Request.Context()
	authUser := getAuthUser(c)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	fileHeaders := form.File["uploadFileList"]
	files := make([]economicindicatordata.FileEntry, 0, len(fileHeaders))
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Status:  http.StatusBadRequest,
				Error:   "BAD_REQUEST",
				Message: "ファイル読み込みに失敗しました: " + fh.Filename,
			})
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				_ = err
			}
		}()
		files = append(files, economicindicatordata.FileEntry{
			FileName: fh.Filename,
			Reader:   f,
			FileSize: fh.Size,
		})
	}

	results, err := ctrl.importUseCase.Execute(ctx, files, authUser.Sub)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

func parseIDAndPublication(c *gin.Context) (int64, time.Time, bool) {
	idStr := c.Param("economicIndicatorId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "invalid economicIndicatorId",
		})
		return 0, time.Time{}, false
	}

	pubStr := c.Param("publication")
	publication, err := time.Parse("2006-01-02 15:04:05", pubStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "invalid publication: expected yyyy-MM-dd HH:mm:ss",
		})
		return 0, time.Time{}, false
	}

	return id, publication, true
}
