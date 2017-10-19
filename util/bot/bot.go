package bot

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	//"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

func Start(redis *redis.Client) {
	commands.Setup(redis)

	commands.Discord.State.TrackEmojis = false
	commands.Discord.State.TrackVoice = false
	//commands.Discord.State.TrackChannels = false
	commands.Discord.State.MaxMessageCount = 1
	commands.Discord.LogLevel = 1
	commands.Discord.SyncEvents = true
	commands.Discord.Compress = true
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
