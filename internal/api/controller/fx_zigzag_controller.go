package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	fxrequest "sandbox-api-gin/internal/api/dto/request/fx"
	"sandbox-api-gin/internal/api/dto/response"
	fxresponse "sandbox-api-gin/internal/api/dto/response/fx"
	fxcommand "sandbox-api-gin/internal/application/command/fx"
	fxdto "sandbox-api-gin/internal/application/dto/fx"
	zigzagusecase "sandbox-api-gin/internal/application/usecase/fx/zigzag"
	fxmodel "sandbox-api-gin/internal/domain/model/fx"
	"sandbox-api-gin/internal/domain/model/fx/zigzag"
)

type ZigZagController struct {
	searchUseCase   *zigzagusecase.SearchZigZagUseCase
	statusUseCase   *zigzagusecase.GetZigZagStatusUseCase
	generateUseCase *zigzagusecase.GenerateZigZagUseCase
	barDataUseCase  *zigzagusecase.GetZigZagBarDataUseCase
}

func NewZigZagController(
	searchUseCase *zigzagusecase.SearchZigZagUseCase,
	statusUseCase *zigzagusecase.GetZigZagStatusUseCase,
	generateUseCase *zigzagusecase.GenerateZigZagUseCase,
	barDataUseCase *zigzagusecase.GetZigZagBarDataUseCase,
) *ZigZagController {
	return &ZigZagController{
		searchUseCase:   searchUseCase,
		statusUseCase:   statusUseCase,
		generateUseCase: generateUseCase,
		barDataUseCase:  barDataUseCase,
	}
}

// Search POST /v1/fx/zigzag
func (ctrl *ZigZagController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.ZigZagSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(req.BarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	cmd := fxcommand.ZigZagSearchCommand{
		BarType:              barType,
		Symbol:               req.Symbol,
		Depth:                int(req.Depth),
		BarDateTimeMin:       req.BarDateTimeMin,
		BarDateTimeMax:       req.BarDateTimeMax,
		Wave:                 req.Wave,
		PreviousWave:         req.PreviousWave,
		NextWave:             req.NextWave,
		Next2Wave:            req.Next2Wave,
		Direction4h200:       req.Direction4h200,
		Direction4h75:        req.Direction4h75,
		Direction4h20:        req.Direction4h20,
		Direction1h200:       req.Direction1h200,
		Direction15m200:      req.Direction15m200,
		Wave4h:               req.Wave4h,
		DirectionTarget4h200: req.DirectionTarget4h200,
		Page:                 req.Page,
		Size:                 req.Size,
	}

	result, err := ctrl.searchUseCase.Execute(ctx, cmd)
	if err != nil {
		handleError(c, err)
		return
	}

	list := make([]fxresponse.ZigZagResult, len(result.List))
	for i, item := range result.List {
		list[i] = toZigZagResult(item)
	}

	c.JSON(http.StatusOK, fxresponse.ZigZagSearchResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  result.TotalCount,
		SearchCount: result.TotalCount,
		TotalPage:   result.TotalPage(),
		List:        list,
	})
}

