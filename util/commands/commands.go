package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"log"
	"runtime/debug"
	"strings"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

var Redis *redis.Client
var Commands map[string]CommandFunction
var HelpStrings map[string]string

var HelpCache string
var Discord *discordgo.Session

func init() {
	Commands = make(map[string]CommandFunction)
	HelpStrings = make(map[string]string)
	Discord, _ = discordgo.New()
	Discord.AddHandler(OnMessageCreate)
	RegisterCommand("help", HelpCommand)
	RegisterHelp("help", "Shows you all the commands this bot has.")

}

type Context struct {
	Args       []string
	Content    string
	ChannelID  string
	GuildID    string
	Type       discordgo.ChannelType
	HasPrefix  bool
	HasMention bool
}

func HelpCommand(session *discordgo.Session, message *discordgo.MessageCreate, ctx *Context) error {
	defer debug.FreeOSMemory()

	com := Commands
	if true {
		prefix, _ := Redis.Get("prefix").Result()

		maxlen := 0

		for name := range com {
			if len(name) > maxlen {
				maxlen = len(name)
			}
		}

		header := "KittehBotGO!"
		resp := "```md\n"
		resp += header + "\n" + strings.Repeat("-", len(header)) + "\n\n"

		for name := range com {
			resp += fmt.Sprintf("<%s>\n", prefix+name+strings.Repeat(" ", maxlen+2-len(name))+HelpStrings[name])
		}

		resp += "```\n"
		HelpCache = resp
	}

	session.ChannelMessageSend(message.ChannelID, HelpCache)

	return nil
}

func Setup(r *redis.Client) {             
	Redis = r
}

func RegisterCommand(Name string, Function CommandFunction) {
	Commands[Name] = Function
}

func RegisterHelp(Name string, Help string) {
	HelpStrings[Name] = Help
}

func GetCommand(msg string) (CommandFunction, string, []string) {

	args := strings.Fields(msg)
	if len(args) == 0 {
		return nil, "", nil
	}
	return Commands[args[0]], args[0], args[1:]

}

func OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	var err error

	var channel *discordgo.Channel
	channel, err = session.State.Channel(message.ChannelID)
	if err != nil {
		channel, err = session.Channel(message.ChannelID)
		if err != nil {
			log.Printf("Can't fetch channel.")
			return
		}
		err = session.State.ChannelAdd(channel)
		if err != nil {
			log.Printf("Can't add channel to state.")
		}
	}

	prefix, _ := Redis.Get("prefix").Result()

	if len(prefix) > 0 {

		if strings.HasPrefix(message.Content, prefix) {

			message.Content = strings.TrimPrefix(message.Content, prefix)

			command, name, args := GetCommand(message.Content)
			if command != nil {

				ctx := &Context{
					Content:   strings.TrimPrefix(message.Content, prefix+name),
					ChannelID: message.ChannelID,
					GuildID:   channel.GuildID,
					Type:      channel.Type,
					HasPrefix: true,
					Args:      args,
				}
				if len(message.Mentions) > 0 {
					ctx.HasMention = true
				}
				command(session, message, ctx)
				return
			}
		}
	}
	return
}