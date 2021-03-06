package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/discordgo"
	"strings"
	//"github.com/go-errors/errors"
)

func init() {
	commands.RegisterCommand("echo", echoCommand)
	commands.RegisterHelp("echo", "Echos what you say.")
}

func echoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	go s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", m.Author.ID, m.Author.Avatar),
		},
		Color:       s.State.UserColor(m.Author.ID, m.ChannelID),
		Description: strings.Join(ctx.Args[0:], " "),
	})
	return nil

}