// Status POST /v1/fx/zigzag/status
func (ctrl *ZigZagController) Status(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.ZigZagStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	if req.SymbolType != "Trade" && req.SymbolType != "Analyze" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST",
			Message: "symbolTypeは'Trade'または'Analyze'のみ有効です",
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(req.BarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	statusList, err := ctrl.statusUseCase.Execute(ctx, fxcommand.ZigZagStatusCommand{
		SymbolType: req.SymbolType,
		BarType:    barType,
		Depth:      int(req.Depth),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	items := make([]fxresponse.ZigZagStatusItem, len(statusList))
	for i, s := range statusList {
		items[i] = toStatusItem(s)
	}

	c.JSON(http.StatusOK, fxresponse.ZigZagStatusResponse{
		ReturnCode:  response.ReturnCodeOk,
		TotalCount:  len(items),
		SearchCount: len(items),
		TotalPage:   1,
		List:        items,
	})
}

// Generate POST /v1/fx/zigzag/generate
func (ctrl *ZigZagController) Generate(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.ZigZagGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(req.BarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	result, err := ctrl.generateUseCase.Execute(ctx, fxcommand.ZigZagGenerateCommand{
		Symbol:      req.Symbol,
		BarType:     barType,
		Depth:       int(req.Depth),
		BarDateTime: req.BarDateTime,
		LoadSize:    req.LoadSize,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	rc := response.ReturnCodeOk
	var msg string
	if result.Warn {
		rc = response.ReturnCodeWarn
		msg = result.Status.Message
	}

	c.JSON(http.StatusOK, fxresponse.ZigZagGenerateResponse{
		ReturnCode: rc,
		Message:    msg,
		Status:     toStatusItem(result.Status),
	})
}

// BarData POST /v1/fx/zigzag/bar-data
func (ctrl *ZigZagController) BarData(c *gin.Context) {
	ctx := c.Request.Context()

	var req fxrequest.ZigZagBarDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST", Message: err.Error(),
		})
		return
	}

	barType, err := fxmodel.BarTypeOf(req.BarType)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: http.StatusBadRequest, Error: "BAD_REQUEST",
			Message: "barTypeは'15M','1H','4H','1D'のいずれかを指定してください",
		})
		return
	}

	result, err := ctrl.barDataUseCase.Execute(ctx, fxcommand.ZigZagBarDataCommand{
		BarType:   barType,
		Symbol:    req.Symbol,
		Depth:     req.Depth,
		WaveStart: req.WaveStart,
		Wave:      req.Wave,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	barDataList := make([]fxresponse.ZigZagBarData, len(result.List))
	for i, row := range result.List {
		barDataList[i] = fxresponse.ZigZagBarData{
			BarDateTime: row.BarDateTime,
			OpenPrice:   row.OpenPrice,
			HighPrice:   row.HighPrice,
			LowPrice:    row.LowPrice,
			ClosePrice:  row.ClosePrice,
			Sma200:      row.Sma200,
			Sma75:       row.Sma75,
			Sma20:       row.Sma20,
		}
	}

	c.JSON(http.StatusOK, fxresponse.ZigZagBarDataResponse{
		ReturnCode:        response.ReturnCodeOk,
		BarType:           string(result.BarType),
		Symbol:            result.Symbol,
		Depth:             result.Depth,
		Wave:              result.Wave,
		ZigZagBarDataList: barDataList,
	})
}

func toZigZagResult(item *fxdto.ZigZagSearchItem) fxresponse.ZigZagResult {
	fractalWaves := make([]fxresponse.ZigZagFractalWave, len(item.FractalWaveList))
	for i, fw := range item.FractalWaveList {
		t := fw.WaveStart
		fractalWaves[i] = fxresponse.ZigZagFractalWave{WaveStart: &t, Wave: fw.Wave}
	}
	return fxresponse.ZigZagResult{
		Symbol:          item.Symbol,
		Depth:           item.Depth,
		Target4h:        toInfoSmaFib(item.Target4h),
		Current:         toInfoSmaFib(item.Current),
		Previous:        toInfo(item.Previous),
		Next:            toInfo(item.Next),
		Next2:           toInfo(item.Next2),
		NextRsRate:      item.NextRsRate,
		Next2MaxRate:    item.Next2MaxRate,
		WaveDxy4h:       item.WaveDxy4h,
		WaveDxy1h:       item.WaveDxy1h,
		FractalWaveList: fractalWaves,
	}
}

func toInfoSmaFib(src fxdto.ZigZagInfoSmaFibonacci) fxresponse.ZigZagInfoSmaFib {
	var fib *fxresponse.ZigZagFibonacci
	if src.Fibonacci != nil {
		f := fxresponse.ZigZagFibonacci{
			F1:         src.Fibonacci.F1,
			F7:         src.Fibonacci.F7,
			F6:         src.Fibonacci.F6,
			F5:         src.Fibonacci.F5,
			F3:         src.Fibonacci.F3,
			F2:         src.Fibonacci.F2,
			F0:         src.Fibonacci.F0,
			PriceRange: src.Fibonacci.PriceRange,
		}
		fib = &f
	}

	r := fxresponse.ZigZagInfoSmaFib{
		Wave:       src.Wave,
		Resistance: src.Resistance,
		Support:    src.Support,
		Fibonacci:  fib,
		Sma4h200s:  toSmaResponse(src.Sma4h200s),
		Sma4h75s:   toSmaResponse(src.Sma4h75s),
		Sma4h20s:   toSmaResponse(src.Sma4h20s),
		Sma1h200s:  toSmaResponse(src.Sma1h200s),
		Sma15m200s: toSmaResponse(src.Sma15m200s),
	}
	if !src.WaveStart.IsZero() {
		t := src.WaveStart
		r.WaveStart = &t
	}
	if !src.WaveEnd.IsZero() {
		t := src.WaveEnd
		r.WaveEnd = &t
	}
	return r
}

func toInfo(src fxdto.ZigZagInfo) fxresponse.ZigZagInfo {
	r := fxresponse.ZigZagInfo{
		Wave:       src.Wave,
		Resistance: src.Resistance,
		Support:    src.Support,
	}
	if !src.WaveStart.IsZero() {
		t := src.WaveStart
		r.WaveStart = &t
	}
	if !src.WaveEnd.IsZero() {
		t := src.WaveEnd
		r.WaveEnd = &t
	}
	return r
}

func toSmaResponse(src fxdto.ZigZagSma) fxresponse.ZigZagSma {
	return fxresponse.ZigZagSma{
		PriceS:    src.PriceS,
		PriceE:    src.PriceE,
		Deviation: src.Deviation,
		Fibonacci: src.Fibonacci,
		Direction: src.Direction,
		Position:  src.Position,
	}
}

func toStatusItem(s *zigzag.ZigZagStatus) fxresponse.ZigZagStatusItem {
	return fxresponse.ZigZagStatusItem{
		SymbolType:           s.SymbolType,
		BarType:              s.BarType,
		Symbol:               s.Symbol,
		Depth:                s.Depth,
		BarDateTimeMin:       s.BarDateTimeMin,
		BarDateTimeMax:       s.BarDateTimeMax,
		BarCount:             s.BarCount,
		BarDateTimeMinZigZag: s.BarDateTimeMinZigZag,
		BarDateTimeMaxZigZag: s.BarDateTimeMaxZigZag,
		ZigzagCount:          s.ZigzagCount,
		BreakResistanceCount: s.BreakResistanceCount,
		BreakSupportCount:    s.BreakSupportCount,
		Message:              s.Message,
	}
}
