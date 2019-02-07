package music

import (
	"errors"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/dca"
	"github.com/jonas747/discordgo"
	"github.com/rylio/ytdl"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	Players     = make(map[int64]*Player)
	playersLock sync.Mutex

	encodeOptions *dca.EncodeOptions
)

func init() {
	encodeOptions = &dca.EncodeOptions{}
	*encodeOptions = *dca.StdEncodeOptions

	encodeOptions.Bitrate = 128
	encodeOptions.RawOutput = true
}

func GetPlayer(guildID int64) *Player {
	playersLock.Lock()
	p := Players[guildID]
	playersLock.Unlock()
	return p
}

type Player struct {
	sync.Mutex
	GuildID int64
	vc      *discordgo.VoiceConnection

	queue          []*ytdl.VideoInfo
	shuffle        bool
	nextQueueIndex int

	currentEncodeSession *dca.EncodeSession
	currentStream        *dca.StreamingSession
	CurrentlyPlaying     *ytdl.VideoInfo
	downloadWriter       io.Closer

	running bool
	EvtChan chan interface{}
}

func CreatePlayer(guildID int64, channelID int64) (*Player, error) {
	playersLock.Lock()
	if p, ok := Players[guildID]; ok {
		playersLock.Unlock()
		return p, nil
	}
	defer playersLock.Unlock()

	player := &Player{
		GuildID: guildID,
		EvtChan: make(chan interface{}),
	}

	vc, err := commands.Discord.GatewayManager.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	player.vc = vc
	player.running = true
	go player.run()

	Players[guildID] = player
	return player, nil
}

func (p *Player) run() {

	log.Println("Waiting for voice connected...")
	//<-p.vc.Connected
	log.Println("Voice connected!")

	defer func() {
		p.Lock()
		p.running = false
		p.Unlock()
	}()

	ticker := time.NewTicker(time.Second)
	for {

		select {
		case <-ticker.C:
			// check for stream status
			p.Lock()
			if p.currentStream == nil {
				p.Unlock()
				continue
			}

			fin, err := p.currentStream.Finished()
			if !fin {
				p.Unlock()
				continue
			}
			// Finished huh? amazing
			if err != nil {
				if p.vc != nil {
					p.vc.Disconnect()
				}

				log.Println("An error occured:", err)
				p.Unlock()
				return
			}
			p.playNext()
			p.Unlock()
		case evt := <-p.EvtChan:
			p.Lock()
			p.handleEvent(evt)
			p.Unlock()
		}
	}
}

func (p *Player) handleEvent(evt interface{}) {
	switch t := evt.(type) {
	case *PlayerEvtPause:
		if p.currentStream != nil {
			p.currentStream.SetPaused(true)
		}
	case *PlayerEvtKill: // >:O
		if p.currentStream != nil {
			p.currentStream.SetPaused(true)
			if p.downloadWriter != nil {
				p.downloadWriter.Close()
			}
		}

		if p.currentEncodeSession != nil {
			err := p.currentEncodeSession.Stop()
			if err != nil {
				log.Println("Error stopping player:", err)
			}
			// Clean up
			p.currentEncodeSession.Truncate()
		}

		if p.vc != nil {
			err := p.vc.Disconnect()
			if err != nil {
				log.Println("Error disconnecting:", err)
			}
		}

		playersLock.Lock()
		delete(Players, p.GuildID)
		playersLock.Unlock()
	case *PlayerEvtNext:

		if t.Index != -1 {
			p.nextQueueIndex = t.Index
		} else if t.Random && len(p.queue) > 0 {
			p.nextQueueIndex = rand.Intn(len(p.queue))
		}

		p.playNext()
	case *PlayerEvtResume:
		if p.currentStream != nil {
			p.currentStream.SetPaused(false)
		}
	case *PlayerEvtQeue:
		if len(p.queue) > 50 {
			break
		}

		p.queue = append(p.queue, t.Info)
		if p.currentEncodeSession == nil {
			p.playNext()
		}
	case *PlayerEvtRemove:
		if t.Index >= len(p.queue) {
			log.Println("Attempted to remove element out of bounds")
			break
		}
		p.queue = append(p.queue[:t.Index], p.queue[t.Index+1:]...)
	default:
		log.Println("UNKNOWN PLAYER EVENT", evt)
	}
}

