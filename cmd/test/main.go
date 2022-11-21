package main

import (
	"awesomeProject1/internal/activator"
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/telegram"
	"net/http"
	"os"
	"time"

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
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              ":2112",
		ReadHeaderTimeout: time.Second,
		Handler:           handler,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Err(err).Msg("Server failed")
	}
}
