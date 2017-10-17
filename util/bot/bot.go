package bot

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/commands"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
)

type EventFunc func(*Bot, *discordgo.Session, interface{})

type Bot struct {
	Redis   *redis.Client
	Discord *discordgo.Session
}

func New(redis *redis.Client) *Bot {
	Discord, _ := discordgo.New()
	commands.Setup(redis, Discord)

	Discord.State.TrackEmojis = false
	Discord.State.TrackVoice = false
	//Discord.State.TrackChannels = false
	Discord.State.MaxMessageCount = 1
	Discord.LogLevel = 1
	Discord.SyncEvents = true
	Discord.Compress = false

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

	return &Bot{Redis: redis, Discord: Discord}
}

func (bot *Bot) Start() {
	fmt.Println("Getting token...")
	token, err := bot.Redis.Get("token").Result()
	if err != nil {
		fmt.Println("Token not found, please run with -runSetup to enter setup.")
		panic(err)
	}
	bot.Discord.Token = "Bot " + token

	err = bot.Discord.Open()
	if err != nil {
		panic(err)
	}

}
