package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	//"github.com/go-errors/errors"
	//"github.com/NamedKitten/KittehBotGo/util/i18n"
	"strconv"
	"time"
)

func init() {
	commands.RegisterCommand("userinfo", UserinfoCommand)
}

func UserinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	//channel, err := s.State.Channel(m.ChannelID)
	//if err != nil {
	//	return errors.Wrap(err, 1)
	//}
	var member *discordgo.Member
	if ctx.HasMention {
		member, _ = s.State.Member(ctx.GuildID, m.Mentions[0].ID)
	} else {
		member, _ = s.State.Member(ctx.GuildID, m.Author.ID)
	}
	var game string
	var status string
	presence, error := s.State.Presence(ctx.GuildID, m.Author.ID)
	if error != nil {
		status = ctx.T("command_userinfo_offline")
		game = ctx.T("command_userinfo_none")
	} else {
		if presence.Game == nil {
			game = "None"
		} else if presence.Game.Type == 0 {
			game = ctx.T("command_userinfo_playing", struct{ Game string }{Game: presence.Game.Name})
		} else if presence.Game.Type == 1 {
			game = ctx.T("command_userinfo_streaming", struct{ Game string }{Game: presence.Game.Name})
		}

		if presence.Status == "dnd" {
			status = ctx.T("command_userinfo_donotdisturb")
		} else if presence.Status == discordgo.StatusOnline {
			status = ctx.T("command_userinfo_online")
		} else if presence.Status == discordgo.StatusIdle {
			status = ctx.T("command_userinfo_idle")
		} // else {
		//	status = "Offline"
		//}
		print(status)
	}

	timenow := time.Now()
	_, zone := timenow.Zone()
	joined, _ := discordgo.Timestamp(member.JoinedAt).Parse()
	userSnowflake, _ := strconv.ParseInt(member.User.ID, 10, 64)
	joinedDiscord := time.Unix((((userSnowflake>>22)+1420070400000)/1000)+int64(zone), 0)
	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_userinfo_joindates") + "**:", Value: fmt.Sprintf("**%s**: %s\n**Discord**: %s",
		ctx.T("command_userinfo_thisserver"), 
		humanize.Time(joined),
		humanize.Time(joinedDiscord),
	)})
	//TODO: multilingual version of gohumanize
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL("512"),
		},
		Description: fmt.Sprintf("**%s**: %s\n**ID**: %s\n[**%s**](%s)\n**%s**: %s",
		ctx.T("command_userinfo_displayname"),
		member.User.Username,
		member.User.ID,
		ctx.T("command_userinfo_avatar"),
		member.User.AvatarURL(""),
		ctx.T("command_userinfo_game"),
		game),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("512"),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: ctx.T("command_about_thanks"),
		},
	})
	return nil

}
