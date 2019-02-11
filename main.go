package main

//go:generate $GOPATH/bin/fileb0x b0x.yaml

import (
	"bufio"
	"flag"
	"fmt"
	_ "github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/util/database"
	"github.com/NamedKitten/KittehBotGo/util/bot"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/webdashboard"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"sync"
)


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

	database.Set("prefix", strings.TrimSpace(prefix))
	database.Set("token", strings.TrimSpace(token))
	fmt.Println("The bot is now setup.")
}

func init() {
	updateInterval := flag.Int("updateInterval", 100, "How often the dashboard gets updated in miliseconds.")
	runSetup := flag.Bool("runSetup", false, "Run setup?")
	flag.Bool("runDashboard", true, "Run dashboard?")
	flag.Parse()

	webdashboard.UpdateInterval = *updateInterval

	if *runSetup {
		setup()
	}
}

func main() {
	bot.Start()
	go webdashboard.StartDashboard()

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

	log.Info("Closing Discord Connection...")

	// Cleanly close down the Discord session.
	commands.Discord.Close()
}
