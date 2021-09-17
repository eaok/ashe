package handler

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/lonelyevil/khl"
)

var TextTradePublish = `
[
  {
    "type": "card",
    "theme": "secondary",
    "size": "lg",
    "modules": [
      {
        "type": "header",
        "text": {
          "type": "plain-text",
          "content": "-------------订单 #%d-------------"
        }
      },
      {
        "type": "section",
        "text": {
          "type": "paragraph",
          "cols": 2,
          "fields": [
            {
              "type": "kmarkdown",
              "content": "**订单发布者**\n%s"
            },
            {
              "type": "kmarkdown",
              "content": "**接单者**%s"
            }
          ]
        }
      },
      {
        "type": "divider"
      },
      {
        "type": "section",
        "text": {
          "type": "plain-text",
          "content": "货物：%s\n求购：%s"
        }
      },
      {
        "type": "divider"
      },
      {
        "type": "section",
        "text": {
          "type": "plain-text",
          "content": "当前折扣：%d%%"
        }
      },
	  {
		"type": "section",
		"text": {
		  "type": "plain-text",
		  "content": "发布时间：%s"
		}
	  }
    ]
  }
]`

var OrderNum = 1000

var TradePrice = [8][7]float64{
	{2.5, 2},
	{3.5, 2.5, 2},
	{4.5, 3.5, 2.5, 2},
	{5.5, 4.5, 3.5, 2.5, 2},
	{0, 0, 5.5, 4.5, 3.5, 3},
	{0, 0, 12, 8, 6, 4.5, 4.5},

	{0, 0, 7, 6, 5, 4},       // 10t
	{0, 0, 24, 14, 11, 8, 8}, // 11t
}

type TradeContent struct {
	Cargo    []TradeContentCargo
	Target   TradeContentTarget
	OrderNum int
	MsgID    MsgID
	Time     time.Time
}

type TradeContentTarget struct {
	Grade  int
	Kind   string
	Purple int
	Blue   int
	Orange int
}

type TradeContentCargo struct {
	CargoGrade int
	cargoNum   int
}

type MsgID struct {
	Publish string
	Wait    string
	Accept  string
}

type TradeUsers struct {
	name   string
	nameID string
	emoji  string
	time   time.Time
}

type TradeData struct {
	sync.Mutex
	TradeAddCross    map[int]chan *khl.ReactionAddContext
	TradeAddNum      map[int]chan *khl.ReactionAddContext
	TradeAddCheck    map[int]chan *khl.ReactionAddContext
	TradeDeleteCheck map[int]chan *khl.ReactionDeleteContext
}

var Trade = &TradeData{
	TradeAddCross:    make(map[int]chan *khl.ReactionAddContext, 1),
	TradeAddNum:      make(map[int]chan *khl.ReactionAddContext, 1),
	TradeAddCheck:    make(map[int]chan *khl.ReactionAddContext, 1),
	TradeDeleteCheck: make(map[int]chan *khl.ReactionDeleteContext, 1),
}

