package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/config"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"runtime"
	"runtime/debug"
	"time"
)

func init() {
	commands.RegisterCommand("about", AboutCommand)
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

func AboutCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	fields := make([]*discordgo.MessageEmbedField, 0, 7)
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_kversion") + "**:", Value: config.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_gversion") + "**:", Value: runtime.Version(), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_dversion") + "**:", Value: discordgo.VERSION, Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_memused") + "**:", Value: fmt.Sprintf("%s / %s (%s %s)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc), ctx.T("command_about_garbage")), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_uptime") + "**:", Value: getDurationString(time.Now().Sub(startTime)), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**Goroutines**:", Value: fmt.Sprintf("%d", runtime.NumGoroutine()), Inline: true})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "**" + ctx.T("command_about_servers") + "**:", Value: fmt.Sprintf("%d", len(s.State.Guilds)), Inline: true})

	for _, server := range s.State.Guilds {
		fmt.Println(server.Name)
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Type: "rich",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    ctx.T("command_about_about", struct{ Person string }{Person: s.State.User.Username}),
			IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", s.State.User.ID, s.State.User.Avatar),
		},
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: ctx.T("command_about_thanks"), //"Thanks for using KittehBotGO!",
		},
	})
	return nil
}
