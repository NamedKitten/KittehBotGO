// +build !nodashboard

package webdashboard

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/music"
	"github.com/NamedKitten/KittehBotGo/util/static"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"strings"
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
var server *gosocketio.Server

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
					thumbnail := strings.Replace(player.CurrentlyPlaying.GetThumbnailURL("default").String(), "%20", "", -1)
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
		socketServer := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

		server = socketServer
		server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
			c.Join("mem")
			c.Join("music")
			c.Join("guilds")
			connected = connected + 1
		})
		server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
			connected = connected - 1
		})
	
		go memStatsUpdater()
		go guildsListUpdater()
		go musicPlayerInfoUpdater()

		serveMux := http.NewServeMux()

		serveMux.Handle("/socket.io/", server)

		serveMux.Handle("/", http.FileServer(static.HTTP))

		err := http.ListenAndServe("0.0.0.0:9000", serveMux)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to start http server.")
		}
	}
}