func (p *Player) playNext() {
	if p.currentEncodeSession != nil {
		if p.downloadWriter != nil {
			p.downloadWriter.Close() // Stop the download!
		}
		if p.currentStream != nil {
			p.currentStream.SetPaused(false)
		}
		// Clean up

		p.currentEncodeSession.Truncate()
	}

	if len(p.queue) < 1 {
		p.currentStream = nil
		p.currentEncodeSession = nil
		p.CurrentlyPlaying = nil
		p.downloadWriter = nil
		return
	}

	nextIndex := p.nextQueueIndex
	if nextIndex >= len(p.queue) {
		nextIndex = 0
	}

	next := p.queue[nextIndex]

	p.queue = append(p.queue[:nextIndex], p.queue[nextIndex+1:]...)

	// Update the queue index
	if p.shuffle {
		if len(p.queue) > 0 {
			p.nextQueueIndex = rand.Intn(len(p.queue))
		}
	} else {
		p.nextQueueIndex++
	}

	log.Println("Playing", next.Title)

	reader, writer := io.Pipe()

	go func() {
		bestFormats := next.Formats.Best(ytdl.FormatAudioBitrateKey)
		if len(bestFormats) < 1 {
			log.Println("No format available")
			return
		}

		next.Download(next.Formats.Best(ytdl.FormatAudioEncodingKey)[0], writer)
		writer.Close()
	}()
	p.downloadWriter = writer

	encodeSession, _ := dca.EncodeMem(reader, encodeOptions)
	stream := dca.NewStream(encodeSession, p.vc, nil)

	p.currentStream = stream
	p.currentEncodeSession = encodeSession
	p.CurrentlyPlaying = next
	fmt.Println(encodeSession.FFMPEGMessages())

}

func (p *Player) QueueUp(url string) error {
	if !p.running {
		return errors.New("Player is not running.  ")
	}

	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return err
	}

	bestFormats := info.Formats.Best(ytdl.FormatAudioBitrateKey)
	if len(bestFormats) < 1 {
		return errors.New("No formats availabls")
	}

	log.Println("Sending", info.Title, "to the player...")
	p.EvtChan <- &PlayerEvtQeue{info}
	return nil
}

func (p *Player) ToggleShuffle() bool {
	p.Lock()
	defer p.Unlock()
	p.shuffle = !p.shuffle

	if p.shuffle && len(p.queue) > 0 {
		p.nextQueueIndex = rand.Intn(len(p.queue))
	}

	return p.shuffle
}

type PlayerStatus struct {
	Paused   bool
	Position time.Duration
	Current  *ytdl.VideoInfo
	Queue    []*ytdl.VideoInfo
	Shuffle  bool
}

// Return all the elemns in the queue
func (p *Player) Status() *PlayerStatus {
	p.Lock()
	paused := true
	position := time.Duration(0)

	if p.currentStream != nil {
		paused = p.currentStream.Paused()
		position = p.currentStream.PlaybackPosition()
	}

	status := &PlayerStatus{
		Paused:   paused,
		Position: position,
		Current:  p.CurrentlyPlaying,
		Queue:    p.queue,
		Shuffle:  p.shuffle,
	}
	p.Unlock()

	return status
}

type PlayerEvtQeue struct {
	Info *ytdl.VideoInfo
}

type PlayerEvtPause struct{}
type PlayerEvtResume struct{}
type PlayerEvtKill struct{}
type PlayerEvtNext struct {
	Index  int  // Jumps to specified index if > -1
	Random bool // Jumps to a random element
}
type PlayerEvtRemove struct {
	Index int // What to remove
}
