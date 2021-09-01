package router

import (
	"fmt"
	"time"

	"github.com/eaok/ashe/config"
	"github.com/eaok/ashe/handler"
	"github.com/lonelyevil/khl"
)

func InitAction(s *khl.Session) {
	// 获取指定频道消息列表
	_, err := s.MessageList(config.Data.IDChannelSelectRole)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(msg[1].Reactions[0])

	// 往指定频道发送一条消息10秒后自动删除
	// resp, _ := s.MessageCreate(&khl.MessageCreate{
	// 	MessageCreateBase: khl.MessageCreateBase{
	// 		TargetID: handler.IDChannelSelectRole,
	// 		Content:  handler."EmojiTest",
	// 	},
	// })
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	s.MessageDelete(resp.MsgID)
	// }()

	resp, _ := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			Type:     khl.MessageTypeKMarkdown,
			TargetID: config.Data.IDChannelSelectRole,
			Content: `**选择红星车队等级**
			---
			根据自己红星等级，自助选择要参与红星车队的等级角色！
			:four:: RS4
			:five:: RS5
			:six:: RS6
			:seven:: RS7
			:eight:: RS8
			:nine:: RS9
			:keycap_ten:: RS10
			:restroom:: RS11
			---`,
		},
	})

	go func() {
		time.Sleep(5 * time.Second)
		s.MessageUpdate(&khl.MessageUpdate{
			MessageUpdateBase: khl.MessageUpdateBase{
				MsgID: resp.MsgID,
				Content: `**选择蓝星匹配等级**
				---
				根据自己蓝星战舰的分数，自助选择匹配级别！
				分数可通过 APP HS Compendium 查询，匹配并不是严格按照这个区间匹配的，请自己斟酌选择！
				组队规则：
				1. 主机前后区共计3个区间所在的战舰可以当僚机！
				2. 僚机请用装配避难模组的战舰！
				:one:: BS1 : 0-1000
				:two:: BS1 : 2000-2000
				:three:: BS3 : 2000-3000
				:four:: BS4 : 3000-4000
				:five:: BS5 : 4000-5000
				:six:: BS6 : 5000-6000
				:seven:: BS7 : 6000-7000
				:eight:: BS8 : 7000- ∞
				---`,
			},
		})
	}()
}

func Route(s *khl.Session) {
	s.AddHandler(handler.AutoDelete)
	s.AddHandler(handler.Ping)
	s.AddHandler(handler.AddReaction)
	s.AddHandler(handler.DeleteReaction)
	s.AddHandler(handler.Team)
	s.AddHandler(handler.Help)
}
