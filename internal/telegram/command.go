package telegram

import (
	"gopkg.in/tucnak/telebot.v2"
)

var (
	commands = []telebot.Command{
		{Text: "info", Description: "View accounts info"},
		{Text: "info_table", Description: "View accounts info in table format"},
		{Text: "positions", Description: "View positions"},
		{Text: "card", Description: "Generate information card"},
		{Text: "pnl", Description: "Calculate today's PNL"},
		{Text: "add", Description: "Add api key: /add name api_key secret_key"},
	}
	btnRefreshInfoTable = new(telebot.ReplyMarkup).Data("Refresh", "refresh_info_table")
)
