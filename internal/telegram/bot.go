package telegram

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"
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
	b.bot.Handle("/i", cmdInfo(b))
	b.bot.Handle("/ii", cmdInfoW(b))
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
		from.UTC().Subtract(diff, "day").StartOf("day")
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
		totalUnP, totalAvB, totalMgB, totalWaB := 0.0, 0.0, 0.0, 0.0
		for _, a := range u.Accounts {
			cli := user.GetUserClient(a)
			acc, err := cli.Info(ctx)
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}
			bf.WriteString(fmt.Sprintf("[+] %s:", a.Name))
			unp := util.ParseFloat(acc.TotalUnrealizedProfit)
			totalUnP += unp
			bf.WriteString(fmt.Sprintf(" UnP: %.02f", unp))
			avb := util.ParseFloat(acc.MaxWithdrawAmount)
			totalAvB += avb
			bf.WriteString(fmt.Sprintf(" | AvB: %.02f", avb))
			mgb := util.ParseFloat(acc.TotalMarginBalance)
			totalMgB += mgb
			bf.WriteString(fmt.Sprintf(" | MgB: %.02f", mgb))
			wab := util.ParseFloat(acc.TotalWalletBalance)
			totalWaB += wab
			bf.WriteString(fmt.Sprintf(" | WaB: %.02f", wab))
			bf.WriteString("\n")
		}

		// total
		bf.WriteString(fmt.Sprintf("Total:"))
		bf.WriteString(fmt.Sprintf(" UnP: %.02f", totalUnP))
		bf.WriteString(fmt.Sprintf(" | AvB: %.02f", totalAvB))
		bf.WriteString(fmt.Sprintf(" | WaB: %.02f", totalWaB))
		bf.WriteString(fmt.Sprintf(" | MgB: %.02f", totalMgB))
		bf.WriteString("\n")
		_, _ = b.bot.Send(m.Chat, bf.String())
	}
}

func cmdInfoW(b *Bot) interface{} {
	return func(m *telebot.Message) {
		u := b.loadUser(m.Chat)
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()

		var bf bytes.Buffer
		totalP, totalUnP, totalAvB, totalMgB, totalWaB := 0.0, 0.0, 0.0, 0.0, 0.0
		from, _ := goment.New()
		from.UTC().StartOf("day")
		to, _ := goment.New(from)
		to = to.Add(1, "day").Subtract(1, "second")

		w := tabwriter.NewWriter(&bf, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(w, "Name \tPnL \tUnP \tAvB \tWaB \tMgB \t")
		fmt.Fprintln(w, "------ \t----- \t----- \t----- \t----- \t----- \t")
		for _, a := range u.Accounts {
			cli := user.GetUserClient(a)
			acc, err := cli.Info(ctx)
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}
			incomes, err := cli.ListIncome(ctx, from.ToTime(), to.ToTime())
			if err != nil {
				_, _ = b.bot.Send(m.Chat, "Error: "+err.Error())
				continue
			}
			pnl := user.TotalUserIncome(incomes)
			totalP += pnl

			fmt.Fprintf(w, "%s \t", a.Name)
			unp := util.ParseFloat(acc.TotalUnrealizedProfit)
			totalUnP += unp
			fmt.Fprintf(w, "%.02f \t", pnl)
			fmt.Fprintf(w, "%.02f \t", unp)
			avb := util.ParseFloat(acc.MaxWithdrawAmount)
			totalAvB += avb
			fmt.Fprintf(w, "%.02f \t", avb)
			wab := util.ParseFloat(acc.TotalWalletBalance)
			totalWaB += wab
			fmt.Fprintf(w, "%.02f \t", wab)
			mgb := util.ParseFloat(acc.TotalMarginBalance)
			totalMgB += mgb
			fmt.Fprintf(w, "%.02f \t", mgb)
			fmt.Fprintln(w)
		}

		// total
		fmt.Fprintln(w, "------ \t----- \t----- \t----- \t----- \t----- \t")
		fmt.Fprintf(w, "Total \t")
		fmt.Fprintf(w, "%.02f \t", totalP)
		fmt.Fprintf(w, "%.02f \t", totalUnP)
		fmt.Fprintf(w, "%.02f \t", totalAvB)
		fmt.Fprintf(w, "%.02f \t", totalWaB)
		fmt.Fprintf(w, "%.02f \t", totalMgB)
		fmt.Fprintln(w)
		_ = w.Flush()

		_, err := b.bot.Send(m.Chat, "```"+bf.String()+"```", telebot.ModeMarkdownV2)
		if err != nil {
			log.Debug().Err(err).Send()
		}
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