// 交易订单总函数
func TradeOrder(ctx *khl.TextMessageContext) {
	var trade TradeContent
	dict := map[string]TradeUsers{}
	chanDone := make(chan bool, 1)

	OrderNum++
	TradeGetData(ctx, &trade, OrderNum)

	Trade.Mutex.Lock()
	Trade.TradeAddCross[trade.OrderNum] = make(chan *khl.ReactionAddContext)
	Trade.TradeAddNum[trade.OrderNum] = make(chan *khl.ReactionAddContext)
	Trade.TradeAddCheck[trade.OrderNum] = make(chan *khl.ReactionAddContext)
	Trade.TradeDeleteCheck[trade.OrderNum] = make(chan *khl.ReactionDeleteContext)
	Trade.Mutex.Unlock()

	TradeFirstMessage(ctx, &trade)
	// 发布订单者先加入map中
	dict["1"] = TradeUsers{
		name:   ctx.Extra.Author.Username,
		nameID: ctx.Extra.Author.ID,
		time:   time.Now(),
	}

	for {
		select {
		case delete := <-Trade.TradeAddCross[trade.OrderNum]:
			ctx.Session.Logger.Warn().Interface("delete", delete).Msg("TradeOrderGoroutine")
			TradeOrderDelete(delete, trade, chanDone)
		case addNum := <-Trade.TradeAddNum[trade.OrderNum]:
			ctx.Session.Logger.Warn().Interface("addNum", addNum).Msg("TradeOrderGoroutine")
			TradeOrderAddNum(addNum, &trade, dict, chanDone)
		case acceptAdd := <-Trade.TradeAddCheck[trade.OrderNum]:
			ctx.Session.Logger.Warn().Interface("acceptAdd", acceptAdd).Msg("TradeOrderGoroutine")
			TradeOrderAcceptAdd(acceptAdd, trade, dict)
		case acceptDelete := <-Trade.TradeDeleteCheck[trade.OrderNum]:
			ctx.Session.Logger.Warn().Interface("acceptDelete", acceptDelete).Msg("TradeOrderGoroutine")
			TradeOrderAcceptDelete(acceptDelete, trade, dict)
		case <-chanDone:
			ctx.Session.Logger.Warn().Msg("TradeOrderGoroutine done")
			return
		}
	}
}

// 新加入的接单者获取一个数字emoji
func TradeOrderGetEmojiNum(dict map[string]TradeUsers) string {
	emoji := make(map[string]bool)

	for _, value := range dict {
		emoji[value.emoji] = true
	}

	for i := 0; i < 4; i++ {
		if _, ok := emoji[config.EmojiNum[i]]; !ok {
			return config.EmojiNum[i]
		}
	}

	return ""
}

// 订单发布者取消订单
func TradeOrderDelete(ctx *khl.ReactionAddContext, trade TradeContent, done chan bool) {
	err := ctx.Session.MessageDelete(trade.MsgID.Publish)
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderDelete Publish")
		return
	}
	err = ctx.Session.MessageDelete(trade.MsgID.Wait)
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderDelete Wait")
		return
	}

	done <- true
}

// 订单发布者选择订单接受者
func TradeOrderAddNum(ctx *khl.ReactionAddContext, trade *TradeContent, dict map[string]TradeUsers, done chan bool) {
	err := ctx.Session.MessageDelete(trade.MsgID.Publish)
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAddNum Publish")
		return
	}
	err = ctx.Session.MessageDelete(trade.MsgID.Wait)
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAddNum Wait")
		return
	}

	content := TradeGetBaseContent(*trade, fmt.Sprintf("(met)%s(met)", dict["1"].nameID), fmt.Sprintf("\\n(met)%s(met)", ctx.Extra.UserID))
	TradeSentMessage(trade, ctx.Session, config.Data.IDChannelTradeAccept, content)

	done <- true
}

// 接单者接单
func TradeOrderAcceptAdd(ctx *khl.ReactionAddContext, trade TradeContent, dict map[string]TradeUsers) {
	// 最多4人接单
	if len(dict) == 5 {
		return
	}

	user, err := ctx.Session.UserView(ctx.Extra.UserID, khl.UserViewWithGuildID(ctx.Common.TargetID))
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd UserView")
		return
	}

	dict[ctx.Extra.UserID] = TradeUsers{
		name:   user.Username,
		nameID: user.ID,
		emoji:  TradeOrderGetEmojiNum(dict),
		time:   time.Now(),
	}

	nameMaster, nameElse := TradeOrderGetNames(dict)
	content := TradeGetBaseContent(trade, nameMaster, nameElse)

	// update messge in publish and wait
	err = ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   trade.MsgID.Publish,
			Content: content,
		},
	})
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageUpdate")
		return
	}

	err = ctx.Session.MessageAddReaction(trade.MsgID.Publish, dict[ctx.Extra.UserID].emoji)
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageAddReaction")
		return
	}

	err = ctx.Session.MessageUpdate(&khl.MessageUpdate{
		MessageUpdateBase: khl.MessageUpdateBase{
			MsgID:   trade.MsgID.Wait,
			Content: content,
		},
	})
	if err != nil {
		ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageUpdate")
		return
	}
}

