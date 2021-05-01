package telegram

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/nleeper/goment"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/tucnak/telebot.v2"

	"copytrader/internal/model"
	"copytrader/internal/user"
	"copytrader/internal/util"
)

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
		handleInfoW(b, m)
	}
}
func btnInfoW(b *Bot) interface{} {
	return func(cb *telebot.Callback) {
		handleInfoW(b, cb.Message)
	}
}

func handleInfoW(b *Bot, m *telebot.Message) {
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
		fmt.Fprintf(w, "%.02f (%d) \t", unp, user.TotalPosition(acc.Positions))
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

	_, err := b.bot.Edit(m, "```"+bf.String()+"```", telebot.ModeMarkdownV2, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*btnRefreshInfoTable.Inline()},
		},
	})
	if err == telebot.ErrCantEditMessage {
		_, err = b.bot.Send(m.Chat, "```"+bf.String()+"```", telebot.ModeMarkdownV2, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{
				{*btnRefreshInfoTable.Inline()},
			},
		})
	}
	if err != nil {
		log.Debug().Err(err).Send()
	}
}

type accountInfo struct {
	Info    model.Account
	Account *futures.Account
	Incomes []*futures.IncomeHistory
}

func buildInfoTable(accounts []accountInfo) Table {
	headers := Headers{
		Header{Text: "Name"},
		Header{Text: "Daily PNL", Align: AlignRight},
		Header{Text: "Un. PNL", Align: AlignRight},
		Header{Text: "Avail.", Align: AlignRight},
		// Header{Text: "WaB", Align: AlignRight},
		Header{Text: "Balance", Align: AlignRight},
	}

	totalP, totalUnP, totalAvB, totalMgB, totalWaB := 0.0, 0.0, 0.0, 0.0, 0.0
	var rows Rows
	for _, a := range accounts {
		acc := a.Account
		pnl := user.TotalUserIncome(a.Incomes)
		totalP += pnl
		unp := util.ParseFloat(acc.TotalUnrealizedProfit)
		totalUnP += unp
		avb := util.ParseFloat(acc.MaxWithdrawAmount)
		totalAvB += avb
		wab := util.ParseFloat(acc.TotalWalletBalance)
		totalWaB += wab
		mgb := util.ParseFloat(acc.TotalMarginBalance)
		totalMgB += mgb
		columns := Columns{
			Column{Data: a.Info.Name},
			Column{Data: formatMoney(pnl), Color: profitColor(pnl)},
			Column{Data: formatMoney(unp), Color: profitColor(unp)},
			Column{Data: formatMoney(avb), Color: profitColor(avb)},
			// Column{Data: formatMoney(wab), Color: profitColor(wab)},
			Column{Data: formatMoney(mgb), Color: profitColor(mgb)},
		}
		rows = append(rows, Row{Columns: columns})
	}

	footers := Footers{
		Column{Data: "Total"},
		Column{Data: formatMoney(totalP), Color: profitColor(totalP)},
		Column{Data: formatMoney(totalUnP), Color: profitColor(totalUnP)},
		Column{Data: formatMoney(totalAvB), Color: profitColor(totalAvB)},
		// Column{Data: formatMoney(totalWaB), Color: profitColor(totalWaB)},
		Column{Data: formatMoney(totalMgB), Color: profitColor(totalMgB)},
	}

	return Table{
		Headers: headers,
		Rows:    rows,
		Footers: footers,
	}
}

var (
	mp = message.NewPrinter(language.English)
)

func formatMoney(v float64) string {
	return mp.Sprintf("%.01f", v)
}
