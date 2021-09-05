package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

func CardButton(ctx *khl.MessageButtonClickContext) {
	fmt.Printf("value=%s msgID=%s userID=%s\n", ctx.Extra.Value, ctx.Extra.MsgID, ctx.Extra.UserID)

	// team.ReactionAdd <- ctx
}

func AddReaction(ctx *khl.ReactionAddContext) {
	fmt.Println(ctx.Extra.Emoji.Name, "AddReaction", ctx.Extra.ChannelID)

	// 根据userID获取username
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		log.Error().Err(err).Msg("AddReaction")
	}
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.Data.IDChannelSelectRole {
		if err := AddRoles(ctx); err != nil {
			log.Error().Err(err).Msg("AddReaction")
			return
		}
	}
}

func DeleteReaction(ctx *khl.ReactionDeleteContext) {
	fmt.Println(ctx.Extra.Emoji.Name, "DeleteReaction")

	// 根据userID获取username
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		log.Error().Err(err).Msg("DeleteReaction")
	}
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.Data.IDChannelSelectRole {
		if err := DeleteRoles(ctx); err != nil {
			log.Error().Err(err).Msg("DeleteReaction")
			return
		}
	}
}

// auto delete messages
func AutoDelete(ctx *khl.TextMessageContext) {
	// non bot messages are automatically deleted
	if ok := BotTakeOverGroup(ctx.Extra.ChannelName); !ctx.Extra.Author.Bot && ok {
		go func() {
			time.Sleep(15 * time.Second)
			ctx.Session.MessageDelete(ctx.Common.MsgID)
		}()
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
		team.OrderIn <- ctx
	}
}

func OutTeam(ctx *khl.TextMessageContext) {
	if RemovePrefix(ctx.Common.Content) == "out" {
		team.OrderOut <- ctx
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
