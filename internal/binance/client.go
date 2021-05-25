package binance

import (
	"container/list"
	"context"
	"copytrader/internal/util"
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"

	"copytrader/internal/model"
)

type IdolFollower struct {
	r               *rate.Limiter
	idols           *list.List
	mu              sync.RWMutex
	positions       map[string][]model.Position
	pnls            map[string]*atomic.Float64
	markPriceFeeder PriceFeeder
}

func NewIdolFollower(markPriceFeeder PriceFeeder) *IdolFollower {
	return &IdolFollower{
		r:               rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
		idols:           list.New(),
		positions:       map[string][]model.Position{},
		pnls:            map[string]*atomic.Float64{},
		markPriceFeeder: markPriceFeeder,
	}
}

func (c *IdolFollower) Start(ctx context.Context) {
	go c.startFollow(ctx)
}

func (c *IdolFollower) Follow(ctx context.Context, idol model.Idol) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pnls[idol.UID] = atomic.NewFloat64(0)
	for it := c.idols.Front(); it != nil; it = it.Next() {
		id := it.Value.(model.Idol)
		if id.UID == idol.UID {
			// TODO: duplicate
			return nil
		}
	}
	c.idols.PushBack(idol)
	return nil
}

func (c *IdolFollower) Unfollow(ctx context.Context, idol model.Idol) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for it := c.idols.Front(); it != nil; it = it.Next() {
		id := it.Value.(model.Idol)
		if id.UID != idol.UID {
			continue
		}
		c.idols.Remove(it)
		return nil
	}
	return nil
}

func (c *IdolFollower) startFollow(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.mu.RLock()
			if c.idols.Len() == 0 {
				c.mu.RUnlock()
				time.Sleep(100 * time.Millisecond)
				continue
			}
			e := c.idols.Front()
			c.mu.RUnlock()

			id := e.Value.(model.Idol)
			_ = c.r.Wait(ctx)
			// log.Debug().Interface("idol", id).Time("now", time.Now()).Send()
			if err := c.process(ctx, id); err != nil {
				log.Err(err).Interface("idol", id).Msg("process failed")
			}

			c.mu.Lock()
			c.idols.MoveToBack(e)
			c.mu.Unlock()

			// time.Sleep(1 * time.Second)
		}
	}
}

func (c *IdolFollower) process(ctx context.Context, idol model.Idol) error {
	c.mu.RLock()
	oldPositions := c.positions[idol.UID]
	c.mu.RUnlock()

	resp, err := GetIdolPositions(ctx, idol.UID)
	if err != nil {
		return err
	}
	// if oldPositions != nil {
	diff := model.DiffPositions(oldPositions, resp.Positions)
	if len(diff) > 0 {
		printDiffs(idol, diff)
		for _, d := range diff {
			if d.Status == model.StatusClosed {
				// p := c.markPrices.GetOrStore(d.Symbol, d.OldPosition.MarkPrice)
				p := c.markPriceFeeder.Get(d.Symbol)
				if util.IsZero(p) {
					p = d.OldPosition.MarkPrice
				}
				pnl := (p - d.OldPosition.EntryPrice) * d.OldPosition.Amount
				c.pnls[idol.UID].Add(pnl)
			}
		}
		fmt.Printf("PNL: %.2f\n", c.pnls[idol.UID].Load())
	}
	// }

	c.mu.Lock()
	c.positions[idol.UID] = resp.Positions
	c.mu.Unlock()
	return nil
}

func printDiffs(idol model.Idol, diff []model.PositionDiff) {
	log.Info().Str("now", time.Now().Format("15:04:05.9999")).Msg(idol.Name + " updates")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, d := range diff {
		_, _ = fmt.Fprintf(w, "%v\t%v\t%.4f\t%0.2f\t%.2f\n", d.Symbol, d.Status, d.Change, d.Leverage(), d.PNL())
	}
	_ = w.Flush()
}
