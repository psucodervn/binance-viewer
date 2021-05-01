package telegram

import (
	"regexp"
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

var re = regexp.MustCompile(`\s+`)

func splitArgs(payload string) []string {
	return re.Split(strings.TrimSpace(payload), -1)
}

func getSenderName(m *telebot.User) string {
	name := strings.TrimSpace(m.FirstName + " " + m.LastName)
	if len(name) == 0 {
		name = strings.TrimSpace(m.Username)
	}
	return name
}

func getChatName(m *telebot.Chat) string {
	name := strings.TrimSpace(m.FirstName + " " + m.LastName)
	if len(name) == 0 {
		name = strings.TrimSpace(m.Username)
	}
	return name
}
