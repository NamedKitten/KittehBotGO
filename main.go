package main

//go:generate $GOPATH/bin/fileb0x b0x.yaml

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/config"
	"github.com/NamedKitten/KittehBotGo/util/bot"
	"github.com/NamedKitten/KittehBotGo/util/internaldb"
	"github.com/NamedKitten/KittehBotGo/util/static"
	//"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-redis/redis"
	//_ "net/http/pprof"
	//_ "golang.org/x/mobile/app"
	_ "github.com/NamedKitten/KittehBotGo/util/i18n"
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
)

var RedisClient *redis.Client
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

	RedisClient.Set("prefix", strings.TrimSpace(prefix), 0)
	RedisClient.Set("token", strings.TrimSpace(token), 0)
	fmt.Println("The bot is now setup.")
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 4)
	//debug.SetGCPercent(1)

	updateInterval := flag.Int("updateInterval", 100, "How often the dashboard gets updated in miliseconds.")
	redisIP := flag.String("redisIP", "localhost", "IP for redis server.")
	redisPort := flag.Int("redisPort", 6379, "Port for redis server.")
	redisPassword := flag.String("redisPassword", "", "Password for redis server.")
	redisDB := flag.Int("redisDB", 0, "DB ID for redis server.")
	version := flag.Bool("version", false, "Print version and exit.")
	runSetup := flag.Bool("runSetup", false, "Run setup?")
	useInternalDB := flag.Bool("useInternalDB", false, "Use built in redis?")
	internalDBFile := flag.String("internalDBFile", "", "File to save data to for internal redis server   .")

	flag.Bool("runDashboard", true, "Run dashboard?")

	flag.Parse()
	UpdateInterval = *updateInterval
	if *version {
		fmt.Println(config.VERSION)
	}

	if *useInternalDB {
		go database.Start(*internalDBFile, *redisPort)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", *redisIP, *redisPort),
		Password:     *redisPassword,
		DB:           *redisDB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	pong, err := RedisClient.Ping().Result()
	if err != nil || pong != "PONG" {
		fmt.Println("Couldn't connect to redis...")
		panic(err)
	}

	if *runSetup {
		setup()
	}
}

func main() {
	bot := bot.New(RedisClient)
	bot.Start()
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
							go server.BroadcastTo("mem", "mem stats", fmt.Sprintf("%g\n%g\n%g", using, alloc, cleaned))
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
	bot.Discord.Close()
}
