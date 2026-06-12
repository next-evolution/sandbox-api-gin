package fx

type TradePosition struct {
	ID              *int64  `json:"id"`
	PositionNumber  int16   `json:"positionNumber"`
	SettlementPrice float64 `json:"settlementPrice"`
	SettlementPips  int     `json:"settlementPips"`
	SettlementRatio float64 `json:"settlementRatio"`
	Lot             float64 `json:"lot"`
	ProfitAmount    int     `json:"profitAmount"`
	LossAmount      int     `json:"lossAmount"`
}
