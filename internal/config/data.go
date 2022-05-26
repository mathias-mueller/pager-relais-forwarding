package config

type Config struct {
	TelegramConfig *TelegramConfig
}

type TelegramConfig struct {
	ChatID      int64
	APIToken    string
	MessageFile string
}

type HttpConfig struct {
	Endpoints []*HttpConfigItem
}

type HttpConfigItem struct {
	Name   string
	URL    string
	Method string
}
