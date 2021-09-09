package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

// add order prefix
func RemovePrefix(content string) string {
	if config.Data.Prefix != "" && strings.HasPrefix(content, config.Data.Prefix) {
		return strings.TrimPrefix(content, config.Data.Prefix)
	} else if config.Data.Prefix != "" {
		return "nonono"
	}

	return content
}

// bot takes over the group
func BotTakeOverGroup(ChannelName string) bool {
	if BotTakeOverRSGroup(ChannelName) || BotTakeOverRoleSelectGroup(ChannelName) ||
		BotTakeOverBSGroup(ChannelName) {
		return true
	} else {
		return false
	}
}

func BotTakeOverRoleSelectGroup(ChannelName string) bool {
	switch ChannelName {
	case config.Data.NameChannelSelectRole:
		return true
	default:
		return false
	}
}

func BotTakeOverRSGroup(ChannelName string) bool {
	switch ChannelName {
	case config.Data.NameChannelRS11:
		fallthrough
	case config.Data.NameChannelRS10:
		fallthrough
	case config.Data.NameChannelRS9:
		fallthrough
	case config.Data.NameChannelRS8:
		fallthrough
	case config.Data.NameChannelRS7:
		fallthrough
	case config.Data.NameChannelRS6:
		fallthrough
	case config.Data.NameChannelRS5:
		fallthrough
	case config.Data.NameChannelRS4:
		return true
	default:
		return false
	}
}

func BotTakeOverBSGroup(ChannelName string) bool {
	switch ChannelName {
	case config.Data.NameChannelBS:
		return true
	default:
		return false
	}
}

// 发送临时消息
func SendTempMessage(s *khl.Session, channelID string, text string) {
	msg, _ := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: channelID,
			Content:  text,
		},
	})
	go func() {
		time.Sleep(30 * time.Second)
		s.MessageDelete(msg.MsgID)
	}()
}

// 发送私人消息
func SendDirectMessage(s *khl.Session, targetID string, text string) {
	s.DirectMessageCreate(&khl.DirectMessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: targetID,
			Content:  text,
		},
	})
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
