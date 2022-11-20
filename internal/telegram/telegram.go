package telegram

import (
	"awesomeProject1/internal/config"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type API struct {
	bot       *tgbotapi.BotAPI
	conf      *config.TelegramConfig
	counter   prometheus.Counter
	histogram prometheus.Histogram
}

func Init(conf *config.TelegramConfig) *API {
	if content, err := os.ReadFile(conf.MessageFile); err != nil || len(content) == 0 {
		log.Fatal().
			Err(err).
			Str("file", conf.MessageFile).
			Bytes("content", content).
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
	return &API{
		bot:  bot,
		conf: conf,
		counter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "pager_forwarding_telegram_messages_total",
			Help: "The total number of telegram messages sent",
		}),
		histogram: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "pager_forwarding_telegram_messages_duration_sum",
			Help: "The total time while sending telegram messages",
		}),
	}
}

func (api *API) SendMsgString(text string) {
	api.counter.Inc()
	start := time.Now()
	msg := tgbotapi.NewMessage(api.conf.ChatID, text)

	_, e := api.bot.Send(msg)
	if e != nil {
		log.Err(e).Msg("Cannot send message")
	}
	log.Info().Msg("Msg sent")
	end := time.Now()
	api.histogram.Observe(float64(end.Sub(start).Milliseconds()))
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
