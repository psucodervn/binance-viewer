package telegram

import (
	"bytes"
	"context"
	"copytrader/internal/binance"
	"copytrader/internal/user"
	"copytrader/internal/util"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"

	"copytrader/internal/model"
	"copytrader/internal/storage"
)

type Bot struct {
	bot             *telebot.Bot
	db              *model.Database
	markPriceFeeder binance.PriceFeeder
}

func NewBot(token string, db *model.Database, markPriceFeeder binance.PriceFeeder) *Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return &Bot{bot: bot, db: db, markPriceFeeder: markPriceFeeder}
}

func (b *Bot) Start() error {
	if err := b.bot.SetCommands(commands); err != nil {
		return err
	}

	// admin commands
	b.bot.Handle("/reload", cmdReload(b))

	// user commands
	b.bot.Handle("/start", cmdStart(b))
	b.bot.Handle("/add", cmdAddKey(b))
	b.bot.Handle("/info", cmdInfo(b))
	b.bot.Handle("/info_table", cmdInfoW(b))
	b.bot.Handle("/pnl", cmdPNL(b))
	b.bot.Handle("/card", cmdImageCard(b))
	b.bot.Handle("/positions", cmdPositions(b))

	// user query
	b.bot.Handle(&btnRefreshInfoTable, btnInfoW(b))

	b.bot.Start()
	return nil
}

func cmdReload(b *Bot) interface{} {
	return func(m *telebot.Message) {
		b.db.Renew()
		err := storage.LoadOrCreate(b.db)
		if err != nil {
			_, _ = b.bot.Send(m.Chat, "Reload failed: "+err.Error())
		} else {
			_, _ = b.bot.Send(m.Chat, "Reload succeed!")
		}
	}
}

func cmdStart(b *Bot) interface{} {
	return func(m *telebot.Message) {
		_, _ = b.bot.Send(m.Chat, "Welcome")
	}
}

func cmdAddKey(b *Bot) interface{} {
	return func(m *telebot.Message) {
		u := b.loadUser(m)
		args := splitArgs(m.Payload)
		if len(args) != 3 {
			_, _ = b.bot.Send(m.Chat, "Invalid command. Usage: /add name api_key secret_key")
			return
		}
		acc := model.NewAccount(args[0], args[1], args[2])
		if err := b.db.AddAccount(u, acc); err != nil {
			_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
			return
		}
		_ = storage.Save(b.db)
		_, _ = b.bot.Send(m.Chat, "Account added!")
	}
}

// loadUser load or create a new user
func (b *Bot) loadUser(msg *telebot.Message) model.User {
	uid := int64(msg.Sender.ID)
	if msg.Sender.IsBot {
		uid = msg.Chat.ID
	}
	u, err := b.db.FindUser(uid)
	if err == nil {
		return u
	}

	name := getSenderName(msg.Sender)
	if msg.Sender.IsBot {
		name = getChatName(msg.Chat)
	}
	u = model.NewUser(name, uid, nil)
	_ = b.db.AddUser(u)
	_ = storage.Save(b.db)
	return u
}

func (b *Bot) OnEvent(u model.User, acc model.Account, ev *futures.WsUserDataEvent) {
	if !u.TradeNotify {
		return
	}
	if ev.Event != futures.UserDataEventTypeOrderTradeUpdate {
		return
	}
	upd := ev.OrderTradeUpdate
	log.Debug().Str("user", u.Name).Str("account", acc.Name).
		Str("pnl", upd.RealizedPnL).
		Str("symbol", upd.Symbol).
		Str("fee", upd.Commission).
		Str("price", upd.AveragePrice).
		Str("type", string(upd.Type)).
		Str("status", string(upd.Status)).
		Str("side", string(upd.Side)).
		Str("position_side", string(upd.PositionSide)).
		Str("qty", upd.AccumulatedFilledQty).
		Send()
	log.Debug().Str("user", u.Name).Str("account", acc.Name).Interface("update", upd).Send()

	if upd.Status == futures.OrderStatusTypeNew {
		return
	}

	var bf bytes.Buffer
	switch true {
	case upd.PositionSide == futures.PositionSideTypeLong && upd.Side == futures.SideTypeBuy:
		bf.WriteRune('🛒')
		bf.WriteString(fmt.Sprintf(" [%s] Long <b>%s %s</b> at %.04f",
			acc.Name, upd.AccumulatedFilledQty, trimSymbol(upd.Symbol), util.ParseFloat(upd.AveragePrice),
		))
	case upd.PositionSide == futures.PositionSideTypeShort && upd.Side == futures.SideTypeSell:
		bf.WriteRune('🛒')
		bf.WriteString(fmt.Sprintf(" [%s] Short <b>%s %s</b> at %.04f 🔻",
			acc.Name, upd.AccumulatedFilledQty, trimSymbol(upd.Symbol), util.ParseFloat(upd.AveragePrice),
		))
	default:
		pnl := util.ParseFloat(upd.RealizedPnL)
		fee := util.ParseFloat(upd.Commission)
		fee = 0
		if pnl-fee >= 0 {
			bf.WriteRune('💚')
		} else {
			bf.WriteRune('💔')
		}
		bf.WriteString(fmt.Sprintf(" [%s] <b>%s</b> <b>%+.04f$</b> at %.04f",
			acc.Name, trimSymbol(upd.Symbol), pnl-fee, util.ParseFloat(upd.AveragePrice),
		))
	}

	_, _ = b.bot.Send(newTeleUser(u.TelegramID), bf.String(), telebot.ModeHTML)
}

func cmdPositions(b *Bot) interface{} {
	return func(m *telebot.Message) {
		u := b.loadUser(m)
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()

		var bf bytes.Buffer
		for _, a := range u.Accounts {
			cli := user.GetUserClient(a)
			info, err := cli.Info(ctx)
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}

			ps := b.filterPositions(info.Positions)
			total := 0.0
			for _, p := range ps {
				if p.UnrealizedProfit >= 0 {
					bf.WriteRune('🟩')
				} else {
					bf.WriteRune('🟥')
				}
				total += p.UnrealizedProfit
				percent := 100 * p.Leverage * (p.MarkPrice - p.EntryPrice) / p.EntryPrice
				if p.Side == model.SideShort {
					percent = -percent
				}
				bf.WriteString(fmt.Sprintf(" %s: <b>%+.02f</b>$ (%+.02f%%) | Price: %s → %s | Margin: %.02f$ x%.0f ",
					p.Symbol, p.UnrealizedProfit, percent,
					formatPrice(p.EntryPrice), formatPrice(p.MarkPrice), p.Margin, p.Leverage),
				)
				if p.Side == model.SideShort {
					bf.WriteString("🔻")
				}
				bf.WriteRune('\n')
				if len(ps) <= 30 {
					bf.WriteRune('\n')
				}
			}
			bf.WriteString(fmt.Sprintf("%s's Unrealized PNL: <b>%.02f</b>$ (%d positions)\n\n", a.Name, total, len(ps)))
		}
		_, _ = b.bot.Send(m.Chat, bf.String(), telebot.ModeHTML)
	}
}
