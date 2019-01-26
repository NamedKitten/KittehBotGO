package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/music"

	"github.com/bwmarrin/discordgo"
	//"github.com/go-errors/errors"
	"log"
	//"io"
	//"time"
	//"os"
	"errors"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand("music", MusicCommand)
	commands.Discord.LogLevel = 0  
}

func MusicCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	guild, err := s.State.Guild(ctx.GuildID)
	if err != nil {
		log.Println("Failed finding guild:", err)
		return nil
	}

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

		////////////////////
		// PLAYBACK CONTROL
		////////////////////
		case "play", "resume":
			// Resumes/plays
			player.EvtChan <- &music.PlayerEvtResume{}
			s.ChannelMessageSend(m.ChannelID, "Playing.")
		case "pause", "stop":
			// Pauses the playback
			player.EvtChan <- &music.PlayerEvtPause{}
			s.ChannelMessageSend(m.ChannelID, "Paused.")
		case "add":
			// Adds another element to the queue
			if len(ctx.Args) < 2 {
				err = errors.New("Nothing song specified to add.")
				break
			}

			what := ctx.Args[1]
			err = player.QueueUp(what)
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, "Song added.")
			}
		case "next", "skip":
			// Skips to the next one
			player.EvtChan <- &music.PlayerEvtNext{Index: -1}
			s.ChannelMessageSend(m.ChannelID, "Song skpped.")
		case "randnext", "r next":
			// Skips to a random item in the playlist
			player.EvtChan <- &music.PlayerEvtNext{Index: -1, Random: true}
			s.ChannelMessageSend(m.ChannelID, "Playing random song from queue.")
		case "goto", "item":
			// Skips to a specific item in the playlist
			if len(ctx.Args) < 2 {
				err = errors.New("No queue number specified.")
				break
			}

			index, err := strconv.Atoi(ctx.Args[1])
			if err != nil {
				break
			}

			if index < 0 {
				err = errors.New("No.")
				break
			}

			player.EvtChan <- &music.PlayerEvtNext{Index: index}
			s.ChannelMessageSend(m.ChannelID, "Skiping to track.")

		//////////////////
		// UTILIITIES
		//////////////////
		case "status", "stats":
			// Prints player status
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
			// Enters shuffle mode where the next item is picked randomly
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
				err = errors.New("No songs left in queue to remove.")
				break
			}
			player.EvtChan <- &music.PlayerEvtRemove{Index: index}
		}

		if err != nil {
			log.Println("Error occured:", err)
			s.ChannelMessageSend(m.ChannelID, "Error occured: " + err.Error())
		}

	}

	return nil
}
