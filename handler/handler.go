package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"github.com/eaok/ashe/config"
	"github.com/eaok/ashe/print"
	"github.com/eaok/ashe/widget"
)

// add order prefix
func RemovePrefix(m *discordgo.MessageCreate) string {
	if config.Prefix != "" && strings.HasPrefix(m.Content, config.Prefix) {
		return strings.TrimPrefix(m.Content, config.Prefix)
	}

	return ""
}

// delete not queue messages, just keep 20s
func AutoDelete(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if m.Author.ID != s.State.User.ID || !strings.Contains(m.Content, "RED STAR QUEUE") {
	if m.Author.ID != s.State.User.ID {
		fmt.Println("queue ", m.Content)
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
	if RemovePrefix(m) == "ping" {
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

	if RemovePrefix(m) == "avatar" {
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

	if RemovePrefix(m) == "pic" {
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

	if RemovePrefix(m) == "emoji" {
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

	if RemovePrefix(m) == "username" {
		s.ChannelMessageSend(m.ChannelID, "Your username is "+m.Message.Author.Username)
		fmt.Printf("%-10s\t", m.Content)
		print.ColorPrint(34, 0, "Your username is "+m.Message.Author.Username)
	}
}

func Page(s *discordgo.Session, m *discordgo.MessageCreate) {
	if RemovePrefix(m) == "page" {
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

	if RemovePrefix(m) == "queue" {
		go func() {
			widget.BeginChan <- 1
		}()

		for {
			if <-widget.BeginChan == 1 {
				go func() {
					p := widget.NewPaginator(s, m.ChannelID)

					// Add embed pages to paginator
					p.Add(
						&discordgo.MessageEmbed{
							Description: fmt.Sprintf(widget.Text, 0, 9, "", 0),
						},
						&discordgo.MessageEmbed{
							Description: widget.Text,
						},
						&discordgo.MessageEmbed{
							Description: widget.Text,
						},
					)
					p.SetPageFooters()
					p.Spawn()
				}()
			}
		}
	}
}

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if RemovePrefix(m) == "help" {
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
