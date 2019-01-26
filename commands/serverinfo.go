package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	//"github.com/dustin/go-humanize"
	//"github.com/go-errors/errors"
	//"github.com/NamedKitten/KittehBotGo/util/i18n"
	//"strconv"
	//"time"
	"sort"
	"strings"
)

func init() {
	commands.RegisterCommand("serverinfo", ServerinfoCommand)
}

func ServerinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	guild, err := s.State.Guild(ctx.GuildID)
	if err != nil {
		return err
	}
	owner, err := s.State.Member(ctx.GuildID, guild.OwnerID)
	if err != nil {
		return err
	}

	verification := ""
	switch guild.VerificationLevel {
	case 0:
		verification = "None"
	case 1:
		verification = "Low"
	case 2:
		verification = "Medium"
	case 3:
		verification = "(╯°□°）╯︵ ┻━┻"
	case 4:
		verification = "┻━┻ミヽ(ಠ益ಠ)ﾉ彡┻━┻"
	}

	icon := ""
	description := fmt.Sprintf("**ID**: %s", ctx.GuildID)
	if guild.Icon != "" {
		icon = fmt.Sprintf("https://cdn.discordapp.com/icons/%s/%s.jpg", ctx.GuildID, guild.Icon)
		description = description + fmt.Sprintf("\n[Icon](%s)", icon)
	}

	bots := 0
	humans := 0

	for _, member := range guild.Members {
		if member.User.Bot {
			bots += 1
		} else {
			humans += 1
		}
	}
	var ratio string

	if bots == 0 {
		ratio = "**Bots to Humans ratio**: 1:∞"
	} else if bots < humans {
		ratio = fmt.Sprintf("**Bots to Humans ratio**: 1:%d", humans/bots)
	} else {
		ratio = fmt.Sprintf("**Humans to Bots ratio**: 1: %d", bots/humans)
	}

	fields := make([]*discordgo.MessageEmbedField, 0, 2)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Members**:", Value: fmt.Sprintf("%d", guild.MemberCount)})
	fields = append(fields,
		&discordgo.MessageEmbedField{
			Name: "**Other info**:",
			Value: fmt.Sprintf(
				"**Owner**: %s\n**Region**: %s\n**Verification level**: %s\n%s",
				owner.User.Mention(),
				guild.Region,
				verification,
				ratio,
			),
		},
	)
	guildRoles := discordgo.Roles(guild.Roles)
	sort.Sort(guildRoles)
	roles := []string{}
	for _, role := range guildRoles {
		roles = append(roles, fmt.Sprintf("<@&%s>", role.ID))
	}
	var roleList string
	if guildRoles.Len() > 0 {
		roleList = strings.Join(roles, ", ")
		if len(roleList) <= 1024 {
			fields = append(fields, &discordgo.MessageEmbedField{Name: "**Roles**:", Value: roleList})
		}
	}

	//TODO: multilingual version of gohumanize
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type:        "rich",
		Description: description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: icon,
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGo.",
		},
	})
	return nil

}
