package activator

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type Activator struct {
	currentValue bool
	counter      prometheus.Counter
}

func New() *Activator {
	return &Activator{
		currentValue: false,
		counter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "pager_forwarding_activations_total",
			Help: "The total number of activations",
		}),
	}
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
			activator.counter.Inc()
			for _, activation := range activations {
				go activation.activate()
			}
			activator.currentValue = true
		}
	}
}
