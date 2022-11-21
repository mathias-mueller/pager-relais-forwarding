package config

type Config struct {
	GeneralConfig  *GeneralConfig
	TelegramConfig *TelegramConfig
	GpioConfig     *GpioConfig
}

type GeneralConfig struct {
	MetricsPort int64
	LogLevel    string
}

type TelegramConfig struct {
	ChatID      int64
	APIToken    string
	MessageFile string
}

type GpioConfig struct {
	Pin      int
	Interval int
}

type HTTPConfig struct {
	Endpoints []*HTTPConfigItem
}

type HTTPConfigItem struct {
	Name   string
	URL    string
	Method string
}
