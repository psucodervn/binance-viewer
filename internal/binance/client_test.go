package binance

import (
	"context"
	"testing"
	"time"

	"copytrader/internal/model"
)

func TestIdolFollower_Follow(t *testing.T) {
	ctx := context.Background()
	c := NewIdolFollower()
	_ = c.Follow(ctx, model.IdolFmzcomAutoTrade)
	c.Start(ctx)
	time.Sleep(3 * time.Second)
}
