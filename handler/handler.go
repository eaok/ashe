package handler

import (
	"fmt"
	"time"

	"github.com/eaok/khlashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

func AddReaction(ctx *khl.ReactionAddContext) {
	fmt.Println("go Reaction add!")

	// 根据userID获取username
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
	if err != nil {
		log.Error().Err(err).Msg("AddReaction")
	}
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.IDChannelSelectRole {
		if err := AddRoles(ctx); err != nil {
			log.Error().Err(err).Msg("AddReaction")
			return
		}
	}

	fmt.Println(ctx.Extra.Emoji.Name, "AddReaction")

	// 这3个emoji的动作传到team gorouting
	switch ctx.Extra.Emoji.Name {
	case config.EmojiCheckMark:
		fallthrough
	case config.EmojiCrossMark:
		fallthrough
	case EmojiHexToDec(config.EmojiStopSign):
		team.ReactionAdd <- ctx
	}
}

func DeleteReaction(ctx *khl.ReactionDeleteContext) {
	fmt.Println("go Reaction delete!")

	// 根据userID获取username
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
	if err != nil {
		log.Error().Err(err).Msg("DeleteReaction")
	}
	// bot event ignore
	if uv.Bot {
		return
	}

	// 角色选择频道
	if ctx.Extra.ChannelID == config.IDChannelSelectRole {
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
	// if strings.Contains(ctx.Common.Content, "ping") {
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

func Team(ctx *khl.TextMessageContext) {
	if RemovePrefix(ctx.Common.Content) == "team" {
		if team.running {
			return
		}

		// go startChannelTeam(config.IDChannelRS11, ctx)
		// go startChannelTeam(config.IDChannelRS10, ctx)
		// go startChannelTeam(config.IDChannelRS9, ctx)
		// go startChannelTeam(config.IDChannelRS8, ctx)
		go func() {
			for {
				chanDone := make(chan bool)
				go startChannelTeam(config.IDChannelRS7, ctx, chanDone)

				if <-chanDone {
					fmt.Printf("team has done!")
				}
			}
		}()
	}
}

func Help(ctx *khl.TextMessageContext) {
	// Ignore all messages created by the bot itself
	if ctx.Extra.Author.Bot {
		return
	}

	if RemovePrefix(ctx.Common.Content) == "help" {
		text := "```\n"
		text += fmt.Sprintf("%-5s\t:\t%s\n", "ping", "responds with pong!")
		text += fmt.Sprintf("%-5s\t:\t%s\n", "help", "prints this help menu!")
		text += "```"

		resp, _ := ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				Type:     khl.MessageTypeKMarkdown,
				TargetID: ctx.Common.TargetID,
				Content:  text,
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
