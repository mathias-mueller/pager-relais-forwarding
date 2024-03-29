package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func Load() (*Config, error) {
	cfg, loadErr := ini.InsensitiveLoad("config.ini")
	if loadErr != nil {
		return nil, fmt.Errorf("error loading config file: %w", loadErr)
	}
	generalSection := cfg.Section("general")

	var generalConfig = &GeneralConfig{}
	err := generalSection.StrictMapTo(generalConfig)
	if err != nil {
		return nil, fmt.Errorf("error mapping telegram config section: %w", err)
	}
	if generalConfig.LogLevel == "" {
		return nil, fmt.Errorf("key '%s' not set", "general.LogLevel")
	}
	if generalConfig.MetricsPort == 0 {
		return nil, fmt.Errorf("key '%s' not set", "general.MetricsPort")
	}

	telegramSection := cfg.Section("telegram")

	var telegramConfig = &TelegramConfig{}
	err = telegramSection.StrictMapTo(telegramConfig)
	if err != nil {
		return nil, fmt.Errorf("error mapping telegram config section: %w", err)
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
		return nil, fmt.Errorf("error mapping gpio config section: %w", err)
	}
	if gpioConfig.Pin == 0 {
		return nil, fmt.Errorf("key '%s' not set", "gpio.Pin")
	}
	if gpioConfig.Interval <= 0 {
		return nil, fmt.Errorf("key '%s' not set", "gpio.Interval")
	}

	return &Config{
		GeneralConfig:  generalConfig,
		TelegramConfig: telegramConfig,
		GpioConfig:     gpioConfig,
	}, nil
}
