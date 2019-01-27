package bot

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	//"github.com/jonas747/discordgo"
	"github.com/xuyu/goredis"
)

func Start(redis *goredis.Redis) {
	commands.Setup(redis)

	fmt.Println("Getting token...")
	token, err := commands.Redis.Get("token")
	if err != nil {
		fmt.Println("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}
	commands.Discord.Token = "Bot " + string(token[:])

	err = commands.Discord.Open()
	if err != nil {
		panic(err)
	}
}
