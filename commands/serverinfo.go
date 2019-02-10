package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/jonas747/discordgo"
	"sort"
	"strings"
)

func init() {
	commands.RegisterCommand("serverinfo", serverinfoCommand)
	commands.RegisterHelp("serverinfo", "Shows info on this server.")
}

func serverinfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	guildState := commands.State.Guild(false, ctx.GuildID)
	guild := guildState.Guild
	ownerState := guildState.Member(false, guild.OwnerID)
	owner := ownerState.DGoCopy()

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
	description := fmt.Sprintf("**ID**: %d", ctx.GuildID)
	if guild.Icon != "" {
		icon = fmt.Sprintf("https://cdn.discordapp.com/icons/%d/%s.jpg", ctx.GuildID, guild.Icon)
		description = description + fmt.Sprintf("\n[Icon](%s)", icon)
	}

	bots := 0
	humans := 0

	for _, member := range guild.Members {
		if member.User.Bot {
			bots++
		} else {
			humans++
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
		roles = append(roles, fmt.Sprintf("<@&%d>", role.ID))
	}
	var roleList string
	if guildRoles.Len() > 0 {
		roleList = strings.Join(roles, ", ")
		if len(roleList) <= 1024 {
			fields = append(fields, &discordgo.MessageEmbedField{Name: "**Roles**:", Value: roleList})
		}
	}

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
