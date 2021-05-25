package bot

import (
	"context"
	"copytrader/internal/binance"
	"os/signal"
	"syscall"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"copytrader/internal/config"
	"copytrader/internal/model"
	"copytrader/internal/runner"
	"copytrader/internal/storage"
	"copytrader/internal/telegram"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bot",
		Short: "bot description",
		Run:   run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	cfg := config.MustReadBotConfig()
	if cfg.UseTestnet {
		futures.UseTestnet = true
		storage.DefaultFilePath = "data/data.testnet.json"
	}

	db := model.NewDatabase()
	if err := storage.LoadOrCreate(db); err != nil {
		log.Fatal().Err(err).Send()
	}

	m := runner.NewManager(db)
	m.Start()

	ctx := context.Background()
	ctx, cc := signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGINT)

	feeder := binance.NewMarkPriceFeeder()
	_ = feeder.Start(ctx)

	bot := telegram.NewBot(cfg.BotToken, db, feeder)
	go func() {
		if err := bot.Start(); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()
	cc()
}

// cli := binance.NewIdolFollower()
// _ = cli.Follow(ctx, model.IdolFmzcomAutoTrade)
// _ = cli.Follow(ctx, model.IdolCryptoNifeCatch)
// _ = cli.Follow(ctx, model.IdolHuyLD)
// _ = cli.Follow(ctx, model.IdolHungLM)
// _ = cli.Follow(ctx, model.IdolPDYK)
// _ = cli.Follow(ctx, model.IdolHalfTalkVery)
// _ = cli.Follow(ctx, model.IdolCountyYearBy)
// _ = cli.Follow(ctx, model.IdolDegenerator)
// _ = cli.Follow(ctx, model.IdolGrugLikesRock)
// cli.Start(ctx)
