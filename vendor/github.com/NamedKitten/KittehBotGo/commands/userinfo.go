package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"strconv"
	"time"
	"github.com/go-errors/errors"
)

func UserinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) (error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return errors.Wrap(err, 1)
	}
	var member *discordgo.Member
	if ctx.HasMention {
		member, _ = s.State.Member(channel.GuildID, m.Mentions[0].ID)
	} else {
		member, _ = s.State.Member(channel.GuildID, m.Author.ID)
	}

	var game string
	var status string
	presence, err := s.State.Presence(channel.GuildID, member.User.ID)
	if err != nil {
		status = "Offline"
		game = "None"
	} else {
		fmt.Println(presence.Game)
		if presence.Game == nil {
			game = "None"
		} else if presence.Game.Type == 0 {
			game = "Playing " + presence.Game.Name
		} else if presence.Game.Type == 1 {
			game = "Streaming " + presence.Game.Name
		}
		fmt.Println(presence.Status)
		if presence.Status == "dnd" {
			status = "Do Not Disturb"
		} else if presence.Status == discordgo.StatusOnline {
			status = "Online"
		} else if presence.Status == discordgo.StatusIdle {
			status = "Idle"
		} else {
			status = "Offline"
		}
		fmt.Println(status)
	}

	timenow := time.Now()
	_, zone := timenow.Zone()
	fmt.Println(zone)
	joined, _ := discordgo.Timestamp(member.JoinedAt).Parse()
	userSnowflake, _ := strconv.ParseInt(member.User.ID, 10, 64)
	joinedDiscord := time.Unix((((userSnowflake >> 22) + 1420070400000) / 1000) - int64(zone), 0)
	fmt.Println(joinedDiscord)
	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Join dates**:", Value:
		 fmt.Sprintf("**This server**: %s\n**Discord**: %s",
		 humanize.Time(joined),
		 humanize.Time(joinedDiscord),
		 )})
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: member.User.AvatarURL("512"),
		},
		Description: fmt.Sprintf("**Display name**: %s\n**ID**: %s\n[**Avatar**](%s)\n**Game**: %s", member.User.Username, member.User.ID, member.User.AvatarURL(""), game),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL("512"),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGO!",
		},
	})
	return nil

}
