package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
)

var (
	commands = []telebot.Command{
		{Text: "info", Description: "View accounts info"},
		{Text: "pnl", Description: "Calculate today's PNL"},
		{Text: "add", Description: "Add api key: /add api_key secret_key"},
		{Text: "start", Description: "Start bot"},
	}
)
