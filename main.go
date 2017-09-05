package main

//go:generate $GOPATH/bin/go-bindata -pkg main -o static.go static/...

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/config"
	"github.com/NamedKitten/KittehBotGo/util"
	"github.com/bwmarrin/discordgo"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var redisclient *redis.Client
var dg, _ = discordgo.New()

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

	redisclient.Set("prefix", prefix, 0)
	redisclient.Set("token", strings.TrimSpace(token), 0)
	fmt.Println("The bot is now setup.")
}

func init() {
	redisIP := flag.String("redisIP", "localhost", "IP for redis server.")
	redisPort := flag.Int("redisPort", 6379, "Port for redis server.")
	redisPassword := flag.String("redisPassword", "", "Password for redis server.")
	redisDB := flag.Int("redisDB", 0, "DB ID for redis server.")
	version := flag.Bool("version", false, "Print version and exit.")
	runSetup := flag.Bool("runSetup", false, "Run setup?")
	flag.Bool("runDashboard", true, "Run dashboard?")

	flag.Parse()

	if *version {
		fmt.Println(config.VERSION)
	}

	redisclient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", *redisIP, *redisPort),
		Password:     *redisPassword,
		DB:           *redisDB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	pong, err := redisclient.Ping().Result()
	if err != nil || pong != "PONG" {
		fmt.Println("Couldn't connect to redis...")
		panic(err)
	}

	if *runSetup {
		setup()
	}
}

func main() {
	var CommandHandler = commands.New(redisclient)
	fmt.Println("Getting token...")
	token, err := redisclient.Get("token").Result()
	if err != nil {
		fmt.Println("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}

	dg.Token = "Bot " + token
	dg.LogLevel = 1
	dg.SyncEvents = false

	dg.AddHandler(CommandHandler.OnMessageCreate)
	CommandHandler.RegisterCommand("ping", "Ping!", BotCommands.PingCommand)
	CommandHandler.RegisterCommand("about", "Give info about bot.", BotCommands.AboutCommand)
	CommandHandler.RegisterCommand("echo", "Echo echo echo...", BotCommands.EchoCommand)
	CommandHandler.RegisterCommand("userinfo", "Gives info about a user.", BotCommands.UserinfoCommand)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	if flag.Lookup("runDashboard").Value.(flag.Getter).Get().(bool) {
		go func() {
			http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "static"}))
			http.HandleFunc("/getdata", func(w http.ResponseWriter, r *http.Request) {

				stats := runtime.MemStats{}
				runtime.ReadMemStats(&stats)
				using := float64(stats.Alloc) / 1024 / 1024
				alloc := float64(stats.Sys) / 1024 / 1024
				cleaned := float64(stats.TotalAlloc) / 1024 / 1024

				fmt.Fprintf(w, "%g\n%g\n%g", using, alloc, cleaned)

			})
			err := http.ListenAndServe("127.0.0.1:9000", nil)
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
	dg.Close()
}
