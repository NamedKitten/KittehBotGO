package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/go-errors/errors"
	// "runtime/debug"
	"time"
)

func init() {
	commands.RegisterCommand("ping", PingCommand)
}

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	//defer debug.FreeOSMemory()
	start := time.Now()
	message, err := s.ChannelMessageSend(m.ChannelID, "...")
	elapsed := time.Since(start)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	s.ChannelMessageEdit(m.ChannelID, message.ID, ctx.T("command_ping_result", struct{ Elapsed string }{Elapsed: elapsed.String()}))
	return nil
}
