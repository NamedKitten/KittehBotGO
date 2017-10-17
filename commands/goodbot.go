package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	"strconv"
	//"github.com/go-errors/errors"
)

func GoodBotCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	count, _ := commands.Redis.Get("GoodBot").Result()
	intcount, _ := strconv.ParseInt(count, 10, 64)
	intcount += 1

	go commands.Redis.Set("GoodBot", fmt.Sprintf("%d", intcount), 0)
	go s.ChannelMessageSend(m.ChannelID, ctx.T("command_goodbot_thanks", intcount))

	return nil

}
