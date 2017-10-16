package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
)

func LanguageCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()
	if len(ctx.Args) > 1 {
		if ctx.Args[0] == "set" {
			//if guild.OwnerID != m.Author.ID {
			//	s.ChannelMessageSend(m.ChannelID, "You need to be the owner to use this command.")
			//	return nil
			//} else {
				ctx.Commands.Redis.Set("language_"+ctx.GuildID, ctx.Args[1], 0)
				s.ChannelMessageSend(m.ChannelID, ctx.T("command_language_success"))
			//}
		}
	}
	return nil
}
