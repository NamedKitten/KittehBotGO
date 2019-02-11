package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/database"
	"github.com/NamedKitten/discordgo"
	"runtime/debug"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand("motd", motdCommand)
	commands.RegisterHelp("motd", "Message of the day related commands.")
	commands.Discord.AddHandler(motdEvent)

}

func motdCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	guild := commands.State.Guild(false, ctx.GuildID).Guild
	if len(ctx.Args) > 1 {
		if ctx.Args[0] == "set" {
			if guild.OwnerID != m.Author.ID {
				s.ChannelMessageSend(m.ChannelID, "You need to be the guild owner to use this command")
				return nil
			}
			if ctx.Args[1] == "channel" {
				go database.Set("motd_"+string(ctx.GuildID)+"_channel", fmt.Sprintf("%v", ctx.ChannelID))
				s.ChannelMessageSend(m.ChannelID, "Channel set.")
			} else {
				go database.Set("motd_"+string(ctx.GuildID), strings.Join(ctx.Args[1:], " "))
				s.ChannelMessageSend(m.ChannelID, "MOTD set.")
			}
		} else {
			motd := database.Get("motd_" + string(ctx.GuildID))
			if len(motd) == 0 {
				s.ChannelMessageSend(m.ChannelID, "Please set the MOTD.")
			} else {
				s.ChannelMessageSend(m.ChannelID, motd)
			}
		}
	}
	return nil
}

func motdEvent(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	motd := database.Get("motd_" + string(m.GuildID))

	motdchannel := database.Get("motd_" + string(m.GuildID) + "_channel")
	i, _ := strconv.ParseInt(motdchannel, 10, 64)
	go s.ChannelMessageSend(i, motd)
}
