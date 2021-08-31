package handler

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/eaok/khlashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

type users struct {
	name   string
	nameID string
	time   int64
}

var (
	text1 = "**红星车队当前人数 [%d/4]**\n"
	text2 = "加入的成员：👇  |  🔴红星等级：@%s\n"
	text3 = "%s"
	text4 = "点击 ✅ 加入车队，点击 ❎ 离开车队，点击 🛑 直接发车！\n"
	// text1 = "**RED STAR QUEUE [%d/4]**\n"
	// text2 = "Members joined:👇  |  🔴RS level: @%s\n"
	// text3 = "%s"
	// text4 = "Use ✅ to join, ❎ to leave or 🛑 to start in %d.\n"
	// text5 = "Bored while waiting? Type !guide to refresh your knowledge!\n"
	Text = text1 + text2 + text3 + text4

// TeamBeginChan = make(chan int, 1)
)

type TeamData struct {
	sync.Mutex

	ReactionAdd chan *khl.ReactionAddContext
	TeamStart   chan bool
	Close       chan bool
	running     bool
}

var team = &TeamData{
	ReactionAdd: make(chan *khl.ReactionAddContext, 1),
	TeamStart:   make(chan bool, 1),
	Close:       make(chan bool),
	running:     false,
}

// startChannelTeam rs gorouting
func startChannelTeam(ChannelID string, ctx *khl.TextMessageContext, done chan bool) {
	fmt.Println("startChannelTeam")
	dict := map[string]users{}
	chanRS := make(chan bool, 1)

	team.running = true

	// 发送初始消息
	resp, err := sendFirstMessage(ctx.Session, ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("send first message failed! startChannelTeam")
		return
	}

	for {
		select {
		case reaction := <-team.ReactionAdd:
			// 如果reactionz有效就进入队伍
			if reaction.Extra.MsgID == resp.MsgID {
				fmt.Println(reaction.Extra.Emoji.Name, "startChannelTeam")
				switch reaction.Extra.Emoji.Name {
				case config.EmojiCheckMark:
					teamIn(dict, resp.MsgID, reaction, chanRS)
				case config.EmojiCrossMark:
					teamOut(dict, resp.MsgID, reaction)
				case EmojiHexToDec(config.EmojiStopSign):
					teamDone(dict, resp.MsgID, reaction, chanRS)
				default:
				}
			}
		case <-chanRS:
			done <- true
			return

		}
	}
}

// send init message
func sendFirstMessage(s *khl.Session, channelID string) (*khl.MessageResp, error) {
	fmt.Println("sendFirstMessage")
	c, err := s.ChannelView(channelID)
	if err != nil {
		return nil, err
	}

	resp, err := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: channelID,
			Content:  fmt.Sprintf(Text, 0, c.Name, ""),
		},
	})
	if err != nil {
		return nil, err
	}

	err = s.MessageAddReaction(resp.MsgID, config.EmojiCheckMark)

	return resp, err
}

// func teamMessgeEmojis(ctx *khl.ReactionAddContext) ([]khl.ReactionItem, error) {
// 	msg, err := ctx.Session.MessageList(ctx.Extra.ChannelID)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	// fmt.Println(msg[len(msg)-1].Reactions[0])

// 	return msg[len(msg)-1].Reactions, nil
// }

// func teamRemoveMessageAllEmojis(msgID string, ctx *khl.ReactionAddContext) {
// 	emojis, _ := teamMessgeEmojis(ctx)

// 	fmt.Println(emojis, "resetReaction 11111111")
// 	for index := range emojis {
// 		err := ctx.Session.MessageDeleteReaction(msgID, emojis[index].Emoji.ID, "")
// 		if err != nil {

// 			return
// 		}
// 	}
// }

func showEmojis(ctx *khl.ReactionAddContext) {
	msg, err := ctx.Session.MessageList(ctx.Extra.ChannelID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg[len(msg)-1].Reactions)

}

func resetReaction(dict map[string]users, msgID string, ctx *khl.ReactionAddContext, flag int) {
	fmt.Println("resetReaction", flag)
	if flag == 1 {
		// teamRemoveMessageAllEmojis(resp.MsgID, ctx)

		showEmojis(ctx)

		switch len(dict) {
		case 1:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCrossMark)
			showEmojis(ctx)

		case 2:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCrossMark)
		case 3:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCrossMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiStopSign)
		}
	} else if flag == 2 {
		switch len(dict) {
		case 0:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
		case 1:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCrossMark)
		case 2:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiStopSign, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.EmojiCrossMark)
		}
	} else if flag == 3 {
		switch len(dict) {
		// one person used to test
		case 1:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiStopSign, ctx.Extra.UserID)
		case 3:
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiStopSign, "")
			ctx.Session.MessageDeleteReaction(msgID, config.EmojiStopSign, ctx.Extra.UserID)
		}
	}
}

