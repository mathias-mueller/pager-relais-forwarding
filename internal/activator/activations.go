package activator

import (
	"awesomeProject1/internal/telegram"

	"github.com/rs/zerolog/log"
)

type Activation interface {
	activate()
}

type TelegramActivation struct {
	API *telegram.API
}

func (t TelegramActivation) activate() {
	log.Info().Msg("Sending telegram Message")
	t.API.SendMsg()
}
