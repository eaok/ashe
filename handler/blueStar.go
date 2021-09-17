package handler

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

var BlueText = `[
	{
	  "type": "card",
	  "size": "lg",
	  "theme": "warning",
	  "modules": [
		{
		  "type": "section",
		  "text": {
			"type": "kmarkdown",
			"content": "**主机组：**\n%s"
		  }
		},
		{
		  "type": "divider"
		},
		{
		  "type": "section",
		  "text": {
			"type": "kmarkdown",
			"content": "**僚机组：**%s"
		  }
		},
		{
		  "type": "divider"
		},
		{
		  "type": "countdown",
		  "mode": "day",
		  "endTime": %d
		},
		{
		  "type": "divider"
		},
		{
		  "type": "section",
		  "text": {
			"type": "plain-text",
			"content": "僚机点击:white_check_mark:加入或退出僚机组\n主机点击:white_check_mark:开始"
		  }
		}
	  ]
	}
  ]`

type BSusers struct {
	name    string
	nameID  string
	role    int64
	emoji   string
	time    time.Time
	timeout int64
}

type BSTeamData struct {
	sync.Mutex
	MapBSAddGoroutine    map[string]chan *khl.ReactionAddContext
	MapBSDeleteGoroutine map[string]chan *khl.ReactionDeleteContext
}

var BSteam = &BSTeamData{
	MapBSAddGoroutine:    make(map[string]chan *khl.ReactionAddContext, 1),
	MapBSDeleteGoroutine: make(map[string]chan *khl.ReactionDeleteContext, 1),
}

// 蓝星呼叫僚机总函数
func BSTeam(ctx *khl.TextMessageContext) {
	startTime := time.Now()
	ctx.Session.Logger.Warn().Str("content", ctx.Common.Content).Str("TargetID", ctx.Common.TargetID).Msg("BSTeam")

	if role := BSGetMaxRole(ctx.Session, ctx.Extra.Author.ID, ctx.Extra.GuildID); role != 0 {
		resp, _ := BSFirstMessage(ctx, role, startTime)

		go BSGoroutine(ctx, resp.MsgID, role, startTime)
	} else {
		SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 请先选择蓝星角色！", ctx.Extra.Author.ID))
		ctx.Session.Logger.Warn().Msgf("%s 请先选择蓝星角色！", ctx.Extra.Author.Username)
	}
}

func BSGoroutine(ctx *khl.TextMessageContext, msgID string, role int64, startTime time.Time) {
	chanBS := make(chan bool, 1)
	dict := map[string]BSusers{}

	BSteam.Mutex.Lock()
	BSteam.MapBSAddGoroutine[msgID] = make(chan *khl.ReactionAddContext)
	BSteam.MapBSDeleteGoroutine[msgID] = make(chan *khl.ReactionDeleteContext)
	BSteam.Mutex.Unlock()

	// 主机数据先加入map中
	dict["1"] = BSusers{
		name:    ctx.Extra.Author.Username,
		nameID:  ctx.Extra.Author.ID,
		role:    role,
		emoji:   config.RSEmoji[role],
		time:    time.Now(),
		timeout: startTime.Add(10*time.Minute).UnixNano() / 1e6,
	}
	// fmt.Println(BSteam.MapBSAddGoroutine)
	// fmt.Println(BSteam.MapBSDeleteGoroutine)

	for {
		select {
		case add := <-BSteam.MapBSAddGoroutine[msgID]:
			if dict["1"].nameID == add.Extra.UserID {
				BSTeamDone(dict, add, chanBS, msgID)
			} else {
				BSTeamIn(dict, add, msgID, chanBS, role)
			}
		case delete := <-BSteam.MapBSDeleteGoroutine[msgID]:
			BSTeamOut(dict, delete, msgID)
		case <-time.After(time.Until(startTime.Add(10 * time.Minute))):
			return
		case <-chanBS:
			return
		}
		// fmt.Println(BSteam.MapBSAddGoroutine)
		// fmt.Println(BSteam.MapBSDeleteGoroutine)
	}
}

// 加入僚机组
func BSTeamIn(dict map[string]BSusers, ctx *khl.ReactionAddContext, msgID string, chanBS chan bool, masterRole int64) {
	if role := BSGetMaxRole(ctx.Session, ctx.Common.AuthorID, ctx.Common.TargetID); role != 0 {
		if BSGetIsMatch(role, masterRole) {
			// 获取user信息
			user, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
			if err != nil {
				ctx.Session.Logger.Error().Err("", err).Msg("BSTeamIn UserView")
				return
			}

			// 满足条件，加入僚机组
			dict[ctx.Extra.UserID] = BSusers{
				name:   user.Username,
				nameID: user.ID,
				role:   role,
				emoji:  config.RSEmoji[role],
				time:   time.Now(),
			}

			// 更新消息
			nameMaster, nameElse, timeout := BSTeamGetNames(dict)
			ctx.Session.MessageUpdate(&khl.MessageUpdate{
				MessageUpdateBase: khl.MessageUpdateBase{
					MsgID:   msgID,
					Content: fmt.Sprintf(BlueText, nameMaster, nameElse, timeout),
				},
			})

			if len(dict) == 4 {
				BSTeamDone(dict, ctx, chanBS, msgID)
			}
		} else {
			SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 你们战舰配置相差太大，匹配不到一起！", ctx.Extra.UserID))
			ctx.Session.Logger.Warn().Msgf("%s 你们战舰配置相差太大，匹配不到一起！", ctx.Extra.UserID)
		}
	} else {
		SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 请先选择蓝星角色！", ctx.Extra.UserID))
		ctx.Session.Logger.Warn().Msgf("%s 请先选择蓝星角色！", ctx.Extra.UserID)
	}
}