// 接单者取消接单
func TradeOrderAcceptDelete(ctx *khl.ReactionDeleteContext, trade TradeContent, dict map[string]TradeUsers) {
	// 检查用户是否在接单组中
	if _, ok := dict[ctx.Extra.UserID]; !ok {
		SendTempMessage(ctx.Session, ctx.Common.TargetID, fmt.Sprintf("(met)%s(met) 你没有接单，不能取消接单！", ctx.Extra.UserID))
		ctx.Session.Logger.Warn().Msgf("%s 你没有接单，不能取消接单！", ctx.Extra.UserID)
		return
	} else {
		err := ctx.Session.MessageDeleteReaction(trade.MsgID.Publish, dict[ctx.Extra.UserID].emoji, "")
		if err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageDeleteReaction")
			return
		}
		delete(dict, ctx.Extra.UserID)

		// 更新消息
		nameMaster, nameElse := TradeOrderGetNames(dict)
		content := TradeGetBaseContent(trade, nameMaster, nameElse)

		// update messge in publish and wait
		err = ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   trade.MsgID.Publish,
				Content: content,
			},
		})
		if err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageUpdate Publish")
			return
		}

		err = ctx.Session.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID:   trade.MsgID.Wait,
				Content: content,
			},
		})
		if err != nil {
			ctx.Session.Logger.Error().Err("", err).Msg("TradeOrderAcceptAdd MessageUpdate Wait")
			return
		}
	}
}

// 从map中获取主机和僚机信息
func TradeOrderGetNames(dict map[string]TradeUsers) (string, string) {
	dictTime := make(map[int64]TradeUsers)
	keys := []string{}
	namesList := []TradeUsers{}
	nameMaster := ""
	nameElse := ""

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
			namesList = append(namesList, TradeUsers{
				name:   dictTime[timeKeyInt].name,
				nameID: dictTime[timeKeyInt].nameID,
				emoji:  dictTime[timeKeyInt].emoji,
				time:   dictTime[timeKeyInt].time,
			})
		}
	}

	for i := 0; i < len(dict); i++ {
		if i == 0 {
			nameMaster += namesList[i].name
		} else {
			nameElse += fmt.Sprintf("\\n%s %s", namesList[i].emoji, namesList[i].name)
		}
	}

	return nameMaster, nameElse
}

// 获取订单的数据
func TradeGetData(ctx *khl.TextMessageContext, trade *TradeContent, orderNum int) {
	var targetKind string

	content := strings.Fields(RemovePrefix(ctx.Common.Content))
	// order 6s120 9scot
	if len(content) == 3 {
		cargo := strings.Split(content[1], "s")
		cargoGrade, _ := strconv.Atoi(cargo[0])
		cargoNum, _ := strconv.Atoi(cargo[1])
		target := strings.Split(content[2], "s")
		targetGrade, _ := strconv.Atoi(target[0])
		targetKind = target[1]

		trade.Cargo = []TradeContentCargo{{cargoGrade, cargoNum}}
		trade.Target.Grade = targetGrade
		trade.Target.Kind = targetKind

		// order 8s40 7s30 10scot
	} else if len(content) == 4 {
		cargo1 := strings.Split(content[1], "s")
		cargoGrade1, _ := strconv.Atoi(cargo1[0])
		cargoNum1, _ := strconv.Atoi(cargo1[1])
		cargo2 := strings.Split(content[2], "s")
		cargoGrade2, _ := strconv.Atoi(cargo2[0])
		cargoNum2, _ := strconv.Atoi(cargo2[1])
		target := strings.Split(content[3], "s")
		targetGrade, _ := strconv.Atoi(target[0])
		targetKind = target[1]

		trade.Cargo = []TradeContentCargo{{cargoGrade1, cargoNum1}, {cargoGrade2, cargoNum2}}
		trade.Target.Grade = targetGrade
		trade.Target.Kind = targetKind

	}
	trade.Time = time.Now()
	TradeGetTargetNum(trade)

	trade.OrderNum = orderNum
}

