package main

//go:generate $GOPATH/bin/fileb0x b0x.yaml

import (
	"bufio"
	"flag"
	"fmt"
	_ "github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/util/bot"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/static"
	"github.com/xuyu/goredis"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"
	//"github.com/pkg/profile"
)

var RedisClient *goredis.Redis
var UpdateInterval int
var connected int = 0

//var dg, _ = discordgo.New()

func setup() {
	reader := bufio.NewReader(os.Stdin)

	var prefix string
	fmt.Println("What prefix do you want? ")
	prefix, prefixErr := reader.ReadString('\n')

	if prefixErr != nil {
		panic(prefixErr)
	}

	var token string
	fmt.Println("What is your bots token? ")
	token, tokenErr := reader.ReadString('\n')

	if tokenErr != nil {
		panic(tokenErr)
	}

	RedisClient.Set("prefix", strings.TrimSpace(prefix), 0, 0, false, false)
	RedisClient.Set("token", strings.TrimSpace(token), 0, 0, false, false)
	fmt.Println("The bot is now setup.")
}

func init() {
	//runtime.GOMAXPROCS(runtime.NumCPU() + 4)
	//debug.SetGCPercent(1)

	updateInterval := flag.Int("updateInterval", 100, "How often the dashboard gets updated in miliseconds.")
	redisIP := flag.String("redisIP", "localhost", "IP for redis server.")
	redisPort := flag.Int("redisPort", 6379, "Port for redis server.")
	redisPassword := flag.String("redisPassword", "", "Password for redis server.")
	redisDB := flag.Int("redisDB", 0, "DB ID for redis server.")
	runSetup := flag.Bool("runSetup", false, "Run setup?")
	flag.Bool("runDashboard", true, "Run dashboard?")

	flag.Parse()
	redisPass := *redisPassword
	UpdateInterval = *updateInterval
	var err error

	RedisClient, err = goredis.Dial(&goredis.DialConfig{
		Network: "tcp",
		Address:         fmt.Sprintf("%s:%d", *redisIP, *redisPort),
		Password:     redisPass,
		Database:           *redisDB,
		Timeout:  10 * time.Second,
		MaxIdle:  10,
	})
	if err != nil {
		fmt.Println("Couldn't connect to redis...")
		panic(err)
	}

	if *runSetup {
		setup()
	}
}

func main() {
   // p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	
	bot.Start(RedisClient)
	//log.Println(bot)
	if flag.Lookup("runDashboard").Value.(flag.Getter).Get().(bool) {
		go func() {
			server, sockerr := socketio.NewServer(nil)
			if sockerr != nil {
				log.Fatal(sockerr)
			}
			server.On("connection", func(so socketio.Socket) {
				so.Join("mem")
				log.Println("on connection")
				connected = connected + 1
				so.On("disconnection", func() {
					connected = connected - 1
					log.Println("on disconnect")
				})
			})
			//so.On("mem get", func(msg string) {
			go func() {
				var lock sync.RWMutex

				for {
					lock.Lock()
					go func() {
						go debug.FreeOSMemory()

						time.Sleep(time.Millisecond * time.Duration(UpdateInterval))

						if connected > 0 {
							//debug.FreeOSMemory()
							stats := runtime.MemStats{}
							runtime.ReadMemStats(&stats)
							using := float64(stats.Alloc) / 1024 / 1024
							alloc := float64(stats.Sys) / 1024 / 1024
							cleaned := float64(stats.TotalAlloc) / 1024 / 1024
							server.BroadcastTo("mem", "mem stats", fmt.Sprintf("%g\n%g\n%g", using, alloc, cleaned))
							go debug.FreeOSMemory()
							//})
						}
						debug.FreeOSMemory()
						lock.Unlock()
					}()
				}

			}()

			//})
			server.On("error", func(so socketio.Socket, err error) {
				log.Println("error:", err)
			})

			http.Handle("/socket.io/", server)

			http.Handle("/", http.FileServer(static.HTTP))

			http.HandleFunc("/interval", func(w http.ResponseWriter, r *http.Request) {
				debug.FreeOSMemory()
				fmt.Fprintf(w, "%d.0", 100)
				debug.FreeOSMemory()
			})
			err := http.ListenAndServe("0.0.0.0:9000", nil)
			if err != nil {
				fmt.Println("Error starting http server:", err)
				os.Exit(1)
			}
		}()
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//saveMemMap()

	// Cleanly close down the Discord session.
	commands.Discord.Close()
	//p.Stop()
}
