package widget

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Paginator provides a method for creating a navigatable embed
type Paginator struct {
	sync.Mutex
	Pages []*discordgo.MessageEmbed
	Index int

	// Loop back to the beginning or end when on the first or last page.
	Loop   bool
	Widget *Widget

	Ses *discordgo.Session

	DeleteMessageWhenDone   bool
	DeleteReactionsWhenDone bool
	ColourWhenDone          int

	running bool
}

// NewPaginator returns a new Paginator
//    ses      : discordgo session
//    channelID: channelID to spawn the paginator on
func NewPaginator(ses *discordgo.Session, channelID string) *Paginator {
	p := &Paginator{
		Ses:                     ses,
		Pages:                   []*discordgo.MessageEmbed{},
		Index:                   0,
		Loop:                    false,
		DeleteMessageWhenDone:   false,
		DeleteReactionsWhenDone: false,
		ColourWhenDone:          -1,
		Widget:                  NewWidget(ses, channelID, nil),
	}
	p.addHandlers()

	return p
}

func GetNames() string {
	dictSort := make(map[int64]string)
	keys := []string{}
	namesList := []string{}

	for key, user := range dict {
		dictSort[user.time] = key
	}

	// 按时间排序取名字列表
	for timeKey := range dictSort {
		println("timeKey", timeKey)
		keys = append(keys, strconv.FormatInt(timeKey, 10))
	}
	sort.Strings(keys)

	for _, key := range keys {
		println("key", key)
		i, err := strconv.ParseInt(key, 10, 64)
		if err == nil {
			namesList = append(namesList, dict[dictSort[i]].name)
		}
	}

	names := ""
	for i := 0; i < len(dict); i++ {
		names += EmojiNum[i+1]
		names += " "
		names += namesList[i]
		names += "\n"
	}

	return names
}

func (p *Paginator) addHandlers() {
	// ✅
	p.Widget.Handle(EmojiCheckMark, func(w *Widget, r *discordgo.MessageReaction) {

		if user, ok := dict[r.UserID]; ok {
			// 已经在队列中则回显消息并2秒后删除
			msg, _ := w.Ses.ChannelMessageSendEmbed(w.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf("@%s You're already in the queue!", user.name),
			})
			go func() {
				time.Sleep(2 * time.Second)
				w.Ses.ChannelMessageDelete(w.ChannelID, msg.ID)
			}()
			return
		} else {
			// 加入队列
			userName, _ := w.Ses.User(r.UserID)
			dict[r.UserID] = Users{
				index: len(dict),
				name:  userName.Username,
				time:  time.Now().Unix(),
			}

			n := len(dict)
			names := GetNames()

			p.Pages[n].Description = fmt.Sprintf(Text, n, 9, names, n)
			fmt.Printf(p.Pages[n].Description, n, 9, names, n)
		}

		if err := p.NextPage(); err == nil {
			p.Update()
		}

		// 人数已满
		if len(dict) == 2 {
			for key := range dict {
				delete(dict, key)
			}

			w.Ses.ChannelMessageSendEmbed(w.ChannelID, &discordgo.MessageEmbed{
				Description: "go go go!",
			})

			BeginChan <- 1
		}
	})

	// ❎
	p.Widget.Handle(EmojiCrossMark, func(w *Widget, r *discordgo.MessageReaction) {
		if _, ok := dict[r.UserID]; ok {
			// 已经在队列中则离开队列
			delete(dict, r.UserID)
			n := len(dict)
			names := GetNames()

			p.Pages[n].Description = fmt.Sprintf(Text, n, 9, names, n)
			fmt.Printf(p.Pages[0].Description, n, 9, names, n)

			// 回显消息并1秒后删除
			userName, _ := w.Ses.User(r.UserID)
			msg, _ := w.Ses.ChannelMessageSendEmbed(w.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf("@%s You has left the queue!", userName.Username),
			})
			go func() {
				time.Sleep(1 * time.Second)
				w.Ses.ChannelMessageDelete(w.ChannelID, msg.ID)
			}()
		} else {
			// 不在队列中则回显消息并1秒后删除
			userName, _ := w.Ses.User(r.UserID)
			msg, _ := w.Ses.ChannelMessageSendEmbed(w.ChannelID, &discordgo.MessageEmbed{
				Description: fmt.Sprintf("@%s You're not in the queue!", userName.Username),
			})
			go func() {
				time.Sleep(1 * time.Second)
				w.Ses.ChannelMessageDelete(w.ChannelID, msg.ID)
			}()
			return
		}
		if err := p.PreviousPage(); err == nil {
			p.Update()
		}
	})

	// 🛑
	p.Widget.Handle(EmojiStopSign, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.Goto(0); err == nil {
			p.Update()
		}

		w.Ses.ChannelMessageSendEmbed(w.ChannelID, &discordgo.MessageEmbed{
			Description: "go go go!",
		})

		BeginChan <- 1
	})
}

