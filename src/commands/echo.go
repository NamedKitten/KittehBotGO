package BotCommands

import (
	"../../src/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-errors/errors"
)

func EchoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) (error) {
	content := ctx.Content
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	member, err := s.State.Member(channel.GuildID, m.Author.ID)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", m.Author.ID, m.Author.Avatar),
		},
		Color:       s.State.UserColor(m.Author.ID, m.ChannelID),
		Description: content,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGO!",
		},
	})
	return nil

}
