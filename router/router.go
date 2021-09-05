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
	msg, err := s.MessageList(config.Data.IDChannelSelectRole)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msg[0])

	// 程序启动时往指定频道发送一条消息
	resp, _ := s.MessageCreate(&khl.MessageCreate{
		MessageCreateBase: khl.MessageCreateBase{
			TargetID: config.Data.IDChannelSelectRole,
			Content:  "ashe is start...",
		},
	})
	go func() {
		time.Sleep(5 * time.Minute)
		s.MessageDelete(resp.MsgID)
	}()

	// 启动红星频道排队
	handler.TeamStart(s)
}

func Route(s *khl.Session) {
	if config.Data.RunMod == "debug" {
		s.AddHandler(handler.AutoDelete)
	}
	s.AddHandler(handler.Ping)
	s.AddHandler(handler.CardButton)
	s.AddHandler(handler.AddReaction)
	s.AddHandler(handler.DeleteReaction)
	s.AddHandler(handler.InTeam)
	s.AddHandler(handler.OutTeam)
	s.AddHandler(handler.Help)
}
