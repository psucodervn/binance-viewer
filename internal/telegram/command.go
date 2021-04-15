package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
)

var (
	commands = []telebot.Command{
		{Text: "i", Description: "View accounts info"},
		{Text: "it", Description: "View accounts info in table format"},
		{Text: "pnl", Description: "Calculate today's PNL"},
		{Text: "add", Description: "Add api key: /add name api_key secret_key"},
		{Text: "start", Description: "Start bot"},
	}
)
