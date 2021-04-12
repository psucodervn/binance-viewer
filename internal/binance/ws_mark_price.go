package binance

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"copytrader/internal/model"
)

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

func (c *IdolFollower) startMarkPriceListener(ctx context.Context) {
	cfg := &WsConfig{
		Endpoint:  wsBaseURL + "!markPrice@arr",
		KeepAlive: true,
		Timeout:   time.Second * 10,
	}
	errHandler := func(err error) {
		log.Err(err).Msg("ws failed")
	}
	eventHandler := func(ev model.MarkPriceUpdateEvent) {
		// log.Debug().Str("symbol", ev.Symbol).Float64("price", ev.MarkPrice).Send()
		c.markPrices.Set(ev.Symbol, ev.MarkPrice)
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
	_, _, _ = wsServe(cfg, wsHandler, errHandler)
}
