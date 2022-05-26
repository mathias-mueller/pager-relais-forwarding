package main

import (
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/telegram"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out: os.Stdout,
		},
	)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	conf, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	telegram.Init(conf.TelegramConfig)

	telegram.SendMsgString("Hello World")
}
