package gpio

import (
	"awesomeProject1/internal/config"

	"time"

	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

func Start(config *config.GpioConfig) <-chan bool {
	err := rpio.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init GPIO library")
	}

	output := make(chan bool)

	ticker := time.NewTicker(time.Duration(config.Interval) * time.Millisecond)
	done := make(chan time.Time)
	go func() {
		defer close(output)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				go func() {
					output <- IsPinHigh(config)
				}()
			}
		}
	}()
	return output
}

func IsPinHigh(config *config.GpioConfig) bool {
	pinNumber := config.Pin
	log.Trace().
		Int("number", pinNumber).
		Msg("Reading pin")
	pin := rpio.Pin(pinNumber)
	pin.Input()

	state := pin.Read()
	isHigh := state == rpio.High
	log.Debug().
		Bool("high", isHigh).
		Interface("raw", state).
		Msg("Current pin state")
	return isHigh
}
