package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/dustin/go-humanize"
	"github.com/NamedKitten/discordgo"
	"github.com/NamedKitten/dstate"
	"time"
)

func init() {
	commands.RegisterCommand("userinfo", userinfoCommand)
	commands.RegisterHelp("userinfo", "Shows info on a user.")
}

func userinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	var member *discordgo.Member
	var memberState *dstate.MemberState

	guildState := commands.State.Guild(false, ctx.GuildID)

	if ctx.HasMention {
		memberState = guildState.Member(false, m.Mentions[0].ID)
	} else {
		memberState = guildState.Member(false, m.Author.ID)
	}
	member = memberState.DGoCopy()

	var game string
	var status string

	presenceStatus := memberState.PresenceStatus
	presenceGame := memberState.PresenceGame
	switch {
	case presenceGame == nil:
		game = "None"
	case presenceGame.Type == 0:
		game = "Playing " + presenceGame.Name
	case presenceGame.Type == 1:
		game = "Streaming " + presenceGame.Name
	}

	switch presenceStatus {
	default:
		status = "Offline"
	case 1:
		status = "Online"
	case 2:
		status = "Idle"
	case 3:
		status = "Do Not Disturb"
	}

	timenow := time.Now()
	_, zone := timenow.Zone()
	joined, _ := discordgo.Timestamp(member.JoinedAt).Parse()
	userSnowflake := member.User.ID
	joinedDiscord := time.Unix((((userSnowflake>>22)+1420070400000)/1000)+int64(zone), 0)
	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Joined At**:", Value: fmt.Sprintf("**This Server**: %s\n**Discord**: %s",
		humanize.Time(joined),
		humanize.Time(joinedDiscord),
	)})
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    m.Author.Username,
			IconURL: m.Author.AvatarURL("512"),
		},
		Description: fmt.Sprintf("**Display Name**: %s\n**ID**: %d\n[**Avatar URL**](%s)\n**Currently Playing**: %s\n**Status**: %s",
			member.User.Username,
			member.User.ID,
			member.User.AvatarURL(""),
			game, status),
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
