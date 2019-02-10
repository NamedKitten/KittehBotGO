package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/dustin/go-humanize"
	"github.com/jonas747/discordgo"
	"runtime"
	"time"
)

func init() {
	commands.RegisterCommand("about", aboutCommand)
	commands.RegisterHelp("about", "Tells you about the bot.")
}

var startTime = time.Now()

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

func aboutCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fields := make([]*discordgo.MessageEmbedField, 0, 8)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Go Version**:", Value: runtime.Version(), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Discord Go Version**:", Value: discordgo.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Memory Used**:", Value: fmt.Sprintf("%s / %s (%s %s)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc), "Garbage Collected"), Inline: false})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Uptime**:", Value: getDurationString(time.Now().Sub(startTime)), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Goroutines**:", Value: fmt.Sprintf("%d", runtime.NumGoroutine()), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Servers In**:", Value: fmt.Sprintf("%d", len(commands.State.Guilds)), Inline: true})
	selfUserState := commands.State.User(false)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "About KittehBotGo",
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", selfUserState.User.ID, selfUserState.User.Avatar),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Thanks for using KittehBotGO!",
		},
	})
	return nil
}