// Spawn spawns the paginator in channel p.ChannelID
func (p *Paginator) Spawn() error {
	if p.Running() {
		return ErrAlreadyRunning
	}
	p.Lock()
	p.running = true
	p.Unlock()

	defer func() {
		p.Lock()
		p.running = false
		p.Unlock()

		// Delete Message when done
		if p.DeleteMessageWhenDone && p.Widget.Message != nil {
			p.Ses.ChannelMessageDelete(p.Widget.Message.ChannelID, p.Widget.Message.ID)
		} else if p.ColourWhenDone >= 0 {
			if page, err := p.Page(); err == nil {
				page.Color = p.ColourWhenDone
				p.Update()
			}
		}

		// Delete reactions when done
		if p.DeleteReactionsWhenDone && p.Widget.Message != nil {
			p.Ses.MessageReactionsRemoveAll(p.Widget.ChannelID, p.Widget.Message.ID)
		}
	}()

	page, err := p.Page()
	if err != nil {
		return err
	}
	p.Widget.Embed = page

	return p.Widget.Spawn()
}

// Add a page to the paginator
//    embed: embed page to add.
func (p *Paginator) Add(embeds ...*discordgo.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

// Page returns the page of the current index
func (p *Paginator) Page() (*discordgo.MessageEmbed, error) {
	p.Lock()
	defer p.Unlock()

	if p.Index < 0 || p.Index >= len(p.Pages) {
		return nil, ErrIndexOutOfBounds
	}

	return p.Pages[p.Index], nil
}

// NextPage sets the page index to the next page
func (p *Paginator) NextPage() error {
	p.Lock()
	defer p.Unlock()

	if p.Index+1 >= 0 && p.Index+1 < len(p.Pages) {
		p.Index++
		return nil
	}

	return ErrIndexOutOfBounds
}

// PreviousPage sets the current page index to the previous page.
func (p *Paginator) PreviousPage() error {
	p.Lock()
	defer p.Unlock()

	if p.Index-1 >= 0 && p.Index-1 < len(p.Pages) {
		p.Index--
		return nil
	}

	// Set the queue back to the beginning if Loop is enabled.
	if p.Loop {
		p.Index = len(p.Pages) - 1
		return nil
	}

	return ErrIndexOutOfBounds
}

// Goto jumps to the requested page index
//    index: The index of the page to go to
func (p *Paginator) Goto(index int) error {
	p.Lock()
	defer p.Unlock()
	if index < 0 || index >= len(p.Pages) {
		return ErrIndexOutOfBounds
	}
	p.Index = index
	return nil
}

// Update updates the message with the current state of the paginator
func (p *Paginator) Update() error {
	if p.Widget.Message == nil {
		return ErrNilMessage
	}

	page, err := p.Page()
	if err != nil {
		return err
	}

	_, err = p.Widget.UpdateEmbed(page)
	p.Widget.SetPageReactions(page)
	return err
}

// Running returns the running status of the paginator
func (p *Paginator) Running() bool {
	p.Lock()
	running := p.running
	p.Unlock()
	return running
}

// SetPageFooters sets the footer of each embed to
// Be its page number out of the total length of the embeds.
func (p *Paginator) SetPageFooters() {
	for _, embed := range p.Pages {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("RS9 #%d•今天%d:%d", 44, time.Now().Hour(), time.Now().Minute()),
		}
	}
}
