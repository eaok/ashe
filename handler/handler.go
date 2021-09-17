package handler

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

// auto delete messages
func AutoDelete(ctx *khl.TextMessageContext) {
	ctx.Session.Logger.Warn().Str("content", ctx.Common.Content).Str("userName", ctx.Extra.Author.Username).Str("ChannelName", ctx.Extra.ChannelName).Msg("AutoDelete")

	if config.Data.RunMod == "debug" {
		// non bot messages are automatically deleted
		if !ctx.Extra.Author.Bot && BotTakeOverGroup(ctx.Extra.ChannelName) {
			go func() {
				time.Sleep(15 * time.Second)
				ctx.Session.MessageDelete(ctx.Common.MsgID)
			}()
		}
	}
}

func Ping(ctx *khl.TextMessageContext) {
	if ctx.Common.Type != khl.MessageTypeText || ctx.Extra.Author.Bot {
		return
	}
	if RemovePrefix(ctx.Common.Content) == "ping" {
		resp, _ := ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "pong",
			},
		})

		// bot take over group auto delete message!
		if BotTakeOverGroup(ctx.Extra.ChannelName) {
			go func() {
				time.Sleep(15 * time.Second)
				ctx.Session.MessageDelete(resp.MsgID)
			}()
		}
	}
}

func InTeam(ctx *khl.TextMessageContext) {
	if strings.HasPrefix(RemovePrefix(ctx.Common.Content), "in") {
		if BotTakeOverRSGroup(ctx.Extra.ChannelName) {
			team.OrderIn <- ctx
		} else {
			SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 红星频道专属指令，请在红星频道输入！", ctx.Extra.Author.ID))
		}
	}
}

func OutTeam(ctx *khl.TextMessageContext) {
	if RemovePrefix(ctx.Common.Content) == "out" {
		if BotTakeOverRSGroup(ctx.Extra.ChannelName) {
			team.OrderOut <- ctx
		} else {
			SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 红星频道专属指令，请在红星频道输入！", ctx.Extra.Author.ID))
		}
	}
}

func Blue(ctx *khl.TextMessageContext) {
	if RemovePrefix(ctx.Common.Content) == "bb" {
		if BotTakeOverBSGroup(ctx.Extra.ChannelName) {
			BSTeam(ctx)
		} else {
			SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 蓝星频道专属指令，请在蓝星频道输入！", ctx.Extra.Author.ID))
		}
	}
}

func Rate(ctx *khl.TextMessageContext) {
	if ctx.Common.Type != khl.MessageTypeText || ctx.Extra.Author.Bot {
		return
	}
	if RemovePrefix(ctx.Common.Content) == "rate" {
		resp, _ := ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				TargetID: ctx.Common.TargetID,
				Content:  "https://img.kaiheila.cn/assets/2021-09/Q7l3QGSYvd0po07q.png",
			},
		})

		// bot take over group auto delete message!
		if BotTakeOverGroup(ctx.Extra.ChannelName) {
			go func() {
				time.Sleep(15 * time.Second)
				ctx.Session.MessageDelete(resp.MsgID)
			}()
		}
	}
}

func Order(ctx *khl.TextMessageContext) {
	if strings.HasPrefix(RemovePrefix(ctx.Common.Content), "order") {
		if ctx.Common.TargetID == config.Data.IDChannelTradePublish {
			// https://c.runoob.com/front-end/854/
			if matched, _ := regexp.MatchString(`^order\s([0-9]s[0-9]+\s)+[0-9]+s[tco]{1,3}$`, RemovePrefix(ctx.Common.Content)); matched {
				go TradeOrder(ctx)
			} else {
				SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 输入参数错误！", ctx.Extra.Author.ID))
			}

		} else {
			SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 发布订单频道专属指令，请在发布订单频道输入！", ctx.Extra.Author.ID))
		}
	}
}

func Help(ctx *khl.TextMessageContext) {
	// Ignore all messages created by the bot itself
	if ctx.Extra.Author.Bot {
		return
	}

	if RemovePrefix(ctx.Common.Content) == "help" {
		text := "> `"
		text += fmt.Sprintf("所有指令前缀为：%s", config.Data.Prefix)
		text += "`\n"
		text += "```\n"
		text += fmt.Sprintf("%-7s:    %s\n", "ping", "回复消息：pong!")
		text += fmt.Sprintf("%-7s:    %s\n", "in", "加入车队中，可跟数字[1-3]")
		text += fmt.Sprintf("%-7s:    %s\n", "out", "离开车队")
		text += fmt.Sprintf("%-7s:    %s\n", "bb", "创建一个10分钟的蓝星呼叫僚机队列")
		text += fmt.Sprintf("%-7s:    %s\n", "rate", "查看神器交易价格")
		text += fmt.Sprintf("%-7s:    %s\n", "order", "创建一个交易订单")
		text += fmt.Sprintf("%-7s:    %s\n", "help", "查看指令帮助菜单")
		text += "```"

		resp, _ := ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				Type:     khl.MessageTypeKMarkdown,
				TargetID: ctx.Common.TargetID,
				Content:  text,
			},
		})

		// bot take over group auto delete message!
		if BotTakeOverGroup(ctx.Extra.ChannelName) && config.Data.RunMod == "debug" {
			go func() {
				time.Sleep(15 * time.Second)
				ctx.Session.MessageDelete(resp.MsgID)
			}()
		}
	}
}
