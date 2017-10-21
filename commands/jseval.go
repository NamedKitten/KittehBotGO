package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	//"github.com/go-errors/errors"
	"github.com/robertkrimen/otto"
	"strings"
)

var vm = otto.New()

func init() {
	commands.RegisterCommand("jseval", JSEvalCommand)
}
func JSEvalCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	application, _ := s.Application("@me")

	if application.Owner.ID != m.Author.ID {
		go s.ChannelMessageSend(m.ChannelID, ctx.T("command_jseval_ownererror"))
		return nil
	}

	if len(ctx.Args) > 0 {

		stuff := strings.Join(ctx.Args[0:], " ")

		value, err := vm.Run(stuff)
		if err != nil {
			go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s:\n```js\n%s\n```\n%s:\n```js\n%s\n```", ctx.T("command_jseval_input"), stuff, ctx.T("command_jseval_output"), fmt.Sprint(err)))
		} else {
			go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s:\n```js\n%s\n```\n%s:\n```js\n%s\n```", ctx.T("command_jseval_input"), stuff, ctx.T("command_jseval_output"), value))
		}
	}

	return nil

}
