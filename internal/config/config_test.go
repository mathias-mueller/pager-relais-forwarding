package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_nofile(t *testing.T) {
	os.Remove("config.ini")

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad(t *testing.T) {
	type args struct {
		config string
	}
	const validGeneralSection = "[general]\n" +
		"MetricsPort = 8080\n" +
		"LogLevel = INFO\n"
	const validTelegramSection = "[telegram]\n" +
		"ChatID = 1234567890\n" +
		"APIToken = abcdefghijk\n" +
		"MessageFile = message.txt\n"
	const validGpioSection = "[gpio]\n" +
		"Pin = 10\n" +
		"Interval = 1000\n"
	var expectedGeneralSection = &GeneralConfig{
		MetricsPort: 8080,
		LogLevel:    "INFO",
	}
	var expectedTelgramConfig = &TelegramConfig{
		ChatID:      1234567890,
		APIToken:    "abcdefghijk",
		MessageFile: "message.txt",
	}
	expectedGpioConfig := &GpioConfig{
		Pin:      10,
		Interval: 1000,
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid",
			args: args{
				validGeneralSection +
					validTelegramSection +
					validGpioSection,
			},
			want: &Config{
				GeneralConfig:  expectedGeneralSection,
				TelegramConfig: expectedTelgramConfig,
				GpioConfig:     expectedGpioConfig,
			},
			wantErr: assert.NoError,
		},
		{
			name: "no APIToken",
			args: args{
				validGeneralSection +
					"[telegram]\n" +
					"ChatID = 1234567890\n" +
					"MessageFile = message.txt" +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no ChatID",
			args: args{
				validGeneralSection +
					"[telegram]\n" +
					"APIToken = abcdefghijk\n" +
					"MessageFile = message.txt" +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no API token",
			args: args{
				validGeneralSection +
					"[telegram]\n" +
					"ChatID = abc\n" +
					"MessageFile = message.txt" +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no messageFile",
			args: args{
				validGeneralSection +
					"[telegram]\n" +
					"ChatID = 1234567890\n" +
					"APIToken = abcdefghijk" +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no gpio pin",
			args: args{
				validGeneralSection +
					validTelegramSection +
					"[gpio]\n" +
					"Interval = 1000",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no gpio interval",
			args: args{
				validGeneralSection +
					validTelegramSection +
					"[gpio]\n" +
					"Pin = 10",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no metrics port",
			args: args{
				"[general]\n" +
					"LogLevel = INFO\n" +
					validTelegramSection +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no log level",
			args: args{
				"[general]\n" +
					"MetricsPort = 8080\n" +
					validTelegramSection +
					validGpioSection,
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.config)
			assert.NoError(t, os.WriteFile("config.ini", []byte(tt.args.config), os.ModePerm))
			defer func() { assert.NoError(t, os.Remove("config.ini")) }()
			got, err := Load()
			if !tt.wantErr(t, err, fmt.Sprintf("Load(%+v)", tt.args)) {
				return
			}
			fmt.Println(err)
			assert.Equalf(t, tt.want, got, "Load()")
		})
	}
}
