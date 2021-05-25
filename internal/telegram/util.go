package telegram

import (
	"copytrader/internal/model"
	"copytrader/internal/util"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

var re = regexp.MustCompile(`\s+`)

func splitArgs(payload string) []string {
	return re.Split(strings.TrimSpace(payload), -1)
}

func getSenderName(m *telebot.User) string {
	name := strings.TrimSpace(m.FirstName + " " + m.LastName)
	if len(name) == 0 {
		name = strings.TrimSpace(m.Username)
	}
	return name
}

func getChatName(m *telebot.Chat) string {
	name := strings.TrimSpace(m.FirstName + " " + m.LastName)
	if len(name) == 0 {
		name = strings.TrimSpace(m.Username)
	}
	return name
}

func formatPrice(p float64) string {
	if p > 1_000 {
		return fmt.Sprintf("%.0f", p)
	} else if p > 100 {
		return fmt.Sprintf("%.1f", p)
	} else if p > 10 {
		return fmt.Sprintf("%.2f", p)
	} else {
		return fmt.Sprintf("%.4f", p)
	}
}

func trimSymbol(symbol string) string {
	if strings.HasSuffix(symbol, "USDT") {
		symbol = symbol[:len(symbol)-4]
	}
	return symbol
}

type Position struct {
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`
	EntryPrice       float64 `json:"entryPrice"`
	MarkPrice        float64 `json:"markPrice"`
	Leverage         float64 `json:"leverage"`
	UnrealizedProfit float64 `json:"unrealizedProfit"`
	Margin           float64 `json:"margin"`
}

func (b *Bot) filterPositions(positions []*futures.AccountPosition) []Position {
	var ps []Position
	for i := range positions {
		p := positions[i]
		if !util.IsZero(util.ParseFloat(p.PositionAmt)) {
			side := model.SideLong
			if p.PositionSide == "SHORT" {
				side = model.SideShort
			}
			symbol := p.Symbol
			if strings.HasSuffix(symbol, "USDT") {
				symbol = symbol[:len(symbol)-4]
			}

			fp := Position{
				Symbol:           symbol,
				Side:             side,
				EntryPrice:       util.ParseFloat(p.EntryPrice),
				MarkPrice:        b.markPriceFeeder.Get(p.Symbol),
				Leverage:         util.ParseFloat(p.Leverage),
				UnrealizedProfit: util.ParseFloat(p.UnrealizedProfit),
				Margin:           util.ParseFloat(p.InitialMargin),
			}

			ps = append(ps, fp)
		}
	}
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].UnrealizedProfit > ps[j].UnrealizedProfit
	})
	return ps
}

type teleUser struct {
	ID int64
}

func (t teleUser) Recipient() string {
	return strconv.Itoa(int(t.ID))
}

func newTeleUser(id int64) teleUser {
	return teleUser{ID: id}
}