// 发送初始消息
func TradeFirstMessage(ctx *khl.TextMessageContext, trade *TradeContent) (resp *khl.MessageResp, err error) {

	content := TradeGetBaseContent(*trade, ctx.Extra.Author.Username, "")

	TradeSentMessage(trade, ctx.Session, ctx.Common.TargetID, content)
	TradeSentMessage(trade, ctx.Session, config.Data.IDChannelTradeWait, content)

	return
}

// 获取订单中消息内容
func TradeGetBaseContent(trade TradeContent, order string, accept string) string {
	contentCargo := ""

	for i := 0; i < len(trade.Cargo); i++ {
		contentCargo += fmt.Sprintf("%d×%ds ", trade.Cargo[i].cargoNum, trade.Cargo[i].CargoGrade)
	}
	contentTargetCount := trade.Target.Purple + trade.Target.Blue + trade.Target.Orange
	contentTarget := fmt.Sprintf("%d×%ds (%s)", contentTargetCount, trade.Target.Grade, TradeGetBaseContentKindNum(trade))
	createTime := trade.Time.Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(TextTradePublish, trade.OrderNum, order, accept, contentCargo, contentTarget, TradeGetDiscount(trade), createTime)

	return content
}

// 获取订单目标种类的数量
func TradeGetBaseContentKindNum(trade TradeContent) string {
	kindNum := ""

	if trade.Target.Purple > 0 {
		kindNum += fmt.Sprintf(":purple_heart:×%d  ", trade.Target.Purple)
	}
	if trade.Target.Blue > 0 {
		kindNum += fmt.Sprintf(":blue_heart:×%d  ", trade.Target.Blue)
	}
	if trade.Target.Orange > 0 {
		kindNum += fmt.Sprintf(":yellow_heart:×%d", trade.Target.Orange)
	}

	return kindNum
}

// 发送消息
func TradeSentMessage(trade *TradeContent, session *khl.Session, targetID string, content string) {
	resp, err := session.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeCard,
			TargetID: targetID,
			Content:  content,
		},
	})
	if err != nil {
		session.Logger.Error().Err("", err).Msg("TradeSentMessage MessageCreate")
		return
	}

	switch targetID {
	case config.Data.IDChannelTradePublish:
		trade.MsgID.Publish = resp.MsgID
		session.MessageAddReaction(resp.MsgID, config.Data.EmojiCrossMark)
	case config.Data.IDChannelTradeWait:
		trade.MsgID.Wait = resp.MsgID
		session.MessageAddReaction(resp.MsgID, config.Data.EmojiCheckMark)
	case config.Data.IDChannelTradeAccept:
		trade.MsgID.Accept = resp.MsgID
	}
}

// 获取折扣
func TradeGetDiscount(trade TradeContent) int {
	kind := TradeGetKind(trade.Target.Kind)

	switch kind {
	case 3:
		return 5
	case 2:
		return 2
	default:
		return 0
	}
}

