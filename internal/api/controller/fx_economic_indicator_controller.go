package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	economicindicator "sandbox-api-gin/internal/application/usecase/fx/economicindicator"
)

type EconomicIndicatorController struct {
	searchUseCase *economicindicator.SearchEconomicIndicatorUseCase
	getUseCase    *economicindicator.GetEconomicIndicatorUseCase
	addUseCase    *economicindicator.AddEconomicIndicatorUseCase
	updateUseCase *economicindicator.UpdateEconomicIndicatorUseCase
}

func NewEconomicIndicatorController(
	searchUseCase *economicindicator.SearchEconomicIndicatorUseCase,
	getUseCase *economicindicator.GetEconomicIndicatorUseCase,
	addUseCase *economicindicator.AddEconomicIndicatorUseCase,
	updateUseCase *economicindicator.UpdateEconomicIndicatorUseCase,
) *EconomicIndicatorController {
	return &EconomicIndicatorController{
		searchUseCase: searchUseCase,
		getUseCase:    getUseCase,
		addUseCase:    addUseCase,
		updateUseCase: updateUseCase,
	}
}

// Search POST /v1/fx/economic-indicator/search
func (ctrl *EconomicIndicatorController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.EconomicIndicatorSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	result, err := ctrl.searchUseCase.Execute(ctx, req.Page, req.Size, req.CountryCode, req.Importance, req.Name)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, fxresponse.EconomicIndicatorSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.List,
	})
}

// Get GET /v1/fx/economic-indicator/:countryCode/:id
func (ctrl *EconomicIndicatorController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	countryCode := c.Param("countryCode")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "invalid id",
		})
		return
	}

	dto, err := ctrl.getUseCase.Execute(ctx, countryCode, id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Add POST /v1/fx/economic-indicator
func (ctrl *EconomicIndicatorController) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.EconomicIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.addUseCase.Execute(ctx, req.Indicator); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Update PUT /v1/fx/economic-indicator/:countryCode/:id
func (ctrl *EconomicIndicatorController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	countryCode := c.Param("countryCode")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "invalid id",
		})
		return
	}

	var req fxrequest.EconomicIndicatorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.updateUseCase.Execute(ctx, countryCode, id, req.Indicator); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
