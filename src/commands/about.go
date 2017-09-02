package BotCommands

import (
	"../config"
	"../util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"runtime"
)

func AboutCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) {
	fields := make([]*discordgo.MessageEmbedField, 0, 3)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**KittehBotGo Version**:", Value: config.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Go Version**:", Value: runtime.Version(), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**DiscordGo Version**:", Value: discordgo.VERSION, Inline: true})

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "About " + s.State.User.Username,
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", s.State.User.ID, s.State.User.Avatar),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGO!",
		},
	})
}
