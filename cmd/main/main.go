package main

import (
	"awesomeProject1/internal/activator"
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/gpio"
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

	conf, err := config.Load()
	if err != nil {
		log.Err(err).Msg("Could not load config")
		os.Exit(1)
	}
	var level zerolog.Level
	level, err = zerolog.ParseLevel(conf.GeneralConfig.LogLevel)
	if err != nil {
		log.Warn().
			Err(err).
			Str("wanted", conf.GeneralConfig.LogLevel).
			Msg("Could not parse log level, using INFO")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	a := activator.New()
	telegramAPI := telegram.Init(conf.TelegramConfig)

	rawGpioValues := gpio.Start(conf.GpioConfig)

	go func() {
		a.EnableActivation(rawGpioValues,
			[]activator.Activation{
				&activator.TelegramActivation{API: telegramAPI},
			},
		)
		log.Fatal().Msg("Activation ended")
	}()

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
		os.Exit(1)
	}
}
