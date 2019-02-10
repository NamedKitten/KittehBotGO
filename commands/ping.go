package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/go-errors/errors"
	"github.com/jonas747/discordgo"
	"time"
)

func init() {
	commands.RegisterCommand("ping", pingCommand)
	commands.RegisterHelp("ping", "Shows latency for sending a message.")
}

func pingCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	start := time.Now()
	message, err := s.ChannelMessageSend(m.ChannelID, "...")
	elapsed := time.Since(start)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	s.ChannelMessageEdit(m.ChannelID, message.ID, "Pong! "+elapsed.String())
	return nil
}
