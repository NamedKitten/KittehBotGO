package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"runtime/debug"
	"strings"
	//"github.com/go-errors/errors"
)

func init() {
	commands.RegisterCommand("echo", EchoCommand)
}

func EchoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()
	go s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", m.Author.ID, m.Author.Avatar),
		},
		Color:       s.State.UserColor(m.Author.ID, m.ChannelID),
		Description: strings.Join(ctx.Args[0:], " "),
		Footer: &discordgo.MessageEmbedFooter{
			Text: ctx.T("command_about_thanks"),
		},
	})
	return nil

}
