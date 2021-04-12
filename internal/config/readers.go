// Code generated by verigo. DO NOT EDIT.
package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// ReadBotConfig reads BotConfig from env
func ReadBotConfig(prefix ...string) (BotConfig, error) {
	p := ""
	if len(prefix) > 0 {
		p = prefix[0]
	}
	var cfg BotConfig
	err := envconfig.Process(p, &cfg)
	return cfg, err
}

// MustReadBotConfig reads BotConfig from env, panic if error
func MustReadBotConfig(prefix ...string) BotConfig {
	cfg, err := ReadBotConfig(prefix...)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return cfg
}
