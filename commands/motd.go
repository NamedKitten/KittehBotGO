package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
    "github.com/jonas747/discordgo"
    "fmt"
	"runtime/debug"
	"strings"
	"strconv"
)

func init() {
	commands.RegisterCommand("motd", MotdCommand)
	commands.RegisterHelp("motd", "Message of the day related commands.")
	commands.Discord.AddHandler(MotdEvent)

}

func MotdCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	guild := commands.State.Guild(false, ctx.GuildID).Guild
	if len(ctx.Args) > 1 {
		if ctx.Args[0] == "set" {
			if guild.OwnerID != m.Author.ID {
				s.ChannelMessageSend(m.ChannelID, "You need to be the guild owner to use this command")
				return nil
			}
			if ctx.Args[1] == "channel" {
				go commands.Redis.Set("motd_"+string(ctx.GuildID)+"_channel", fmt.Sprintf("%v", ctx.ChannelID), 0, 0, false, false)
				s.ChannelMessageSend(m.ChannelID, "Channel set.")
			} else {
				go commands.Redis.Set("motd_"+string(ctx.GuildID), strings.Join(ctx.Args[1:], " "), 0, 0, false, false)
				s.ChannelMessageSend(m.ChannelID, "MOTD set.")
			}
		} else {
			motd, err := commands.Redis.Get("motd_" + string(ctx.GuildID))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Please set the MOTD.")
			} else {
				s.ChannelMessageSend(m.ChannelID, string(motd[:]))
			}
		}
	}
	return nil
}

func MotdEvent(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
    motdR, err := commands.Redis.Get("motd_" + string(m.GuildID))
    motd := string(motdR[:])
    
	motdchannelR, channelerr := commands.Redis.Get("motd_" + string(m.GuildID) + "_channel")
    motdchannel := string(motdchannelR[:])
	if err != nil || channelerr != nil {
		return
	} else {
		i, _ := strconv.ParseInt(motdchannel, 10, 64)
		go s.ChannelMessageSend(i, motd)
	}
}
