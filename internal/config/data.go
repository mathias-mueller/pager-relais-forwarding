package config

type Config struct {
	TelegramConfig *TelegramConfig
	GpioConfig     *GpioConfig
}

type TelegramConfig struct {
	ChatID      int64
	APIToken    string
	MessageFile string
}

type GpioConfig struct {
	Pin int
}

type HttpConfig struct {
	Endpoints []*HttpConfigItem
}

type HttpConfigItem struct {
	Name   string
	URL    string
	Method string
}
