package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	"strconv"
)

func init() {
	commands.RegisterCommand("letmein", LetMeInCommand)
}

func waitForMessage(s *discordgo.Session) chan *discordgo.MessageCreate {
	out := make(chan *discordgo.MessageCreate)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		out <- e
	})
	return out
}


func LetMeInCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	application, _ := s.Application("@me")
	
	if application.Owner.ID != m.Author.ID {
		return nil
	}

	message := "```md\n"
	for i, g := range s.State.Guilds {
		message += fmt.Sprintf("[%d](%s)\n", i, g.Name)
	}
	message += "\n```What guild do you wish to generate a invite for?"

	s.ChannelMessageSend(m.ChannelID, message)

	for {
		wm := <-waitForMessage(s)
		if m.ChannelID != wm.ChannelID {
			continue
		} else if m.Author.ID != wm.Author.ID {
			continue
		} else {
			if wm.Content == "exit" || wm.Content == "quit" {
				s.ChannelMessageSend(m.ChannelID, "Exiting.")
			}
			num, err := strconv.Atoi(wm.Content)
			s.ChannelMessageSend(m.ChannelID, strconv.Itoa(num))
			s.ChannelMessageSend(m.ChannelID,strconv.Itoa(len(s.State.Guilds)) )
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid number.")
				return nil
			}


			guild := s.State.Guilds[num]
			for _, c := range guild.Channels {
				if c.Type == 0 {
					invite, ierr := s.ChannelInviteCreate(c.ID, discordgo.Invite{MaxUses: 0, MaxAge: 0, Temporary: false})
					if ierr != nil {
						fmt.Println("Can't create invite.")
						return nil
					} else {
						s.ChannelMessageSend(m.ChannelID, "https://discord.gg/" + invite.Code)
						return nil
					}
				} else {
					fmt.Println(c)
				}
			}
			s.ChannelMessageSend(m.ChannelID, "Can't fetch invite.")
			return nil	

		}

	}

	return nil
}
