package router

import (
	"github.com/bwmarrin/discordgo"
	"github.com/eaok/ashe/handler"
)

func Route(s *discordgo.Session) {
	s.AddHandler(handler.AutoDelete)
	s.AddHandler(handler.Ping)
	s.AddHandler(handler.Avatar)
	s.AddHandler(handler.Pic)
	s.AddHandler(handler.Emoji)
	s.AddHandler(handler.Username)
	s.AddHandler(handler.Page)
	s.AddHandler(handler.Queue)
	s.AddHandler(handler.Help)
}
