package activator

import "github.com/rs/zerolog/log"

type Activator struct {
	currentValue bool
}

func New() *Activator {
	return &Activator{currentValue: false}
}

func (activator *Activator) EnableActivation(inputs <-chan bool, activations []Activation) {
	for input := range inputs {
		log.Trace().
			Bool("currentValue", activator.currentValue).
			Bool("input", input).
			Msg("Current values")
		if !input {
			activator.currentValue = false
			continue
		}
		if !activator.currentValue {
			log.Info().Msg("Activating...")
			for _, activation := range activations {
				go activation.activate()
			}
			activator.currentValue = true
		}
	}
}
