package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad_nofile(t *testing.T) {
	os.Remove("config.ini.default")

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad(t *testing.T) {
	type args struct {
		config string
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
				"[telegram]\n" +
					"ChatID = 1234567890\n" +
					"APIToken = abcdefghijk\n" +
					"MessageFile = message.txt",
			},
			want: &Config{TelegramConfig: &TelegramConfig{
				ChatID:      1234567890,
				APIToken:    "abcdefghijk",
				MessageFile: "message.txt",
			}},
			wantErr: assert.NoError,
		},
		{
			name: "invalid",
			args: args{
				"[telegram]\n" +
					"foo = bar\n",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no APIToken",
			args: args{
				"[telegram]\n" +
					"ChatID = 1234567890\n" +
					"MessageFile = message.txt",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no ChatID",
			args: args{
				"[telegram]\n" +
					"APIToken = abcdefghijk\n" +
					"MessageFile = message.txt",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no API token",
			args: args{
				"[telegram]\n" +
					"ChatID = abc\n" +
					"MessageFile = message.txt",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "no messageFile",
			args: args{
				"[telegram]\n" +
					"ChatID = 1234567890\n" +
					"APIToken = abcdefghijk",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, os.WriteFile("config.ini", []byte(tt.args.config), os.ModePerm))
			defer func() { assert.NoError(t, os.Remove("config.ini")) }()
			got, err := Load()
			if !tt.wantErr(t, err, fmt.Sprintf("Load()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Load()")

		})
	}
}
