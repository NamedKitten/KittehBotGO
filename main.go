package main

//go:generate $GOPATH/bin/fileb0x b0x.yaml

import (
	"bufio"
	"flag"
	"fmt"
	_ "github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/util/bot"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/webdashboard"
	log "github.com/sirupsen/logrus"
	"github.com/xuyu/goredis"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var RedisClient *goredis.Redis

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
	webdashboard.UpdateInterval = *updateInterval
	var err error

	RedisClient, err = goredis.Dial(&goredis.DialConfig{
		Network:  "tcp",
		Address:  fmt.Sprintf("%s:%d", *redisIP, *redisPort),
		Password: redisPass,
		Database: *redisDB,
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
	//defer profile.Start().Stop()

	bot.Start(RedisClient)
	go webdashboard.StartDashboard()

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running..")
	var wG sync.WaitGroup
	wG.Add(1)
	var sc chan os.Signal
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		<-sc
		wG.Done()
	}()
	wG.Wait()

	//saveMemMap()

	// Cleanly close down the Discord session.
	commands.Discord.Close()
	//p.Stop()
}
