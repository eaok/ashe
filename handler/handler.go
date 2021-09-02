package handler

import (
	"fmt"
	"sync"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

func CardButton(ctx *khl.MessageButtonClickContext) {
	fmt.Printf("value=%s msgID=%s userID=%s\n", ctx.Extra.Value, ctx.Extra.MsgID, ctx.Extra.UserID)

	team.ReactionAdd <- ctx
}

func AddReaction(ctx *khl.ReactionAddContext) {
	fmt.Println(ctx.Extra.Emoji.Name, "AddReaction", ctx.Extra.ChannelID)

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
	uv, err := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
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

func teamGoroutin(session *khl.Session, channelID string, reset chan bool) {
	wait := sync.WaitGroup{}

	for {
		wait.Add(1)
		go startChannelTeam(session, channelID, &wait)
		wait.Wait()
		fmt.Printf("%s team has done!\n", config.RSEmoji[config.ChanRole[channelID]])
	}
}

func Team(ctx *khl.TextMessageContext) {
	if RemovePrefix(ctx.Common.Content) == "team" {
		reset := make(chan bool)

		fmt.Printf("team.running=%v\n", team.running)
		if team.running {
			// team.Close <- true
			team.running = false
			return
		}

		go teamGoroutin(ctx.Session, config.Data.IDChannelRS11, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS10, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS9, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS8, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS7, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS6, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS5, reset)
		go teamGoroutin(ctx.Session, config.Data.IDChannelRS4, reset)

		// addreaction发送到指定goroutine
		go func() {
			for {
				fmt.Printf("team.MapGoroutine %v\n", team.MapGoroutine)
				select {
				case button := <-team.ReactionAdd:
					fmt.Printf("team.MapGoroutine[%v]%s\n", button.Extra.TargetID, button.Extra.Value)
					team.MapGoroutine[button.Extra.TargetID] <- button
				case <-team.Close:
					close(reset)
					return
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
