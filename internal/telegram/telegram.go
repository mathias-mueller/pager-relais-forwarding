package telegram

import (
	"awesomeProject1/internal/config"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

type Api struct {
	bot  *tgbotapi.BotAPI
	conf *config.TelegramConfig
}

func Init(conf *config.TelegramConfig) *Api {

	if content, err := os.ReadFile(conf.MessageFile); err == nil || len(content) == 0 {
		log.Fatal().
			Err(err).
			Str("file", conf.MessageFile).
			Msg("Error reading message file or file is empty")
	}

	bot, err := tgbotapi.NewBotAPI(conf.APIToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to Telegram API")
	}
	bot.Debug = true

	log.Info().
		Str("user", bot.Self.UserName).
		Msg("Connected to Telegram API")
	return &Api{
		bot:  bot,
		conf: conf,
	}
}

func (api *Api) SendMsgString(text string) {
	msg := tgbotapi.NewMessage(api.conf.ChatID, text)

	_, e := api.bot.Send(msg)
	if e != nil {
		log.Err(e).Msg("Cannot send message")
	}
	log.Info().Msg("Msg sent")
}

func (api *Api) SendMsg() {
	content, err := os.ReadFile(api.conf.MessageFile)
	if err != nil {
		log.Err(err).Msg("Failed to send telegram message. Cannot read message file")
		return
	}
	api.SendMsgString(string(content))
}
