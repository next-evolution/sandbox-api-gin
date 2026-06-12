package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sandbox-api-gin/internal/api/dto/response"
	fxusecase "sandbox-api-gin/internal/application/usecase/fx"
)

// 有効なsymbolType値
var validSymbolTypes = map[string]struct{}{
	"Trade":   {},
	"Analyze": {},
}

type MasterListController struct {
	useCase *fxusecase.GetMasterUseCase
}

func NewMasterListController(useCase *fxusecase.GetMasterUseCase) *MasterListController {
	return &MasterListController{useCase: useCase}
}

// Symbol GET /v1/fx/master-list/symbol/:symbolType
func (ctrl *MasterListController) Symbol(c *gin.Context) {
	ctx := c.Request.Context()
	symbolType := c.Param("symbolType")

	if _, ok := validSymbolTypes[symbolType]; !ok {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "symbolTypeは'Trade'または'Analyze'のみ有効です",
		})
		return
	}

	list, err := ctrl.useCase.Symbol(ctx, symbolType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Country GET /v1/fx/master-list/country
func (ctrl *MasterListController) Country(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := ctrl.useCase.Country(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CurrencyPair GET /v1/fx/master-list/currency-pair
func (ctrl *MasterListController) CurrencyPair(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := ctrl.useCase.CurrencyPair(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CurrencyIndex GET /v1/fx/master-list/currency-index
func (ctrl *MasterListController) CurrencyIndex(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := ctrl.useCase.CurrencyIndex(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}

// EconomicIndicator GET /v1/fx/master-list/economic-indicator/:countryCode
func (ctrl *MasterListController) EconomicIndicator(c *gin.Context) {
	ctx := c.Request.Context()
	countryCode := c.Param("countryCode")
	list, err := ctrl.useCase.EconomicIndicator(ctx, countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Error:   "INTERNAL_SERVER_ERROR",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, list)
}
