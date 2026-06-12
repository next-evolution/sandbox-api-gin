package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/request"
	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	country "sandbox-api-gin/internal/application/usecase/fx/country"
)

type CountryController struct {
	searchUseCase *country.SearchCountryUseCase
	addUseCase    *country.AddCountryUseCase
	getUseCase    *country.GetCountryUseCase
	updateUseCase *country.UpdateCountryUseCase
}

func NewCountryController(
	searchUseCase *country.SearchCountryUseCase,
	addUseCase *country.AddCountryUseCase,
	getUseCase *country.GetCountryUseCase,
	updateUseCase *country.UpdateCountryUseCase,
) *CountryController {
	return &CountryController{
		searchUseCase: searchUseCase,
		addUseCase:    addUseCase,
		getUseCase:    getUseCase,
		updateUseCase: updateUseCase,
	}
}

// Search POST /v1/fx/country/search
func (ctrl *CountryController) Search(c *gin.Context) {
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

	c.JSON(http.StatusOK, fxresponse.CountrySearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.CountryList,
	})
}

// Add POST /v1/fx/country
func (ctrl *CountryController) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.CountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.addUseCase.Execute(ctx, req.Country); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Get GET /v1/fx/country/:code
func (ctrl *CountryController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	dto, err := ctrl.getUseCase.Get(ctx, code)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Update PUT /v1/fx/country/:code
func (ctrl *CountryController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	baseCode := c.Param("code")

	var req fxrequest.CountryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.updateUseCase.Execute(ctx, baseCode, req.Country); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
