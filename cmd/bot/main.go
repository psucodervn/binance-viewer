package bot

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"copytrader/internal/binance"
	"copytrader/internal/config"
	"copytrader/internal/model"
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
	db := model.NewDatabase()
	if err := storage.LoadOrCreate(db); err != nil {
		log.Fatal().Err(err).Send()
	}

	bot := telegram.NewBot(cfg.BotToken, db)
	if err := bot.Start(); err != nil {
		log.Fatal().Err(err).Send()
	}

	ctx := context.Background()
	ctx, cc := signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGINT)
	cli := binance.NewIdolFollower()
	// _ = cli.Follow(ctx, model.IdolFmzcomAutoTrade)
	_ = cli.Follow(ctx, model.IdolCryptoNifeCatch)
	_ = cli.Follow(ctx, model.IdolHuyLD)
	_ = cli.Follow(ctx, model.IdolHungLM)
	_ = cli.Follow(ctx, model.IdolPDYK)
	_ = cli.Follow(ctx, model.IdolHalfTalkVery)
	_ = cli.Follow(ctx, model.IdolCountyYearBy)
	_ = cli.Follow(ctx, model.IdolDegenerator)
	_ = cli.Follow(ctx, model.IdolGrugLikesRock)
	cli.Start(ctx)
	<-ctx.Done()
	cc()
}
