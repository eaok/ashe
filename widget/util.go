package widget

import (
	"github.com/bwmarrin/discordgo"
)

type Users struct {
	index int
	name  string
	time  int64
}

var (
	// dict  = make(map[string]string, 15)
	dict      = map[string]Users{}
	text1     = "**RED STAR QUEUE [%d/4]**\n"
	text2     = "Members joined:👇  |  🔴RS level: @RS%d\n"
	text3     = "%s"
	text4     = "Use ✅ to join, ❎ to leave or 🛑 to start in %d.\n"
	text5     = "Bored while waiting? Type !guide to refresh your knowledge!\n"
	Text      = text1 + text2 + text3 + text4 + text5
	BeginChan = make(chan int)
)

// NextMessageReactionAddC returns a channel for the next MessageReactionAdd event
func nextMessageReactionAddC(w *Widget) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	w.Ses.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}

// EmbedsFromString splits a string into a slice of MessageEmbeds.
//     txt     : text to split
//     chunklen: How long the text in each embed should be
//               (if set to 0 or less, it defaults to 2048)
func EmbedsFromString(txt string, chunklen int) []*discordgo.MessageEmbed {
	if chunklen <= 0 {
		chunklen = 2048
	}

	embeds := []*discordgo.MessageEmbed{}
	for i := 0; i < int((float64(len(txt))/float64(chunklen))+0.5); i++ {
		start := i * chunklen
		end := start + chunklen
		if end > len(txt) {
			end = len(txt)
		}
		embeds = append(embeds, &discordgo.MessageEmbed{
			Description: txt[start:end],
		})
	}
	return embeds
}
