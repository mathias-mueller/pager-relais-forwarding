package main

import (
	"awesomeProject1/internal/activator"
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/telegram"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("failed to load config")
		os.Exit(1)
	}
	a := activator.New()
	telegramAPI := telegram.Init(conf.TelegramConfig)

	inputs := make(chan bool)
	defer close(inputs)

	go a.EnableActivation(inputs,
		[]activator.Activation{
			&activator.TelegramActivation{API: telegramAPI},
		},
	)

	inputs <- true
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
