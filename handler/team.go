package handler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

type users struct {
	name   string
	nameID string
	time   time.Time
	count  int
}

type TeamData struct {
	sync.Mutex
	OrderIn         chan *khl.TextMessageContext
	OrderOut        chan *khl.TextMessageContext
	MapInGoroutine  map[string]chan *khl.TextMessageContext
	MapOutGoroutine map[string]chan *khl.TextMessageContext
}

var Text = `[{
	"type": "card",
	"size": "lg",
	"theme": "primary",
	"modules": [{
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
		}
	]
}]`

var team = &TeamData{
	OrderIn:         make(chan *khl.TextMessageContext, 1),
	OrderOut:        make(chan *khl.TextMessageContext, 1),
	MapInGoroutine:  make(map[string]chan *khl.TextMessageContext),
	MapOutGoroutine: make(map[string]chan *khl.TextMessageContext),
}

func TeamStart(s *khl.Session) {
	go TeamGoroutin(s, config.Data.IDChannelRS11)
	go TeamGoroutin(s, config.Data.IDChannelRS10)
	go TeamGoroutin(s, config.Data.IDChannelRS9)
	go TeamGoroutin(s, config.Data.IDChannelRS8)
	go TeamGoroutin(s, config.Data.IDChannelRS7)
	go TeamGoroutin(s, config.Data.IDChannelRS6)
	go TeamGoroutin(s, config.Data.IDChannelRS5)
	go TeamGoroutin(s, config.Data.IDChannelRS4)

	// in out 发送到指定goroutine
	go func() {
		for {
			// fmt.Printf("team.MapInGoroutine %v\n", team.MapInGoroutine)
			// fmt.Printf("team.MapOutGoroutine %v\n", team.MapOutGoroutine)
			select {
			case in := <-team.OrderIn:
				team.MapInGoroutine[in.Common.TargetID] <- in
			case out := <-team.OrderOut:
				team.MapOutGoroutine[out.Common.TargetID] <- out
			}
		}
	}()
}

func TeamGoroutin(session *khl.Session, channelID string) {
	wait := sync.WaitGroup{}
	for {
		wait.Add(1)
		go TeamStartChannel(session, channelID, &wait)
		wait.Wait()
		fmt.Printf("%s team has done!\n", config.RSEmoji[config.ChanRole[channelID]])
	}
}

// startChannelTeam rs gorouting
func TeamStartChannel(session *khl.Session, ChannelID string, wait *sync.WaitGroup) {
	fmt.Printf("startChannelTeam ChannelID=%s\n", ChannelID)
	dict := map[string]users{}
	chanRS := make(chan bool, 1)

	teamIn := make(chan *khl.TextMessageContext, 1)
	teamOut := make(chan *khl.TextMessageContext, 1)

	// 填充频道和chan通道的map，实现往指定goroutine发送数据
	// 并发访问map不安全，会出现fatal error: concurrent map writes
	team.Mutex.Lock()
	team.MapInGoroutine[ChannelID] = teamIn
	team.MapOutGoroutine[ChannelID] = teamOut
	team.Mutex.Unlock()

	// 发送初始消息
	// TeamSendTempMessage(session, ChannelID, "team start......")

	for {
		select {
		case in := <-teamIn:
			TeamIn(dict, in, chanRS)
			fmt.Printf("dict %v\n", dict)
		case out := <-teamOut:
			TeamOut(dict, out)
			fmt.Printf("dict %v\n", dict)
		case <-chanRS:
			wait.Done()
			return
		}
	}
}

func TeamSendTempMessage(s *khl.Session, channelID string, text string) {
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

func TeamGetSortNames(dict map[string]users) string {
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
				count:  dictTime[timeKeyInt].count,
			})
		}
	}

	names := ""
	EmojiIndex := 0
	for i := 0; i < len(dict); i++ {
		// timeSub := time.Since(namesList[i].time)
		// value := fmt.Sprintf("%v", timeSub.Round(time.Second))
		// value := fmt.Sprintf("%v", namesList[i].time.Round(time.Second))
		value := namesList[i].time.Format("15 : 04")
		for j := 0; j < namesList[i].count; j++ {
			names += fmt.Sprintf("%s %s %15s\\n", config.EmojiNum[EmojiIndex], namesList[i].name, value)
			EmojiIndex++
		}
	}
	fmt.Println("names", names, "TeamGetSortNames")

	return names
}

