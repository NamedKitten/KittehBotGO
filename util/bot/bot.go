package bot

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	log "github.com/sirupsen/logrus"
	"github.com/xuyu/goredis"
)

func Start(redis *goredis.Redis) {
	commands.Setup(redis)

	token, err := commands.Redis.Get("token")
	if err != nil {
		log.Fatal("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}
	commands.Discord.Token = "Bot " + string(token[:])

	err = commands.Discord.Open()
	if err != nil {
		panic(err)
	}
}
