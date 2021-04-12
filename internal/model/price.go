package model

type MarkPriceUpdateEvent struct {
	Type                 string  `json:"e"`
	Time                 int64   `json:"E"`
	Symbol               string  `json:"s"`
	MarkPrice            float64 `json:"p"`
	IndexPrice           float64 `json:"i"`
	EstimatedSettlePrice float64 `json:"P"`
}
