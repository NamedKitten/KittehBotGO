package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	//"github.com/go-errors/errors"
	"github.com/Shopify/go-lua"
	"strings"
)

var l = lua.NewState()

func LuaEvalCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()
	lua.OpenLibraries(l)

	application, _ := s.Application("@me")

	if application.Owner.ID != m.Author.ID {
		go s.ChannelMessageSend(m.ChannelID, ctx.T("command_luaeval_ownererror"))
		return nil
	}

	if len(ctx.Args) > 0 {

		stuff := strings.Join(ctx.Args[0:], " ")

		err := lua.DoString(l, stuff)
		if err != nil {			
			go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("$s:\n```lua\n%s\n```\nOutput:\n```lua\n%s\n```", stuff, err))

		} else {
			go s.ChannelMessageSend(m.ChannelID, ctx.T("command_luaeval_done"))
		}
	}

	return nil

}
