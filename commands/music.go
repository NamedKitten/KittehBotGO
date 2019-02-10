package BotCommands

import (
	"errors"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/music"
	"github.com/jonas747/discordgo"
	"log"
	"strconv"
)

func init() {
	commands.RegisterCommand("music", musicCommand)
	commands.RegisterHelp("music", "Music commands.")
	commands.Discord.LogLevel = 0
}

func musicCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	var err error
	guild := commands.State.Guild(false, ctx.GuildID).Guild

	var vs *discordgo.VoiceState

	if len(ctx.Args) > 0 {

		for _, v := range guild.VoiceStates {
			if v.UserID == m.Author.ID {
				vs = v
				break
			}
		}

		if ctx.Args[0] == "join" {
			_, err = music.CreatePlayer(guild.ID, vs.ChannelID)
			if err != nil {
				log.Println("Error creating player:", err)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Joining voice channel...")
			}
			return nil
		}

		player := music.GetPlayer(guild.ID)
		if player == nil {
			return nil
		}

		switch ctx.Args[0] {
		case "die", "kill", "leave":
			player.EvtChan <- &music.PlayerEvtKill{}
			s.ChannelMessageSend(m.ChannelID, "Leaving voice channel...")

		case "play", "resume":
			player.EvtChan <- &music.PlayerEvtResume{}
			s.ChannelMessageSend(m.ChannelID, "Playing.")
		case "pause", "stop":
			player.EvtChan <- &music.PlayerEvtPause{}
			s.ChannelMessageSend(m.ChannelID, "Paused.")
		case "add":
			if len(ctx.Args) < 2 {
				err = errors.New("no song specified")
				break
			}

			what := ctx.Args[1]
			err = player.QueueUp(what)
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, "Song added.")
			}
		case "next", "skip":
			player.EvtChan <- &music.PlayerEvtNext{Index: -1}
			s.ChannelMessageSend(m.ChannelID, "Song skpped.")
		case "randnext", "r next":
			player.EvtChan <- &music.PlayerEvtNext{Index: -1, Random: true}
			s.ChannelMessageSend(m.ChannelID, "Playing random song from queue.")
		case "goto", "item":
			if len(ctx.Args) < 2 {
				err = errors.New("queue number not specified")
				break
			}

			index, err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				break
			}

			if index < 0 {
				err = errors.New("no item to skip to")
				break
			}

			player.EvtChan <- &music.PlayerEvtNext{Index: index}
			s.ChannelMessageSend(m.ChannelID, "Skiping to track.")

		case "status", "stats":
			status := player.Status()

			itemDuration := "None"
			itemName := "No song currently playing."
			if status.Current != nil {
				itemDuration = status.Current.Duration.String()
				itemName = status.Current.Title
			}

			out := fmt.Sprintf("**Player status:**\n**Paused:** %v\n**Title:** %s\n**Position:** %s/%s\n**Shuffle:** %v\n", status.Paused, itemName, status.Position.String(), itemDuration, status.Shuffle)

			if len(status.Queue) > 0 {
				out += "\n\n**Queue:**\n"
			}

			for k, v := range status.Queue {
				out += fmt.Sprintf("**#%d:** %s - %s (<%s>)\n", k, v.Title, v.Duration.String(), "https://www.youtube.com/watch?v="+v.ID)
			}

			s.ChannelMessageSend(m.ChannelID, out)

		case "shuffle":
			shuffle := player.ToggleShuffle()
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Shuffle: %v", shuffle))
		case "remove":
			// Removes an element in the playlist
			index := 0
			index, err = strconv.Atoi(ctx.Args[1])
			if err != nil {
				break
			}

			if index < 0 {
				err = errors.New("no songs in queue to remove")
				break
			}
			player.EvtChan <- &music.PlayerEvtRemove{Index: index}
		}

		if err != nil {
			log.Println("Error occured:", err)
			s.ChannelMessageSend(m.ChannelID, "Error occured: "+err.Error())
		}

	} else {
		helpMsg := "```md\n"
		helpMsg += "# Music Command Help\n"
		helpMsg += "[music join](Joins the voice channel you are currently in)\n"
		helpMsg += "[music leave](Leaves the voice channel the bot is in)\n"
		helpMsg += "[music play/resume](Plays or Resumes playback of the currently playing track)\n"
		helpMsg += "[music pause/stop](Pauses playback of the currently playing track)\n"
		helpMsg += "[music status/stats](Shows stats for the current session)\n"
		helpMsg += "[music next/skip](Skips a track)\n"
		helpMsg += "[music rnext](Skips to a random item in the playlist)\n"
		helpMsg += "[music goto](Skips to a specified item by position in the playlist)\n"
		helpMsg += "[music shuffle](Enter shuffle mode)\n"
		helpMsg += "[music remove](Removes a item by position in playlist.)\n"
		helpMsg += "```"
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}

	return nil
}
