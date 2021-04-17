package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
)

type AccountRunner struct {
	cli *futures.Client
}

func NewAccountRunner(apiKey string, secretKey string) *AccountRunner {
	cli := futures.NewClient(apiKey, secretKey)
	return &AccountRunner{cli: cli}
}

func (r *AccountRunner) OnUpdate() error {
	ctx := context.Background()
	listenKey, err := r.cli.NewStartUserStreamService().Do(ctx)
	if err != nil {
		return err
	}

	go func() {
		for range time.Tick(10 * time.Minute) {
			_ = r.cli.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(ctx)
		}
	}()

	go func() {
		for {
			doneC, stopC, err := futures.WsUserDataServe(listenKey, func(ev *futures.WsUserDataEvent) {
				fmt.Printf("%+v\n", ev)
			}, func(err error) {
				log.Err(err).Send()
			})
			if err != nil {
				log.Err(err).Send()
				return
			}
			select {
			case <-doneC:
				log.Debug().Msg("done")
				continue
			case <-stopC:
				log.Debug().Msg("stop")
				return
			}
		}
	}()

	return err
}
