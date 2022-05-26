package gpio

import (
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)
import "time"

func Start() <-chan time.Time {
	err := rpio.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init GPIO library")
	}
	IsPinHigh()

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
					if IsPinHigh() {
						output <- time.Now()
					}
				}()
			}
		}
	}()
	return output
}

func IsPinHigh() bool {
	pinNumber := 10
	log.Trace().
		Int("number", pinNumber).
		Msg("Reading pin")
	pin := rpio.Pin(pinNumber)

	state := pin.Read()
	isHigh := state == rpio.High
	log.Debug().
		Bool("high", isHigh).
		Interface("raw", state).
		Msg("Current pin state")
	return isHigh
}
