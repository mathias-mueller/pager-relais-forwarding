package telegram

import (
	"awesomeProject1/internal/config"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const retryDelay = time.Second * 10
const retryLimit = 3

type API struct {
	bot  *tgbotapi.BotAPI
	conf *config.TelegramConfig
}

func Init(conf *config.TelegramConfig) *API {
	if content, err := os.ReadFile(conf.MessageFile); err != nil || len(content) == 0 {
		log.Fatal().
			Err(err).
			Str("file", conf.MessageFile).
			Bytes("content", content).
			Msg("Error reading message file or file is empty")
	}
	remainingTries := retryLimit
	var bot *tgbotapi.BotAPI
	var err error
	for remainingTries > 0 {
		bot, err = tgbotapi.NewBotAPI(conf.APIToken)
		if err == nil {
			break
		}
		log.Warn().Err(err).Msg("Cannot connect to TelegramAPI. Retrying after delay...")
		remainingTries--
		time.Sleep(retryDelay)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to Telegram API")
	}
	bot.Debug = true

	log.Info().
		Str("user", bot.Self.UserName).
		Msg("Connected to Telegram API")
	return &API{
		bot:  bot,
		conf: conf,
	}
}

func (api *API) SendMsgString(text string) {
	msg := tgbotapi.NewMessage(api.conf.ChatID, text)

	_, e := api.bot.Send(msg)
	if e != nil {
		log.Err(e).Msg("Cannot send message")
	}
	log.Info().Msg("Msg sent")
}

func (api *API) SendMsg() {
	log.Info().
		Str("file", api.conf.MessageFile).
		Msg("Sending Telegram message from file")
	content, err := os.ReadFile(api.conf.MessageFile)
	if err != nil {
		log.Err(err).
			Str("file", api.conf.MessageFile).
			Msg("Failed to send telegram message. Cannot read message file")
		return
	}
	log.Info().Msg("Sending message")
	api.SendMsgString(string(content))
}
