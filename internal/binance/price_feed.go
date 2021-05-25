package binance

import (
	"context"
	"copytrader/internal/model"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type OnNewPrice func(symbol string, price float64)

type PriceFeeder interface {
	Get(symbol string) float64
	Subscribe(fn OnNewPrice)
}

type MarkPriceFeeder struct {
	wsBaseURL string
	prices    *PriceMap
}

const (
	wsBaseURL = "wss://fstream.binance.com/ws/"
)

type rawMarkPriceUpdateEvent struct {
	Type                 string `json:"e"`
	Time                 int64  `json:"E"`
	Symbol               string `json:"s"`
	MarkPrice            string `json:"p"`
	IndexPrice           string `json:"i"`
	EstimatedSettlePrice string `json:"P"`
}

func NewMarkPriceFeeder() *MarkPriceFeeder {
	return &MarkPriceFeeder{
		wsBaseURL: wsBaseURL,
		prices:    NewPriceMap(),
	}
}

func (m *MarkPriceFeeder) Start(ctx context.Context) error {
	cfg := &WsConfig{
		Endpoint:  m.wsBaseURL + "!markPrice@arr",
		KeepAlive: true,
		Timeout:   time.Second * 10,
	}

	errHandler := func(err error) {
		log.Err(err).Msg("ws failed")
	}
	eventHandler := func(ev model.MarkPriceUpdateEvent) {
		m.prices.Set(ev.Symbol, ev.MarkPrice)
	}
	wsHandler := func(rawMessage []byte) {
		var rawEvents []rawMarkPriceUpdateEvent
		if err := json.Unmarshal(rawMessage, &rawEvents); err != nil {
			errHandler(err)
			return
		}
		if len(rawEvents) == 0 {
			return
		}
		for _, rawEv := range rawEvents {
			ev := model.MarkPriceUpdateEvent{
				Symbol: rawEv.Symbol,
				Time:   rawEv.Time,
				Type:   rawEv.Type,
			}
			ev.MarkPrice, _ = strconv.ParseFloat(rawEv.MarkPrice, 64)
			eventHandler(ev)
		}
	}

	doneC, stopC, err := wsServe(cfg, wsHandler, errHandler)
	go func() {
		select {
		case <-ctx.Done():
			stopC <- struct{}{}
			return
		case <-doneC:
			return
		}
	}()

	return err
}

func (m *MarkPriceFeeder) Get(symbol string) float64 {
	return m.prices.Get(symbol)
}

func (m *MarkPriceFeeder) Subscribe(fn OnNewPrice) {
	panic("implement me")
}
