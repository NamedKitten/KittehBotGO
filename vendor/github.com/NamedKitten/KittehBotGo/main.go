package main

import (
	"github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/config"
	"github.com/NamedKitten/KittehBotGo/util"
	"bufio"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	//"github.com/wcharczuk/go-chart"
	//"io/ioutil"
)

var redisclient *redis.Client
var dg, _ = discordgo.New()
/*
var memAllocMapX []float64
var memAllocMapY []float64
var memSysMapX []float64
var memSysMapY []float64
var memTotalAllocMapX []float64
var memTotalAllocMapY []float64
*/

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
/*
func saveMemMap() {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{Show: true},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "MemAlloc",
				XValues: memAllocMapX,
				YValues: memAllocMapY,
			},
			chart.ContinuousSeries{
				Name:    "MemSys",
				XValues: memSysMapX,
				YValues: memSysMapY,
			},
			chart.ContinuousSeries{
				Name:    "MemTotalAlloc",
				XValues: memTotalAllocMapX,
				YValues: memTotalAllocMapY,
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	_ := graph.Render(chart.PNG, buffer)
	ioutil.WriteFile("out.png", buffer, 0644)
}
*/

func init() {
	redisIP := flag.String("redisIP", "localhost", "IP for redis server.")
	redisPort := flag.Int("redisPort", 6379, "Port for redis server.")
	redisPassword := flag.String("redisPassword", "", "Password for redis server.")
	redisDB := flag.Int("redisDB", 0, "DB ID for redis server.")
	version := flag.Bool("version", false, "Print version and exit.")
	runSetup := flag.Bool("runSetup", false, "Run setup?")

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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//saveMemMap()

	// Cleanly close down the Discord session.
	dg.Close()
 }
