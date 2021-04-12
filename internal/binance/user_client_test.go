package binance

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newUserClientFromEnv() *UserClient {
	return NewUserClient(os.Getenv("TEST_API_KEY"), os.Getenv("TEST_SECRET_KEY"))
}

func TestUserClient_Info(t *testing.T) {
	type fields struct {
		ApiKey    string
		SecretKey string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{fields: fields{ApiKey: os.Getenv("TEST_API_KEY"), SecretKey: os.Getenv("TEST_SECRET_KEY")}, wantErr: false},
		{fields: fields{ApiKey: "fake", SecretKey: "fake"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewUserClient(tt.fields.ApiKey, tt.fields.SecretKey)
			ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFunc()
			gotAcc, err := c.Info(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.NotNil(t, gotAcc)
			}
		})
	}
}

func TestUserClient_Trades(t *testing.T) {
	c := newUserClientFromEnv()
	ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cc()
	gotTrades, err := c.Trades(ctx)
	if err != nil {
		t.Errorf("Trades() error = %v", err)
		return
	}
	require.NotNil(t, t, gotTrades)
}

func TestUserClient_HistoryTrades(t *testing.T) {
	c := newUserClientFromEnv()
	ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cc()
	gotTrades, err := c.HistoryTrades(ctx)
	if err != nil {
		t.Errorf("HistoryTrades() error = %v", err)
		return
	}
	require.NotNil(t, t, gotTrades)
}

func TestUserClient_ListIncome(t *testing.T) {
	c := newUserClientFromEnv()
	ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cc()

	gotTrades, err := c.ListIncome(ctx, time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		t.Errorf("ListIncome() error = %v", err)
		return
	}
	require.NotNil(t, t, gotTrades)
}
