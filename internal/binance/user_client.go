package binance

import (
	"context"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

type UserClient struct {
	ApiKey    string
	SecretKey string

	cli *futures.Client
}

func NewUserClient(apiKey string, secretKey string) *UserClient {
	cli := binance.NewFuturesClient(apiKey, secretKey)
	return &UserClient{
		ApiKey:    apiKey,
		SecretKey: secretKey,
		cli:       cli,
	}
}

func (c *UserClient) Info(ctx context.Context) (acc *futures.Account, err error) {
	acc, err = c.cli.NewGetAccountService().Do(ctx)
	return
}

func (c *UserClient) Trades(ctx context.Context) (trades []*futures.AccountTrade, err error) {
	trades, err = c.cli.NewListAccountTradeService().Limit(20).Do(ctx)
	return
}

func (c *UserClient) HistoryTrades(ctx context.Context) (trades []*futures.Trade, err error) {
	return c.cli.NewHistoricalTradesService().Symbol("BATUSDT").Limit(20).Do(ctx)
}

func (c *UserClient) ListIncome(ctx context.Context, from, to time.Time) (incomes []*futures.IncomeHistory, err error) {
	trades, err := c.cli.NewGetIncomeHistoryService().StartTime(from.UnixNano() / 1_000_000).EndTime(to.UnixNano() / 1_000_000).Limit(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	return trades, nil
}
