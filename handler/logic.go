package handler

import (
	"strconv"
	"strings"

	"github.com/eaok/khlashe/config"
)

// add order prefix
func RemovePrefix(content string) string {
	if config.Data.Prefix != "" && strings.HasPrefix(content, config.Data.Prefix) {
		return strings.TrimPrefix(content, config.Data.Prefix)
	}

	return ""
}

// bot takes over the group
func BotTakeOverGroup(ChannelName string) bool {
	switch ChannelName {
	case config.NameChannelSelectRole:
		fallthrough
	case config.NameChannelRS11:
		fallthrough
	case config.NameChannelRS10:
		fallthrough
	case config.NameChannelRS9:
		fallthrough
	case config.NameChannelRS8:
		fallthrough
	case config.NameChannelRS7:
		fallthrough
	case config.NameChannelRS6:
		fallthrough
	case config.NameChannelRS5:
		fallthrough
	case config.NameChannelRS4:
		return true
	default:
		return false
	}
}

// emoji编码转10进制格式
func EmojiHexToDec(emoji string) (str string) {
	for i := 0; i < len([]rune(emoji)); i++ {
		str += "[#"
		str += strconv.Itoa(int([]rune(emoji)[i]))
		str += ";]"
	}

	return
}
