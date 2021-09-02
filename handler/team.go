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

	ReactionAdd  chan *khl.MessageButtonClickContext
	MapGoroutine map[string]chan *khl.MessageButtonClickContext
	TeamStart    chan bool
	Close        chan bool
	running      bool
}

var Text = `[
	{
	  "type": "card",
	  "size": "lg",
	  "theme": "warning",
	  "modules": [
		{
		  "type": "header",
		  "text": {
			"type": "plain-text",
			"content": "红星车队%s"
		  }
		},
		{
		  "type": "divider"
		},
		{
		  "type": "section",
		  "text": {
			"type": "kmarkdown",
			"content": "%s"
		  }
		},
		{
		  "type": "action-group",
		  "elements": [
			{
			  "type": "button",
			  "theme": "primary",
			  "value": "ok",
			  "click": "return-val",
			  "text": {
				"type": "plain-text",
				"content": "加入"
			  }
			},
			{
			  "type": "button",
			  "theme": "danger",
			  "value": "cancel",
			  "click": "return-val",
			  "text": {
				"type": "plain-text",
				"content": "离开"
			  }
			},
			{
			  "type": "button",
			  "theme": "primary",
			  "value": "begin",
			  "click": "return-val",
			  "text": {
				"type": "plain-text",
				"content": "开始"
			  }
			}
		  ]
		}
	  ]
	}
  ]`

// var (
// 	text1 = "**红星车队当前人数 [%d/4]**\n"
// 	text2 = "加入的成员：👇  |  🔴红星等级：(rol)%d(rol)\n"
// 	text3 = "%s"
// 	text4 = "点击 ✅ 加入车队，点击 ❎ 离开车队，点击 🛑 直接发车！\n"
// 	Text  = text1 + text2 + text3 + text4
// )

var team = &TeamData{
	ReactionAdd:  make(chan *khl.MessageButtonClickContext, 1),
	MapGoroutine: make(map[string]chan *khl.MessageButtonClickContext),
	TeamStart:    make(chan bool, 1),
	Close:        make(chan bool),
	running:      false,
}

// startChannelTeam rs gorouting
func startChannelTeam(session *khl.Session, ChannelID string, wait *sync.WaitGroup) {
	fmt.Printf("startChannelTeam ChannelID=%s\n", ChannelID)
	dict := map[string]users{}
	chanRS := make(chan bool, 1)

	team.running = true

	buttonChan := make(chan *khl.MessageButtonClickContext, 1)

	// 填充频道和chan通道的map，实现往指定goroutine发送数据
	// 并发访问map不安全，会出现fatal error: concurrent map writes
	session.RWMutex.Lock()
	team.MapGoroutine[ChannelID] = buttonChan
	session.RWMutex.Unlock()

	// 发送初始消息
	resp, err := sendFirstMessage(session, ChannelID)
	if err != nil {
		log.Error().Err(err).Msg("send first message failed! startChannelTeam")
		return
	}

	// channelID获取channelName guildID
	ch, _ := session.ChannelView(ChannelID)

	for {
		startTime := time.Now()
		select {
		case button := <-buttonChan:
			if button.Extra.MsgID == resp.MsgID {
				switch button.Extra.Value {
				case "ok":
					teamIn(dict, resp.MsgID, button, chanRS)
				case "cancel":
					teamOut(dict, resp.MsgID, button)
				case "begin":
					teamDone(dict, resp.MsgID, button, chanRS)
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
						Content: fmt.Sprintf(Text, ch.Name, names),
					},
				})
			}

		case <-chanRS:
			wait.Done()
			return
		}
	}
}

// send init message
func sendFirstMessage(s *khl.Session, channelID string) (*khl.MessageResp, error) {
	// channelID获取channelName
	ch, err := s.ChannelView(channelID)
	if err != nil {
		return nil, err
	}

	resp, err := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeCard,
			TargetID: channelID,
			Content:  fmt.Sprintf(Text, ch.Name, ""),
		},
	})
	if err != nil {
		return nil, err
	}

	return resp, err
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

		// if i != 0 {
		// 	names += "\n"
		// }
		names += fmt.Sprintf("%s %s %10s\\n", config.EmojiNum[i], namesList[i].name, value)
	}

	return names
}

func teamIn(dict map[string]users, msgID string, ctx *khl.MessageButtonClickContext, close chan bool) error {
	fmt.Println("teamIn")
	// channelID获取channelName guildID
	ch, err := ctx.Session.ChannelView(ctx.Extra.TargetID)
	if err != nil {
		return err
	}
	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.UserID]; ok {
		sendTempMessage(ctx.Session, ctx.Extra.TargetID, fmt.Sprintf("@%s You're already in the team!", u.name))
		return nil
	} else {
		// 根据userID获取username
		uv, err := ctx.Session.UserView(ctx.Extra.UserID, ch.GuildID)
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
			fmt.Println(names)
			ctx.Session.MessageUpdate(&khl.MessageUpdate{
				MessageUpdateBase: khl.MessageUpdateBase{
					MsgID:   msgID,
					Content: fmt.Sprintf(Text, ch.Name, names),
				},
			})
		}
	}

	return nil
}

func teamOut(dict map[string]users, msgID string, ctx *khl.MessageButtonClickContext) error {
	fmt.Println("teamOut")
	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.UserID]; !ok {
		sendTempMessage(ctx.Session, ctx.Extra.TargetID, fmt.Sprintf("@%s You're not in the team!", u.name))
		return nil
	} else {
		// leave the team
		delete(dict, ctx.Extra.UserID)

		// channelID获取channelName guildID
		ch, err := ctx.Session.ChannelView(ctx.Extra.TargetID)
		if err != nil {
			return err
		}

		names := teamGetSortNames(dict)
		ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   msgID,
				Content: fmt.Sprintf(Text, ch.Name, names),
			},
		})
	}

	return nil
}

func teamDone(dict map[string]users, msgID string, ctx *khl.MessageButtonClickContext, close chan bool) error {
	fmt.Println("teamDone")

	// channelID获取channelName guildID
	ch, err := ctx.Session.ChannelView(ctx.Extra.TargetID)
	if err != nil {
		return err
	}

	names := teamGetSortNames(dict)
	ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   msgID,
			Content: fmt.Sprintf(Text, ch.Name, names),
		},
	})

	ment := ""
	for key := range dict {
		ment += "(met)"
		ment += dict[key].nameID
		ment += "(met)"
	}
	ctx.Session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: ch.ID,
			Content:  fmt.Sprintf("%s车队已经出发啦。。。 \n---\n", ment),
		},
	})

	close <- true

	return nil
}
