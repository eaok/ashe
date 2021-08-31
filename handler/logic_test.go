package handler

import (
	"testing"

	"github.com/eaok/khlashe/config"
)

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
				config.EmojiNine,
			},
			wantStr: "[#57;][#65039;][#8419;]",
		},
		{
			name: "10",
			args: args{
				config.EmojiTen,
			},
			wantStr: "[#128287;]",
		},
		{
			name: "✅",
			args: args{
				config.EmojiCheckMark,
			},
			wantStr: "[#128287;]",
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
