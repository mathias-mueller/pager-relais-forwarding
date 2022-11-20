package main

import (
	"awesomeProject1/internal/activator"
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/gpio"
	"awesomeProject1/internal/telegram"
	"fmt"
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

	conf, err := config.Load()
	if err != nil {
		fmt.Printf("error loading config: %+v\n", err)
		os.Exit(1)
	}
	var level zerolog.Level
	level, err = zerolog.ParseLevel(conf.GeneralConfig.LogLevel)
	if err != nil {
		fmt.Printf("Could not parse log level '%s', using INFO", conf.GeneralConfig.LogLevel)
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

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