// 计算能换取的目标数量
func TradeGetTargetNum(trade *TradeContent) {
	var count float64

	// 备份结构体trade
	tradeBak := *trade
	tradeBak.Cargo = make([]TradeContentCargo, len(trade.Cargo))
	copy(tradeBak.Cargo, trade.Cargo)

	for i := 0; i < len(trade.Cargo); i++ {
		if ((tradeBak.Target.Grade == 10) || (tradeBak.Target.Grade == 11)) && strings.Contains(trade.Target.Kind, "t") {
			switch TradeGetKind(trade.Target.Kind) {
			case 1:
				trade.Target.Purple += int(math.Round(float64(tradeBak.Cargo[i].cargoNum) / TradePrice[tradeBak.Target.Grade-4][tradeBak.Cargo[i].CargoGrade-4]))
			case 2:
				trade.Target.Purple += int(math.Round(float64(tradeBak.Cargo[i].cargoNum) / (TradePrice[tradeBak.Target.Grade-4][tradeBak.Cargo[i].CargoGrade-4] * 2)))
				tradeBak.Cargo[i].cargoNum -= int(TradePrice[tradeBak.Target.Grade-4][tradeBak.Cargo[i].CargoGrade-4] * float64(trade.Target.Purple))
			case 3:
				trade.Target.Purple += int(math.Round(float64(tradeBak.Cargo[i].cargoNum) / (TradePrice[tradeBak.Target.Grade-4][tradeBak.Cargo[i].CargoGrade-4] * 3)))
				tradeBak.Cargo[i].cargoNum -= int(TradePrice[tradeBak.Target.Grade-4][tradeBak.Cargo[i].CargoGrade-4] * float64(trade.Target.Purple))
			}
		}
		count += float64(tradeBak.Cargo[i].cargoNum) / TradePrice[tradeBak.Target.Grade-6][tradeBak.Cargo[i].CargoGrade-4]
	}
	discount := TradeGetDiscount(*trade)
	count = count * (1.0 + float64(discount+2)/100)
	countNum := int(math.Round(float64(count)))

	if ((tradeBak.Target.Grade == 10) || (tradeBak.Target.Grade == 11)) && strings.Contains(trade.Target.Kind, "t") {
		trade.Target.Purple = int(float64(trade.Target.Purple) * (1.0 + float64(discount+2)/100))
	}

	switch TradeGetKind(trade.Target.Kind) {
	case 1:
		if ((tradeBak.Target.Grade == 10) || (tradeBak.Target.Grade == 11)) && strings.Contains(trade.Target.Kind, "t") {
			break
		}
		switch trade.Target.Kind {
		case "t":
			trade.Target.Purple = countNum
		case "c":
			trade.Target.Blue = countNum
		case "o":
			trade.Target.Orange = countNum
		}
	case 2:
		for i := 0; i < len(trade.Target.Kind); i++ {
			if ((tradeBak.Target.Grade == 10) || (tradeBak.Target.Grade == 11)) && strings.Contains(trade.Target.Kind, "t") {
				switch string(trade.Target.Kind[i]) {
				case "c":
					trade.Target.Blue = countNum
				case "o":
					trade.Target.Orange = countNum
				}
			} else {
				if i == 0 {
					switch string(trade.Target.Kind[i]) {
					case "t":
						trade.Target.Purple += countNum / 2
					case "c":
						trade.Target.Blue += countNum / 2
					case "o":
						trade.Target.Orange += countNum / 2
					}
				} else {
					switch string(trade.Target.Kind[i]) {
					case "t":
						trade.Target.Purple = countNum - trade.Target.Purple - trade.Target.Blue - trade.Target.Orange
					case "c":
						trade.Target.Blue = countNum - trade.Target.Purple - trade.Target.Blue - trade.Target.Orange
					case "o":
						trade.Target.Orange = countNum - trade.Target.Purple - trade.Target.Blue - trade.Target.Orange
					}
				}
			}
		}
	case 3:
		if ((tradeBak.Target.Grade == 10) || (tradeBak.Target.Grade == 11)) && strings.Contains(trade.Target.Kind, "t") {
			for i := 0; i < countNum; i++ {
				switch i % 2 {
				case 0:
					trade.Target.Orange++
				case 1:
					trade.Target.Blue++
				}
			}
		} else {
			for i := 0; i < countNum; i++ {
				switch i % 3 {
				case 0:
					trade.Target.Purple++
				case 1:
					trade.Target.Blue++
				case 2:
					trade.Target.Orange++
				}
			}
		}
	}
}

// 获取换取目标的种类
func TradeGetKind(targetKind string) int {
	kind := 0
	if strings.Contains(targetKind, "t") {
		kind++
	}
	if strings.Contains(targetKind, "c") {
		kind++
	}
	if strings.Contains(targetKind, "o") {
		kind++
	}

	return kind
}
