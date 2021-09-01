package handler

import (
	"os"
	"testing"

	"github.com/eaok/ashe/config"
)

func TestMain(m *testing.M) {
	config.ReadConfig("../config/config.ini")
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestEmojiHexToDec(t *testing.T) {
	type args struct {
		emoji string
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		// TODO: Add test cases.
		{
			name: "9",
			args: args{
				config.Data.EmojiNine,
			},
			wantStr: "[#57;][#65039;][#8419;]",
		},
		{
			name: "10",
			args: args{
				config.Data.EmojiTen,
			},
			wantStr: "[#128287;]",
		},
		{
			name: "✅",
			args: args{
				config.Data.EmojiCheckMark,
			},
			wantStr: "[#9989;]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotStr := EmojiHexToDec(tt.args.emoji); gotStr != tt.wantStr {
				t.Errorf("EmojiHexToDec() = %v, want %v", gotStr, tt.wantStr)
			}
		})
	}
}
