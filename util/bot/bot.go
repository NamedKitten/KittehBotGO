package bot

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	//"github.com/jonas747/discordgo"
	"github.com/go-redis/redis"
)

func Start(redis *redis.Client) {
	commands.Setup(redis)

	fmt.Println("Getting token...")
	token, err := commands.Redis.Get("token").Result()
	if err != nil {
		fmt.Println("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}
	commands.Discord.Token = "Bot " + token

	err = commands.Discord.Open()
	if err != nil {
		panic(err)
	}
}
