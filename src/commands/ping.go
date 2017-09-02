package BotCommands

import (
	"../../src/util"
	"github.com/bwmarrin/discordgo"
	"time"
)

func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) {
	start := time.Now()
	message, err := s.ChannelMessageSend(m.ChannelID, "...")
	elapsed := time.Since(start)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageEdit(m.ChannelID, message.ID, "Ping! Time taken to send message: `"+elapsed.String()+"`.")
}
