package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"strconv"
	"time"
)

func init() {
	commands.RegisterCommand("userinfo", UserinfoCommand)
}

func UserinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
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
		status = "Offline"
		game = "None"
	} else {
		switch {
		case presence.Game == nil:
			game = "None"
		case presence.Game.Type == 0:
			game = "Playing " + presence.Game.Name
		case presence.Game.Type == 0:
			game = "Streaming " + presence.Game.Name
		}

		switch string(presence.Status) {
		case "dnd":
			status ="Do Not Disturb"
		case "online":
			status = "Offline"
		case "idle":
			status = "Idle"
		}
	}
	print(status)

	timenow := time.Now()
	_, zone := timenow.Zone()
	joined, _ := discordgo.Timestamp(member.JoinedAt).Parse()
	userSnowflake, _ := strconv.ParseInt(member.User.ID, 10, 64)
	joinedDiscord := time.Unix((((userSnowflake>>22)+1420070400000)/1000)+int64(zone), 0)
	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Joined At**:", Value: fmt.Sprintf("**%s**: %s\n**Discord**: %s",
		"This Server",
		humanize.Time(joined),
		humanize.Time(joinedDiscord),
	)})
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL("512"),
		},
		Description: fmt.Sprintf("**%s**: %s\n**ID**: %s\n[**%s**](%s)\n**%s**: %s",
			"Display Name",
			member.User.Username,
			member.User.ID,
			"Avatar URL",
			member.User.AvatarURL(""),
			"Currently Playing",
			game),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("512"),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGo",
		},
	})
	return nil

}
