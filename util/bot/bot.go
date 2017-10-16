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
	Redis          *redis.Client
	CommandHandler *commands.Commands
	Discord        *discordgo.Session
}

func New(redis *redis.Client) *Bot {
	CommandHandler := commands.New(redis)
	Discord, _ := discordgo.New()
 	Discord.State.TrackEmojis = false
	Discord.State.TrackVoice = false
	//Discord.State.TrackChannels = false
	Discord.State.MaxMessageCount = 1
	Discord.LogLevel = 1
	Discord.SyncEvents = true
	Discord.Compress = false

	Discord.AddHandler(CommandHandler.OnMessageCreate)
	CommandHandler.RegisterCommand("ping", "Ping!", BotCommands.PingCommand)
	CommandHandler.RegisterCommand("about", "Give info about bot.", BotCommands.AboutCommand)
	CommandHandler.RegisterCommand("echo", "Echo echo echo...", BotCommands.EchoCommand)
	CommandHandler.RegisterCommand("userinfo", "Gives info about a user.", BotCommands.UserinfoCommand)
	//CommandHandler.RegisterCommand("weeb", "Does weeb stuff.", BotCommands.WeebCommand)
	CommandHandler.RegisterCommand("motd", "Message of the day.", BotCommands.MotdCommand)
	CommandHandler.RegisterCommand("goodbot", "Call me a good bot.", BotCommands.GoodBotCommand)
	CommandHandler.RegisterCommand("xkcd", "Get a xkcd.", BotCommands.XkcdCommand)
	CommandHandler.RegisterCommand("luaeval", "Eval lua.", BotCommands.LuaEvalCommand)
	CommandHandler.RegisterCommand("jseval", "Eval js.", BotCommands.JSEvalCommand)
	CommandHandler.RegisterCommand("language", "Set language.", BotCommands.LanguageCommand)
	
	Discord.AddHandler(CommandHandler.MotdEvent)
	
	return &Bot{CommandHandler: CommandHandler, Redis: redis, Discord: Discord}
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
