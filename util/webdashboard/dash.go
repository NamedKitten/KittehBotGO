package webdashboard

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/music"
	"github.com/NamedKitten/KittehBotGo/util/static"
	"github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"
)

// UpdateInterval is how often in milliseconds the web dashboard is updated.
var UpdateInterval int

// connected is the amount of clients connected to the websocket.
var connected int
var server *socketio.Server

type guildInfo struct {
	Icon    string `json:"icon"`
	Name    string `json:"name"`
	Members int    `json:"members"`
}

type memoryStatsInfo struct {
	Using     float64 `json:"using"`
	Allocated float64 `json:"allocated"`
	Cleaned   float64 `json:"cleaned"`
}

type musicPlayerInfo struct {
	GuildName string `json:"guildName"`
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title"`
}

func musicPlayerInfoUpdater() {

	for {
		if connected > 0 {
			var musicPlayerInfoList []musicPlayerInfo
			for gid, player := range music.Players {
				if player.CurrentlyPlaying != nil {
					guildName := commands.State.Guild(false, gid).Guild.Name
					playerStatus := player.Status()
					thumbnail := player.CurrentlyPlaying.GetThumbnailURL("best").String()
					title := playerStatus.Current.Title
					musicPlayerInfoList = append(musicPlayerInfoList, musicPlayerInfo{guildName, thumbnail, title})
				}
			}
			jsonMusicPlayersInfo, _ := json.Marshal(musicPlayerInfoList)
			server.BroadcastTo("music", "music stats", string(jsonMusicPlayersInfo))
		}
		time.Sleep(time.Second)
	}
}

func guildsListUpdater() {

	for {
		if connected > 0 {
			var guildsInfoList []guildInfo
			for _, guild := range commands.State.Guilds {
				guildIcon := fmt.Sprintf("https://cdn.discordapp.com/icons/%d/%s.jpg", guild.ID, guild.Guild.Icon)
				guildsInfoList = append(guildsInfoList, guildInfo{guildIcon, guild.Guild.Name, guild.Guild.MemberCount})
			}
			jsonGuildInfo, _ := json.Marshal(guildsInfoList)
			server.BroadcastTo("guilds", "guilds stats", string(jsonGuildInfo))
		}
		time.Sleep(time.Second * 2)
	}
}

func memStatsUpdater() {

	for {
		debug.FreeOSMemory()
		if connected > 0 {
			stats := runtime.MemStats{}
			runtime.ReadMemStats(&stats)
			memStatsInfo := memoryStatsInfo{
				float64(stats.Alloc) / 1024 / 1024,
				float64(stats.Sys) / 1024 / 1024,
				float64(stats.TotalAlloc) / 1024 / 1024,
			}
			jsonStatsInfo, _ := json.Marshal(memStatsInfo)

			server.BroadcastTo("mem", "mem stats", string(jsonStatsInfo))
		}
		debug.FreeOSMemory()
		time.Sleep(time.Millisecond * time.Duration(UpdateInterval))
	}
}

// StartDashboard starts the web dashboard.
func StartDashboard() {
	if flag.Lookup("runDashboard").Value.(flag.Getter).Get().(bool) {
		socketServer, sockerr := socketio.NewServer(nil)
		if sockerr != nil {
			log.WithFields(log.Fields{
				"err": sockerr,
			}).Error("A error occurred creating the socket server.")
		}
		server = socketServer
		server.On("connection", func(so socketio.Socket) {
			so.Join("mem")
			so.Join("music")
			so.Join("guilds")
			connected = connected + 1
			so.On("disconnection", func() {
				connected = connected - 1
			})
		})
		go memStatsUpdater()
		go guildsListUpdater()
		go musicPlayerInfoUpdater()

		server.On("error", func(so socketio.Socket, err error) {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("A error occurred in the socket.")
		})

		http.Handle("/socket.io/", server)

		http.Handle("/", http.FileServer(static.HTTP))

		err := http.ListenAndServe("0.0.0.0:9000", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to start http server.")
		}
	}
}
