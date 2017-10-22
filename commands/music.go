package BotCommands

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	//"github.com/go-errors/errors"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"log"
	//"io"
	//"time"
	 //"fmt"
	 //"os"
)

func init() {
	commands.RegisterCommand("music", MusicCommand)
}

func MusicCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	log.Println("ok")
	
	g, aaerr := s.State.Guild(ctx.GuildID)
	if aaerr != nil {
		// Could not find guild.
		log.Fatal(aaerr)
	}

	channelID := ""
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			channelID = vs.ChannelID
		}
	}
	voiceConnection, _ := s.ChannelVoiceJoin(ctx.GuildID, channelID, false, true)
	voiceConnection.LogLevel =  5
	log.Println("ok1")


options := dca.StdEncodeOptions
options.RawOutput = true
options.Bitrate = 120



videoInfo, err := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
if err != nil {
   panic(err)
}

format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
downloadURL, err := videoInfo.GetDownloadURL(format)
if err != nil {
    panic(err)
}

encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
if err != nil {
    panic(err)
}
defer encodingSession.Cleanup()
    
done := make(chan error)    
dca.NewStream(encodingSession, voiceConnection, done)
derr := <- done
if derr != nil {
	log.Println(encodingSession.FFMPEGMessages())
	
    panic(derr)
}
log.Println("nya?")


//voiceConnection.Speaking(false)
//voiceConnection.Disconnect()


return nil
}

