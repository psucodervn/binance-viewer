package user

import (
	"github.com/adshao/go-binance/v2/futures"
)

type AccountRunner struct {
	cli *futures.Client
}

func NewAccountRunner(apiKey string, secretKey string) *AccountRunner {
	cli := futures.NewClient(apiKey, secretKey)
	return &AccountRunner{cli: cli}
}

func (r *AccountRunner) OnUpdate() error {
	// ctx := context.Background()
	// listenKey, err := r.cli.NewStartUserStreamService().Do(ctx)
	// if err != nil {
	//   return err
	// }
	return nil
}
