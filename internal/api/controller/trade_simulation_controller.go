package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxusecase "sandbox-api-gin/internal/application/usecase/fx"
	"sandbox-api-gin/internal/domain/model/fx"
)

type TradeSimulationController struct {
	useCase *fxusecase.TradeSimulationUseCase
}

func NewTradeSimulationController(useCase *fxusecase.TradeSimulationUseCase) *TradeSimulationController {
	return &TradeSimulationController{useCase: useCase}
}

// Simulation POST /v1/fx/trade/simulation
func (ctrl *TradeSimulationController) Simulation(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.TradeSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  http.StatusBadRequest,
			Error:   "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	positionList := make([]*fx.TradePosition, 0, len(req.PositionList))
	for i := range req.PositionList {
		positionList = append(positionList, req.PositionList[i].ToDomain())
	}

	cmd := &fxcommand.TradeSimulationCommand{
		RiskAmount:    req.RiskAmount,
		FirstLotRatio: req.FirstLotRatio,
		Entry:         req.Entry.ToDomain(),
		PositionList:  positionList,
	}

	result, err := ctrl.useCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.TradeSimulationResponse{
		ApiResponse:  response.ApiResponse{ReturnCode: response.ReturnCodeOk},
		Entry:        result.Entry,
		PositionList: result.PositionList,
	})
}
