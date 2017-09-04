package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/config"
	"github.com/NamedKitten/KittehBotGo/util"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"runtime"
	"time"
)

var startTime = time.Now()

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

func AboutCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fields := make([]*discordgo.MessageEmbedField, 0, 6)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**KittehBotGo Version**:", Value: config.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Go Version**:", Value: runtime.Version(), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**DiscordGo Version**:", Value: discordgo.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Memory used**:", Value: fmt.Sprintf("%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc)), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Uptime**:", Value: getDurationString(time.Now().Sub(startTime)), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Goroutines**:", Value: fmt.Sprintf("%d", runtime.NumGoroutine()), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Servers**:", Value: fmt.Sprintf("%d", len(s.State.Guilds)), Inline: true})

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
	return nil
}
