package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/request"
	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	summertime "sandbox-api-gin/internal/application/usecase/fx/summertime"
)

type SummerTimeController struct {
	searchUseCase *summertime.SearchSummerTimeUseCase
	addUseCase    *summertime.AddSummerTimeUseCase
	getUseCase    *summertime.GetSummerTimeUseCase
	updateUseCase *summertime.UpdateSummerTimeUseCase
}

func NewSummerTimeController(
	searchUseCase *summertime.SearchSummerTimeUseCase,
	addUseCase *summertime.AddSummerTimeUseCase,
	getUseCase *summertime.GetSummerTimeUseCase,
	updateUseCase *summertime.UpdateSummerTimeUseCase,
) *SummerTimeController {
	return &SummerTimeController{
		searchUseCase: searchUseCase,
		addUseCase:    addUseCase,
		getUseCase:    getUseCase,
		updateUseCase: updateUseCase,
	}
}

// Search POST /v1/fx/summer-time/search
func (ctrl *SummerTimeController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ApiSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	result, err := ctrl.searchUseCase.Execute(ctx, req.Page, req.Size)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, fxresponse.SummerTimeSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.SummerTimeList,
	})
}

// Add POST /v1/fx/summer-time
func (ctrl *SummerTimeController) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.SummerTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.addUseCase.Execute(ctx, req.SummerTime); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Get GET /v1/fx/summer-time/:targetYear
func (ctrl *SummerTimeController) Get(c *gin.Context) {
	ctx := c.Request.Context()

	targetYear, err := parseTargetYear(c.Param("targetYear"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	dto, err := ctrl.getUseCase.Get(ctx, targetYear)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Update PUT /v1/fx/summer-time/:targetYear
func (ctrl *SummerTimeController) Update(c *gin.Context) {
	ctx := c.Request.Context()

	baseYear, err := parseTargetYear(c.Param("targetYear"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	var req fxrequest.SummerTimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.updateUseCase.Execute(ctx, baseYear, req.SummerTime); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func parseTargetYear(s string) (int16, error) {
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("targetYear は整数で指定してください: %s", s)
	}
	return int16(v), nil
}
