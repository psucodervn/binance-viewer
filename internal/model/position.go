package model

import (
	"math"

	"copytrader/internal/util"
)

const (
	SideBuy         = "BUY"
	SideSell        = "SELL"
	StatusNew       = "NEW"
	StatusIncreased = "INCREASED"
	StatusDecreased = "DECREASED"
	StatusClosed    = "CLOSED"
)

var (
	DefaultLeverage = 3.0
)

type Position struct {
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice"`
	MarkPrice       float64 `json:"markPrice"`
	PNL             float64 `json:"pnl"`
	ROE             float64 `json:"roe"`
	Amount          float64 `json:"amount"`
	UpdateTimeStamp int64   `json:"updateTimeStamp"`
}

func (p Position) Leverage() float64 {
	lv := math.Abs(p.ROE * p.EntryPrice / (p.MarkPrice - p.EntryPrice))
	if math.IsNaN(lv) || math.IsInf(lv, 0) {

	}
	lv = math.Max(lv, 1)
	return lv
}

func (p Position) Valid() bool {
	return len(p.Symbol) > 0
}

type PositionDiff struct {
	Symbol      string   `json:"symbol"`
	Amount      float64  `json:"amount"`
	Change      float64  `json:"change"`
	Status      string   `json:"status"`
	OldPosition Position `json:"oldPosition"`
	NewPosition Position `json:"newPosition"`
}

func (d PositionDiff) Equal(to PositionDiff) bool {
	return d.Symbol == to.Symbol &&
		util.IsEqual(d.Amount, to.Amount) &&
		util.IsEqual(d.Change, to.Change) &&
		d.Status == to.Status
}

func (d PositionDiff) Leverage() float64 {
	if d.NewPosition.Valid() {
		return d.NewPosition.Leverage()
	} else if d.OldPosition.Valid() {
		return d.OldPosition.Leverage()
	} else {
		return DefaultLeverage
	}
}

func (d PositionDiff) PNL() float64 {
	if d.NewPosition.Valid() {
		return d.NewPosition.PNL
	} else {
		return d.OldPosition.PNL
	}
}

func DiffPositions(oldPositions, newPositions []Position) []PositionDiff {

	pos := map[string]Position{}
	for _, p := range oldPositions {
		pos[p.Symbol] = p
	}

	// updated positions
	var res []PositionDiff
	for _, p := range newPositions {
		if _, exists := pos[p.Symbol]; !exists {

			res = append(res, PositionDiff{
				Symbol:      p.Symbol,
				Amount:      p.Amount,
				Change:      p.Amount,
				Status:      StatusNew,
				NewPosition: p,
			})
			continue
		}

		old := pos[p.Symbol]
		if util.IsEqual(old.Amount, p.Amount) {

			delete(pos, p.Symbol)
			continue
		}

		d := PositionDiff{
			Symbol:      p.Symbol,
			Amount:      p.Amount,
			Change:      p.Amount - old.Amount,
			OldPosition: old,
			NewPosition: p,
		}
		switch {
		case util.IsZero(d.Amount):
			d.Status = StatusClosed
		case d.Change > 0 && d.Amount > 0:
			d.Status = StatusIncreased
		case d.Change < 0 && d.Amount > 0:
			d.Status = StatusDecreased
		case d.Change < 0 && d.Amount < 0:
			d.Status = StatusIncreased
		case d.Change > 0 && d.Amount < 0:
			d.Status = StatusIncreased
		}
		res = append(res, d)

		delete(pos, p.Symbol)
	}

	for _, p := range pos {
		res = append(res, PositionDiff{
			Symbol:      p.Symbol,
			Amount:      0,
			Change:      -p.Amount,
			Status:      StatusClosed,
			OldPosition: p,
		})
	}

	return res
}
