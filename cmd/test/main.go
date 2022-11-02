package main

import (
	"awesomeProject1/internal/config"
	"awesomeProject1/internal/gpio"
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

	actuations := gpio.Start(conf.GpioConfig)

	for actuation := range actuations {
		log.Info().
			Time("time", actuation).
			Msg("Actuation")
	}
}