// 退出僚机组
func BSTeamOut(dict map[string]BSusers, ctx *khl.ReactionDeleteContext, msgID string) {
	// 检查用户是否在僚机组中
	if _, ok := dict[ctx.Extra.UserID]; !ok {
		SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 你没有在僚机队伍中！", ctx.Extra.UserID))
		ctx.Session.Logger.Warn().Msgf("%s 你没有在僚机队伍中！", ctx.Extra.UserID)
		return
	} else {
		// 离开僚机组
		delete(dict, ctx.Extra.UserID)

		// 更新消息
		nameMaster, nameElse, timeout := BSTeamGetNames(dict)
		ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   msgID,
				Content: fmt.Sprintf(BlueText, nameMaster, nameElse, timeout),
			},
		})
	}
}

// 组队完成
func BSTeamDone(dict map[string]BSusers, ctx *khl.ReactionAddContext, chanBS chan bool, msgID string) {
	// 更新消息
	nameMaster, nameElse, _ := BSTeamGetNames(dict)
	ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   msgID,
			Content: fmt.Sprintf(BlueText, nameMaster, nameElse, time.Now().Add(5*time.Second).UnixNano()/1e6),
		},
	})

	// 创建一个临时语音频道
	// names := strings.Fields(nameMaster)
	// c, err := ctx.Session.ChannelCreate(&khl.ChannelCreate{
	// 	GuildID:     ctx.Common.TargetID,
	// 	ParentID:    config.Data.IDChannelGroupBS,
	// 	Name:        names[1],
	// 	Type:        khl.ChannelTypeVoice,
	// 	LimitAmount: 4,
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// go func() {
	// 	time.Sleep(1 * time.Minute)
	// 	ctx.Session.ChannelDelete(c.ID)
	// }()

	// 发送已经完成的消息
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
			Content:  fmt.Sprintf("%s蓝星呼叫僚机已经完成。。。 \n---\n", ment),
			// Content:  fmt.Sprintf("%s蓝星呼叫僚机已经完成。。。 \n已经创建语音频道(chn)%s(chn)，点击可以加入 \n---\n", ment, c.ID),
		},
	})

	BSteam.Lock()
	delete(BSteam.MapBSAddGoroutine, msgID)
	delete(BSteam.MapBSDeleteGoroutine, msgID)
	BSteam.Unlock()

	chanBS <- true
}

// 从map中获取主机和僚机信息
func BSTeamGetNames(dict map[string]BSusers) (string, string, int64) {
	dictTime := make(map[int64]BSusers)
	keys := []string{}
	namesList := []BSusers{}
	nameMaster := ""
	nameElse := ""
	var timeout int64

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
			namesList = append(namesList, BSusers{
				name:   dictTime[timeKeyInt].name,
				nameID: dictTime[timeKeyInt].nameID,
				role:   dictTime[timeKeyInt].role,
				emoji:  dictTime[timeKeyInt].emoji,
				time:   dictTime[timeKeyInt].time,
			})
		}
	}

	for i := 0; i < len(dict); i++ {
		value := namesList[i].time.Format("15 : 04")
		if i == 0 {
			nameMaster += fmt.Sprintf("%s %s %15s", namesList[i].emoji, namesList[i].name, value)
			timeout = namesList[i].timeout
		} else {
			nameElse += fmt.Sprintf("\\n%s %s %15s", namesList[i].emoji, namesList[i].name, value)
		}
	}

	return nameMaster, nameElse, timeout
}

// 两个角色是否能够匹配
func BSGetIsMatch(role, masterRole int64) bool {
	roleIndex := 0
	masterRoleIndex := 0

	for i := 0; i < 8; i++ {
		if config.BSRoleNum[i] == role {
			roleIndex = i
		}
		if config.BSRoleNum[i] == masterRole {
			masterRoleIndex = i
		}
	}

	return math.Abs(float64(roleIndex-masterRoleIndex)) < 3
}

// 发送初始消息
func BSFirstMessage(ctx *khl.TextMessageContext, role int64, startTime time.Time) (resp *khl.MessageResp, err error) {
	name := fmt.Sprintf("%s %s %15s", config.RSEmoji[role], ctx.Extra.Author.Username, time.Now().Format("15 : 04"))
	timeout := startTime.Add(10*time.Minute).UnixNano() / 1e6
	resp, err = ctx.Session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeCard,
			TargetID: ctx.Common.TargetID,
			Content:  fmt.Sprintf(BlueText, name, "", timeout),
		},
	})

	ctx.Session.MessageAddReaction(resp.MsgID, config.Data.EmojiCheckMark)

	return
}

// 获取用户最大的蓝星角色
func BSGetMaxRole(s *khl.Session, userID string, guildID string) int64 {
	user, err := s.UserView(userID, khl.UserViewWithGuildID(guildID))
	if err != nil {
		s.Logger.Error().Err("", err).Msg("BSGetMaxRole UserView")
	}

	for i := 7; i >= 0; i-- {
		for _, role := range user.Roles {
			if int64(role) == config.BSRoleNum[i] {
				return role
			}
		}
	}

	return 0
}
