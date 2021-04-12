package config

type BotConfig struct {
	BotToken string `json:"token" split_words:"true" required:"true"`
}
