package telegram

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/nleeper/goment"
	"github.com/rs/zerolog/log"
	"gopkg.in/tucnak/telebot.v2"

	"copytrader/internal/model"
	"copytrader/internal/storage"
	"copytrader/internal/user"
	"copytrader/internal/util"
)

type Bot struct {
	bot *telebot.Bot
	db  *model.Database
}

func NewBot(token string, db *model.Database) *Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return &Bot{bot: bot, db: db}
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
	b.bot.Handle("/pnl", cmdPNL(b))
	b.bot.Start()
	return nil
}

func cmdPNL(b *Bot) interface{} {
	return func(m *telebot.Message) {
		u := b.loadUser(m.Chat)
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()

		args := splitArgs(m.Payload)
		diff := 0
		if len(args) >= 1 {
			diff = int(util.ParseInt(args[0], 0))
		}

		from, _ := goment.New()
		from.Subtract(diff, "day").StartOf("day").SetHour(7)
		to, _ := goment.New(from)
		to = to.Add(1, "day").Subtract(1, "second")
		log.Debug().Str("from", from.ToString()).Str("to", to.ToString()).Send()

		var bf bytes.Buffer
		total := 0.0
		for _, a := range u.Accounts {
			cli := user.GetUserClient(a)
			incomes, err := cli.ListIncome(ctx, from.ToTime(), to.ToTime())
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}
			pnl := user.TotalUserIncome(incomes)
			total += pnl
			bf.WriteString(fmt.Sprintf("[+] %s: %.02f$\n", a.Name, pnl))
		}
		bf.WriteString(fmt.Sprintf("Total: %.02f$", total))
		_, _ = b.bot.Send(m.Chat, bf.String())
	}
}

func cmdInfo(b *Bot) interface{} {
	return func(m *telebot.Message) {
		u := b.loadUser(m.Chat)
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()

		var bf bytes.Buffer
		for _, a := range u.Accounts {
			cli := user.GetUserClient(a)
			acc, err := cli.Info(ctx)
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}
			bf.WriteString(fmt.Sprintf("[+] %s:", a.Name))
			bf.WriteString(fmt.Sprintf(" UnP: %.02f", util.ParseFloat(acc.TotalUnrealizedProfit)))
			bf.WriteString(fmt.Sprintf(" | AvB: %.02f", util.ParseFloat(acc.MaxWithdrawAmount)))
			bf.WriteString(fmt.Sprintf(" | MgB: %.02f", util.ParseFloat(acc.TotalMarginBalance)))
			bf.WriteString(fmt.Sprintf(" | WaB: %.02f", util.ParseFloat(acc.TotalWalletBalance)))
			bf.WriteString("\n")
		}
		_, _ = b.bot.Send(m.Chat, bf.String())
	}
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
		u := b.loadUser(m.Chat)
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
func (b *Bot) loadUser(chat *telebot.Chat) model.User {
	u, err := b.db.FindUser(chat.ID)
	if err == nil {
		return u
	}
	u = model.NewUser(getName(chat), chat.ID, nil)
	_ = b.db.AddUser(u)
	_ = storage.Save(b.db)
	return u
}
