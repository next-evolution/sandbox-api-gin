package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	symbol "sandbox-api-gin/internal/application/usecase/fx/symbol"
)

type SymbolController struct {
	searchUseCase *symbol.SearchSymbolUseCase
	addUseCase    *symbol.AddSymbolUseCase
	getUseCase    *symbol.GetSymbolUseCase
	updateUseCase *symbol.UpdateSymbolUseCase
}

func NewSymbolController(
	searchUseCase *symbol.SearchSymbolUseCase,
	addUseCase *symbol.AddSymbolUseCase,
	getUseCase *symbol.GetSymbolUseCase,
	updateUseCase *symbol.UpdateSymbolUseCase,
) *SymbolController {
	return &SymbolController{
		searchUseCase: searchUseCase,
		addUseCase:    addUseCase,
		getUseCase:    getUseCase,
		updateUseCase: updateUseCase,
	}
}

// CurrencyPairList GET /v1/fx/symbol/currency-pair-list
func (ctrl *SymbolController) CurrencyPairList(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := ctrl.searchUseCase.Execute(ctx, "Trade", 1, 500)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result.SymbolList)
}

// CurrencyIndexList GET /v1/fx/symbol/currency-index-list
func (ctrl *SymbolController) CurrencyIndexList(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := ctrl.searchUseCase.Execute(ctx, "Analyze", 1, 500)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result.SymbolList)
}

// Search POST /v1/fx/symbol/search
func (ctrl *SymbolController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.SymbolSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if req.SymbolType != "Trade" && req.SymbolType != "Analyze" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: "symbolTypeは'Trade'または'Analyze'のみ有効です",
		})
		return
	}

	result, err := ctrl.searchUseCase.Execute(ctx, req.SymbolType, req.Page, req.Size)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, fxresponse.SymbolSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        result.SymbolList,
	})
}

// Add POST /v1/fx/symbol
func (ctrl *SymbolController) Add(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.SymbolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.addUseCase.Execute(ctx, req.Symbol); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Get GET /v1/fx/symbol/:symbol
func (ctrl *SymbolController) Get(c *gin.Context) {
	ctx := c.Request.Context()
	symbol := c.Param("symbol")

	dto, err := ctrl.getUseCase.Get(ctx, symbol)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Update PUT /v1/fx/symbol/:symbol
func (ctrl *SymbolController) Update(c *gin.Context) {
	ctx := c.Request.Context()
	baseSymbol := c.Param("symbol")

	var req fxrequest.SymbolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	if err := ctrl.updateUseCase.Execute(ctx, baseSymbol, req.Symbol); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
