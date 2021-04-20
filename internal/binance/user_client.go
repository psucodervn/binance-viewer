package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type AccountClient struct {
	ApiKey    string
	SecretKey string

	cli *futures.Client
}

func NewAccountClient(apiKey string, secretKey string) *AccountClient {
	cli := binance.NewFuturesClient(apiKey, secretKey)
	return &AccountClient{
		ApiKey:    apiKey,
		SecretKey: secretKey,
		cli:       cli,
	}
}

func (c *AccountClient) Info(ctx context.Context) (acc *futures.Account, err error) {
	acc, err = c.cli.NewGetAccountService().Do(ctx)
	return
}

func (c *AccountClient) Trades(ctx context.Context) (trades []*futures.AccountTrade, err error) {
	trades, err = c.cli.NewListAccountTradeService().Limit(20).Do(ctx)
	return
}

func (c *AccountClient) HistoryTrades(ctx context.Context) (trades []*futures.Trade, err error) {
	return c.cli.NewHistoricalTradesService().Symbol("BATUSDT").Limit(20).Do(ctx)
}

func (c *AccountClient) ListIncome(ctx context.Context, from, to time.Time) (incomes []*futures.IncomeHistory, err error) {
	const Limit = 1000
	var (
		start  = from.UnixNano() / 1_000_000
		end    = to.UnixNano() / 1_000_000
		lastID = int64(0)
	)
	for {
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		trades, err := c.cli.NewGetIncomeHistoryService().StartTime(start).EndTime(end).Limit(Limit).Do(ctx)
		cc()
		if err != nil {
			return incomes, err
		}
		if len(trades) == 0 {
			break
		}
		i := 0
		for i < len(trades) && trades[i].TranID == lastID {
			i++
		}
		incomes = append(incomes, trades[i:]...)
		if len(trades) < Limit {
			break
		}
		lastID = trades[len(trades)-1].TranID
		start = trades[len(trades)-1].Time
	}
	return incomes, nil
}
