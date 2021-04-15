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
	trades, err := c.cli.NewGetIncomeHistoryService().StartTime(from.UnixNano() / 1_000_000).EndTime(to.UnixNano() / 1_000_000).Limit(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
