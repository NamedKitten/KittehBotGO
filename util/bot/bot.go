package bot

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)


func Start(redis *redis.Client)  {
	Discord, _ := discordgo.New()
	commands.Setup(redis, Discord)

	Discord.State.TrackEmojis = false
	Discord.State.TrackVoice = false
	//Discord.State.TrackChannels = false
	Discord.State.MaxMessageCount = 1
	Discord.LogLevel = 1
	Discord.SyncEvents = true
	Discord.Compress = true

	Discord.AddHandler(commands.OnMessageCreate)
	commands.RegisterCommand("ping", BotCommands.PingCommand)
	commands.RegisterCommand("about", BotCommands.AboutCommand)
	commands.RegisterCommand("echo", BotCommands.EchoCommand)
	commands.RegisterCommand("userinfo", BotCommands.UserinfoCommand)
	commands.RegisterCommand("motd", BotCommands.MotdCommand)
	commands.RegisterCommand("goodbot", BotCommands.GoodBotCommand)
	commands.RegisterCommand("xkcd", BotCommands.XkcdCommand)
	commands.RegisterCommand("luaeval", BotCommands.LuaEvalCommand)
	commands.RegisterCommand("jseval", BotCommands.JSEvalCommand)
	commands.RegisterCommand("language", BotCommands.LanguageCommand)
	Discord.AddHandler(BotCommands.MotdEvent)

	fmt.Println("Getting token...")
	token, err := commands.Redis.Get("token").Result()
	if err != nil {
		fmt.Println("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}
	Discord.Token = "Bot " + token

	err = Discord.Open()
	if err != nil {
		panic(err)
	}
}