func sendTempMessage(s *khl.Session, channelID string, text string) {
	msg, _ := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			TargetID: channelID,
			Content:  text,
		},
	})
	go func() {
		time.Sleep(2 * time.Second)
		s.MessageDelete(msg.MsgID)
	}()
}

func teamGetSortNames(dict map[string]users) string {
	dictSort := make(map[int64]string)
	keys := []string{}
	namesList := []string{}

	for key, user := range dict {
		dictSort[user.time] = key
	}

	// 按时间排序取名字列表
	for timeKey := range dictSort {
		keys = append(keys, strconv.FormatInt(timeKey, 10))
	}
	sort.Strings(keys)

	for _, key := range keys {
		i, err := strconv.ParseInt(key, 10, 64)
		if err == nil {
			namesList = append(namesList, dict[dictSort[i]].name)
		}
	}

	names := ""
	for i := 0; i < len(dict); i++ {
		names += config.EmojiNum[i]
		names += " "
		names += namesList[i]
		names += "\n"
	}

	return names
}

func teamIn(dict map[string]users, msgID string, ctx *khl.ReactionAddContext, close chan bool) error {
	fmt.Println("teamIn")
	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.UserID]; ok {
		sendTempMessage(ctx.Session, ctx.Extra.ChannelID, fmt.Sprintf("@%s You're already in the team!", u.name))
		ctx.Session.MessageDeleteReaction(msgID, config.EmojiCheckMark, ctx.Extra.UserID)
		return nil
	} else {
		// 根据userID获取username
		uv, err := ctx.Session.UserView(ctx.Extra.UserID, ctx.Common.TargetID)
		if err != nil {
			return err
		}

		// join the team
		dict[ctx.Extra.UserID] = users{
			name:   uv.Username,
			nameID: ctx.Extra.UserID,
			time:   time.Now().Unix(),
		}

		// update message
		if len(dict) == 4 {
			teamDone(dict, msgID, ctx, close)
		} else {
			c, err := ctx.Session.ChannelView(ctx.Extra.ChannelID)
			if err != nil {
				return err
			}
			names := teamGetSortNames(dict)
			ctx.Session.MessageUpdate(&khl.MessageUpdate{
				MessageUpdateBase: khl.MessageUpdateBase{
					MsgID:   msgID,
					Content: fmt.Sprintf(Text, len(dict), c.Name, names),
				},
			})
			// reset reaction
			resetReaction(dict, msgID, ctx, 1)
		}
	}

	return nil
}

func teamOut(dict map[string]users, msgID string, ctx *khl.ReactionAddContext) error {
	fmt.Println("teamOut")
	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.UserID]; !ok {
		sendTempMessage(ctx.Session, ctx.Extra.ChannelID, fmt.Sprintf("@%s You're not in the team!", u.name))
		ctx.Session.MessageDeleteReaction(msgID, config.EmojiCrossMark, ctx.Extra.UserID)
		return nil
	} else {
		// leave the team
		delete(dict, ctx.Extra.UserID)

		// update message
		c, err := ctx.Session.ChannelView(ctx.Extra.ChannelID)
		if err != nil {
			return err
		}
		names := teamGetSortNames(dict)
		ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   msgID,
				Content: fmt.Sprintf(Text, len(dict), c.Name, names),
			},
		})

		// reset reaction
		resetReaction(dict, msgID, ctx, 2)
	}

	return nil
}

func teamDone(dict map[string]users, msgID string, ctx *khl.ReactionAddContext, close chan bool) error {
	fmt.Println("teamDone")
	c, err := ctx.Session.ChannelView(ctx.Extra.ChannelID)
	if err != nil {
		return err
	}
	names := teamGetSortNames(dict)
	ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   msgID,
			Content: fmt.Sprintf(Text, len(dict), c.Name, names),
		},
	})
	resetReaction(dict, msgID, ctx, 3)

	ctx.Session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: ctx.Extra.ChannelID,
			Content:  fmt.Sprintf("车队已经出发啦。。。 %s\n---\n", "hello"),
		},
	})

	close <- true

	return nil
}
