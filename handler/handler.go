package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"github.com/eaok/ashe/print"
)

// delete not queue messages, just keep 20s
func AutoDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != s.State.User.ID || !strings.Contains(m.Content, "RED STAR QUEUE") {
		go func() {
			time.Sleep(20 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}()
	}
}

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, "pong")
	}
}

func Avatar(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "avatar" {
		s.ChannelMessageSend(m.ChannelID, m.Author.AvatarURL("2048"))
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, m.Author.AvatarURL("2048"))
	}
}

func Pic(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "pic" {
		s.ChannelMessageSend(m.ChannelID, "https://cdn.jsdelivr.net/gh/eaok/img/docker/components-of-kubernetes.png")
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, "https://cdn.jsdelivr.net/gh/eaok/img/docker/components-of-kubernetes.png")
	}
}

func Emoji(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "emoji" {
		s.ChannelMessageSend(m.ChannelID, ":one:")
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, ":one:")
	}
}

func Username(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "username" {
		s.ChannelMessageSend(m.ChannelID, "Your username is "+m.Message.Author.Username)
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, "Your username is "+m.Message.Author.Username)
	}
}

func Page(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "page" {
		p := dgwidgets.NewPaginator(s, m.ChannelID)

		// Add embed pages to paginator
		p.Add(&discordgo.MessageEmbed{Description: "Page one"},
			&discordgo.MessageEmbed{Description: "Page two"},
			&discordgo.MessageEmbed{Description: "Page three"})

		// Sets the footers of all added pages to their page numbers.
		p.SetPageFooters()

		// When the paginator is done listening set the colour to yellow
		p.ColourWhenDone = 0xffff

		// Stop listening for reaction events after five minutes
		p.Widget.Timeout = time.Minute * 5

		// Add a custom handler for the gun reaction.
		p.Widget.Handle("🔫", func(w *dgwidgets.Widget, r *discordgo.MessageReaction) {
			s.ChannelMessageSend(m.ChannelID, "Bang!")
		})

		p.Spawn()
	}
}

func Queue(s *discordgo.Session, m *discordgo.MessageCreate) {
	// log.Println(s.ShardID)
	// log.Println(s.State.Guilds)
	// log.Println(s.State.User.Bot)
	// log.Println(s.State.User.Username)
	// log.Println(m.Author)    // 发送人名字
	// log.Println(m.ChannelID) // 频道id 848379994279116831
	// log.Println(m.Content)
	// log.Println(m.GuildID) // 工会id 848081055029788673
	// log.Println(m.ID)      // 消息id 848535718790955008

	// var reaction *discordgo.MessageReaction
	// for {
	// 	k := <-nextMessageReactionAddC(s)
	// 	reaction = k.MessageReaction

	// 	// Ignore reactions sent by bot
	// 	if reaction.MessageID != m.Message.ID || s.State.User.ID == reaction.UserID {
	// 		continue
	// 	}

	// 	s.AddHandler(EmojiCheckMark, func(w *Widget, r *discordgo.MessageReaction) {
	// 		if err := p.NextPage(); err == nil {
	// 			p.Update()
	// 		}
	// 	})
	// }

	if m.Content == "queue" {
		text := fmt.Sprintf("**RED STAR QUEUE [%d/4]**\n", 1)
		text += fmt.Sprintf("Members joined:👇  |  🔴RS level: @RS%d\n", 9)
		text += fmt.Sprintf("Use ✅ to join, ❎ to leave or 🛑 to start in %d.\n", 2)
		text += "Bored while waiting? Type !guide to refresh your knowledge!\n"
		text += fmt.Sprintf("RS9 #%d•今天%d:%d", 44, time.Now().Hour(), time.Now().Minute())
		s.ChannelMessageSend(m.ChannelID, text)
		// s.MessageReactionAdd(m.ChannelID, m.Reference().MessageID, "❎")
		// s.MessageReactionAdd(m.ChannelID, m.Reference().MessageID, "🛑")

		fmt.Printf("%-10s\n", m.Content)
		print.ColorPrint(34, 0, text)
	} else if m.Author.ID == s.State.User.ID && strings.Contains(m.Content, "RED STAR QUEUE") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		// s.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		// 	log.Println(s.State.User.Username)
		// })
	}
}

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "help" {
		text := ""
		text += fmt.Sprintf("%-10s\t:\t%s\n", "ping", "responds with pong!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "avatar", "responds with your avatar!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "pic", "responds with a picture!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "emoji", "responds with a emoji!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "username", "responds with your name!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "page", "responds with a page text!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "queue", "responds with a queue!")
		text += fmt.Sprintf("%-10s\t:\t%s\n", "help", "prints this help menu!")

		s.ChannelMessageSend(m.ChannelID, "```"+text+"```")
		fmt.Printf("%-10s\n", m.Content)
		print.ColorPrint(34, 0, text)
	}
}
