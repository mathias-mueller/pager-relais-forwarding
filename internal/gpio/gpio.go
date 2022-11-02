package gpio

import (
	"awesomeProject1/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)
import "time"

func Start(config *config.GpioConfig) <-chan time.Time {
	err := rpio.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init GPIO library")
	}

	output := make(chan time.Time, 0)

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan time.Time)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				go func() {
					if IsPinHigh(config) {
						output <- time.Now()
					}
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
