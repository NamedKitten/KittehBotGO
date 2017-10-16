package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	"strings"
	//"github.com/go-errors/errors"
)

func MotdCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()
	prefix, _ := ctx.Commands.Redis.Get("prefix").Result()

	guild, _ := s.State.Guild(ctx.GuildID)
	if len(ctx.Args) > 1 {
		if ctx.Args[0] == "set" {
			if guild.OwnerID != m.Author.ID {
				s.ChannelMessageSend(m.ChannelID, ctx.T("command_motd_ownererror"))
				return nil
			}
			if ctx.Args[1] == "channel" {
				go ctx.Commands.Redis.Set("motd_"+ctx.GuildID+"_channel", ctx.ChannelID, 0)
				s.ChannelMessageSend(m.ChannelID, ctx.T("command_motd_channel", struct{ ChannelID string }{ChannelID: ctx.ChannelID}))
			} else {
				go ctx.Commands.Redis.Set("motd_"+ctx.GuildID, strings.Join(ctx.Args[1:], " "), 0)
				s.ChannelMessageSend(m.ChannelID, ctx.T("command_motd_set"))
			}
		} else {
			motd, err := ctx.Commands.Redis.Get("motd_" + ctx.GuildID).Result()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ctx.T("command_motd_unsetwarning", struct{ Prefix string }{Prefix: prefix}))
			} else {
				s.ChannelMessageSend(m.ChannelID, motd)
			}
		}
	}
	return nil
}