// 组队人数是否已满
func TeamMember(dict map[string]users) int {
	teamMember := 0
	for key := range dict {
		teamMember += dict[key].count
	}

	return teamMember
}

func TeamIn(dict map[string]users, ctx *khl.TextMessageContext, close chan bool) error {
	fmt.Println("teamIn", ctx.Common.Content)

	// 处理指令参数
	var count int
	if len(RemovePrefix(ctx.Common.Content)) > 2 {
		content := strings.Fields(RemovePrefix(ctx.Common.Content))
		count, _ = strconv.Atoi(content[1])
		if count < 1 || count > 4 {
			TeamSendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 输入参数错误！", ctx.Extra.Author.ID))
			return nil
		}
	} else if len(RemovePrefix(ctx.Common.Content)) == 2 {
		count = 1
	}

	// Check user whether it is in team
	if u, ok := dict[ctx.Extra.Author.ID]; ok {
		TeamSendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 你已经在队伍中！", u.nameID))
		return nil
	} else {
		// 判断count和已有人数是否超额
		if TeamMember(dict)+count > 4 {
			TeamSendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 人数已经超过4人！", u.nameID))
			return nil
		}

		// join the team
		dict[ctx.Extra.Author.ID] = users{
			name:   ctx.Extra.Author.Username,
			nameID: ctx.Extra.Author.ID,
			time:   time.Now(),
			count:  count,
		}

		// channelID获取channelName guildID
		ch, err := ctx.Session.ChannelView(ctx.Common.TargetID)
		if err != nil {
			return err
		}

		// send new message
		names := TeamGetSortNames(dict)
		ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				Type:     khl.MessageTypeCard,
				TargetID: ctx.Common.TargetID,
				Content:  fmt.Sprintf(Text, ch.Name, names),
			},
		})

		if TeamMember(dict) == 4 {
			teamDone(dict, ctx, close)
		}
	}

	return nil
}

func TeamOut(dict map[string]users, ctx *khl.TextMessageContext) error {
	fmt.Println("teamOut", ctx.Common.Content)

	// 处理指令参数 out 没有参数
	if len(RemovePrefix(ctx.Common.Content)) > 3 {
		TeamSendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 输入参数错误！", ctx.Extra.Author.ID))
		return nil
	}

	// Check user whether it is in team
	if _, ok := dict[ctx.Extra.Author.ID]; !ok {
		TeamSendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 你没有再队伍中！", ctx.Extra.Author.ID))
		return nil
	} else {
		// leave the team
		delete(dict, ctx.Extra.Author.ID)

		// channelID获取channelName guildID
		ch, err := ctx.Session.ChannelView(ctx.Common.TargetID)
		if err != nil {
			return err
		}

		// send new message
		names := TeamGetSortNames(dict)
		ctx.Session.MessageCreate(&khl.MessageCreate{
			MessageCreateBase: khl.MessageCreateBase{
				Type:     khl.MessageTypeCard,
				TargetID: ctx.Common.TargetID,
				Content:  fmt.Sprintf(Text, ch.Name, names),
			},
		})
	}

	return nil
}

func teamDone(dict map[string]users, ctx *khl.TextMessageContext, close chan bool) error {
	fmt.Println("teamDone")

	ment := ""
	for key := range dict {
		ment += "(met)"
		ment += dict[key].nameID
		ment += "(met)"
	}
	ctx.Session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: ctx.Common.TargetID,
			Content:  fmt.Sprintf("%s车队已经出发啦。。。 \n---\n", ment),
		},
	})

	close <- true

	return nil
}
