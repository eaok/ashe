package handler

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
	"github.com/phuslu/log"
)

type users struct {
	name   string
	nameID string
	time   time.Time
}

type TeamData struct {
	sync.Mutex

	ReactionAdd  chan *khl.ReactionAddContext
	MapGoroutine map[string]chan *khl.ReactionAddContext
	TeamStart    chan bool
	Close        chan bool
	running      bool
}

var (
	text1 = "**红星车队当前人数 [%d/4]**\n"
	text2 = "加入的成员：👇  |  🔴红星等级：(rol)%d(rol)\n"
	text3 = "%s"
	text4 = "点击 ✅ 加入车队，点击 ❎ 离开车队，点击 🛑 直接发车！\n"
	Text  = text1 + text2 + text3 + text4
)

var team = &TeamData{
	ReactionAdd:  make(chan *khl.ReactionAddContext, 1),
	MapGoroutine: make(map[string]chan *khl.ReactionAddContext),
	TeamStart:    make(chan bool, 1),
	Close:        make(chan bool),
	running:      false,
}

// startChannelTeam rs gorouting
func startChannelTeam(session *khl.Session, ChannelID string, done chan bool) {
	fmt.Printf("startChannelTeam ChannelID=%s\n", ChannelID)
	dict := map[string]users{}
	chanRS := make(chan bool, 1)

	team.running = true

	// 填充频道和chan通道的map，实现往指定goroutine发送数据
	reactionAdd := make(chan *khl.ReactionAddContext, 1)
	team.MapGoroutine[ChannelID] = reactionAdd

	// 发送初始消息
	resp, err := sendFirstMessage(session, ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("send first message failed! startChannelTeam")
		return
	}

	for {
		startTime := time.Now()
		select {
		case reaction := <-reactionAdd:
			// 如果reactionz有效就进入队伍
			if reaction.Extra.MsgID == resp.MsgID {
				fmt.Println(reaction.Extra.Emoji.Name, "startChannelTeam")
				switch reaction.Extra.Emoji.Name {
				case config.Data.EmojiCheckMark:
					teamIn(dict, resp.MsgID, reaction, chanRS)
				case config.Data.EmojiCrossMark:
					teamOut(dict, resp.MsgID, reaction)
				case EmojiHexToDec(config.Data.EmojiStopSign):
					teamDone(dict, resp.MsgID, reaction, chanRS)
				default:
				}
				fmt.Printf("dict %v", dict)
			}
		case <-time.After(time.Until(startTime.Add(time.Minute))):
			// 每分钟刷新一下组队信息，主要是排队时间
			if len(dict) > 0 {
				names := teamGetSortNames(dict)

				session.MessageUpdate(&khl.MessageUpdate{
					MessageUpdateBase: khl.MessageUpdateBase{
						MsgID:   resp.MsgID,
						Content: fmt.Sprintf(Text, len(dict), config.ChanRole[ChannelID], names),
					},
				})
			}

		case <-chanRS:
			done <- true
			return
		}
	}
}

// send init message
func sendFirstMessage(s *khl.Session, channelID string) (*khl.MessageResp, error) {
	resp, err := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: channelID,
			Content:  fmt.Sprintf(Text, 0, config.ChanRole[channelID], ""),
		},
	})
	if err != nil {
		return nil, err
	}

	err = s.MessageAddReaction(resp.MsgID, config.Data.EmojiCheckMark)

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
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCrossMark)
			showEmojis(ctx)

		case 2:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCrossMark)
		case 3:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCrossMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiStopSign)
		}
	} else if flag == 2 {
		switch len(dict) {
		case 0:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
		case 1:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCrossMark)
		case 2:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiStopSign, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, ctx.Extra.UserID)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCheckMark)
			ctx.Session.MessageAddReaction(msgID, config.Data.EmojiCrossMark)
		}
	} else if flag == 3 {
		switch len(dict) {
		// one person used to test
		case 1:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiStopSign, ctx.Extra.UserID)
		case 3:
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiStopSign, "")
			ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiStopSign, ctx.Extra.UserID)
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
	dictTime := make(map[int64]users)
	keys := []string{}
	namesList := []users{}

	// 转换dict为time索引
	for key, user := range dict {
		dictTime[user.time.Unix()] = dict[key]
	}

	// 按时间排序取名字列表
	for timeKey := range dictTime {
		keys = append(keys, strconv.FormatInt(timeKey, 10))
	}
	sort.Strings(keys)

	for _, key := range keys {
		// string to int64
		timeKeyInt, err := strconv.ParseInt(key, 10, 64)
		if err == nil {
			namesList = append(namesList, users{
				name:   dictTime[timeKeyInt].name,
				nameID: dictTime[timeKeyInt].nameID,
				time:   dictTime[timeKeyInt].time,
			})
		}
	}

	names := ""
	for i := 0; i < len(dict); i++ {
		timeSub := time.Since(namesList[i].time)
		value := fmt.Sprintf("%v", timeSub.Round(time.Second))

		names += config.EmojiNum[i]
		names += " "
		names += namesList[i].name
		names += " "
		names += value
		names += "\n"
	}

	return names
}

func teamIn(dict map[string]users, msgID string, ctx *khl.ReactionAddContext, close chan bool) error {
	fmt.Println("teamIn")
	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.UserID]; ok {
		sendTempMessage(ctx.Session, ctx.Extra.ChannelID, fmt.Sprintf("@%s You're already in the team!", u.name))
		ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCheckMark, ctx.Extra.UserID)
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
			time:   time.Now(),
		}

		// update message
		if len(dict) == 4 {
			teamDone(dict, msgID, ctx, close)
		} else {
			names := teamGetSortNames(dict)
			ctx.Session.MessageUpdate(&khl.MessageUpdate{
				MessageUpdateBase: khl.MessageUpdateBase{
					MsgID:   msgID,
					Content: fmt.Sprintf(Text, len(dict), config.ChanRole[ctx.Extra.ChannelID], names),
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
		ctx.Session.MessageDeleteReaction(msgID, config.Data.EmojiCrossMark, ctx.Extra.UserID)
		return nil
	} else {
		// leave the team
		delete(dict, ctx.Extra.UserID)

		names := teamGetSortNames(dict)
		ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   msgID,
				Content: fmt.Sprintf(Text, len(dict), config.ChanRole[ctx.Extra.ChannelID], names),
			},
		})

		// reset reaction
		resetReaction(dict, msgID, ctx, 2)
	}

	return nil
}

func teamDone(dict map[string]users, msgID string, ctx *khl.ReactionAddContext, close chan bool) error {
	fmt.Println("teamDone")
	names := teamGetSortNames(dict)
	ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   msgID,
			Content: fmt.Sprintf(Text, len(dict), config.ChanRole[ctx.Extra.ChannelID], names),
		},
	})
	resetReaction(dict, msgID, ctx, 3)

	ment := ""
	for key := range dict {
		ment += "(met)"
		ment += dict[key].nameID
		ment += "(met)"
	}
	ctx.Session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: ctx.Extra.ChannelID,
			Content:  fmt.Sprintf("%s车队已经出发啦。。。 \n---\n", ment),
		},
	})

	close <- true

	return nil
}
