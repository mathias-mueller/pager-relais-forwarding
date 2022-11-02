package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

func Load() (*Config, error) {
	cfg, loadErr := ini.InsensitiveLoad("config.ini")
	if loadErr != nil {
		return nil, fmt.Errorf("error loading config file: %+v", loadErr)
	}

	telegramSection := cfg.Section("telegram")

	var telegramConfig = &TelegramConfig{}
	err := telegramSection.StrictMapTo(telegramConfig)
	if err != nil {
		return nil, fmt.Errorf("error mapping telegram config section: %+v", err)
	}
	if telegramConfig.APIToken == "" {
		return nil, fmt.Errorf("key '%s' not set", "telegram.APIToken")
	}
	if telegramConfig.ChatID == 0 {
		return nil, fmt.Errorf("key '%s' not set", "telegram.ChatID")
	}
	if telegramConfig.MessageFile == "" {
		return nil, fmt.Errorf("key '%s' not set", "telegram.MessageFile")
	}

	gpioSection := cfg.Section("gpio")
	var gpioConfig = &GpioConfig{}
	err = gpioSection.StrictMapTo(gpioConfig)
	if err != nil {
		return nil, fmt.Errorf("error mapping gpio config section: %+v", err)
	}

	return &Config{
		TelegramConfig: telegramConfig,
		GpioConfig:     gpioConfig,
	}, nil
}
