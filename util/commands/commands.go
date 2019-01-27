package commands

import (
	"fmt"
	"github.com/jonas747/discordgo"
	"github.com/xuyu/goredis"
	"go/build"
	"github.com/go-errors/errors"
	"github.com/jonas747/dstate"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

var Redis *goredis.Redis
var Commands map[string]CommandFunction
var HelpStrings map[string]string

var HelpCache string
var Discord *discordgo.Session

var State *dstate.State

func init() {
	var err error
	Commands = make(map[string]CommandFunction)
	HelpStrings = make(map[string]string)
	Discord, err= discordgo.New()
	if (err != nil) {
		log.Fatal(err)
	}
	State = dstate.NewState()
	Discord.StateEnabled = false
	Discord.SyncEvents = true
	State.MaxChannelMessages = 1000
	State.MaxMessageAge = time.Hour
	State.ThrowAwayDMMessages = true
	State.TrackPrivateChannels = true
	State.CacheExpirey = time.Minute * 10
	Discord.AddHandler(State.HandleEvent)
	Discord.AddHandler(OnMessageCreate)
	RegisterCommand("help", HelpCommand)
	RegisterHelp("help", "Shows you all the commands this bot has.")
}

type Context struct {
	Args       []string
	Content    string
	ChannelID  int64
	GuildID    int64
	Type       discordgo.ChannelType
	HasPrefix  bool
	HasMention bool
}

func HelpCommand(session *discordgo.Session, message *discordgo.MessageCreate, ctx *Context) error {
	defer debug.FreeOSMemory()

	com := Commands
	if true {
		pre, _ := Redis.Get("prefix")
		prefix := string(pre[:])
		

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

func Setup(r *goredis.Redis) {             
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

	var channel *dstate.ChannelState
	channel = State.Channel(false, message.ChannelID)
	if err != nil {
		log.Printf("Can't fetch channel.")
		return
	}

	pre, _ := Redis.Get("prefix")
	prefix := string(pre[:])
	
	if len(prefix) > 0 {

		if strings.HasPrefix(message.Content, prefix) {

			message.Content = strings.TrimPrefix(message.Content, prefix)

			command, name, args := GetCommand(message.Content)
			if command != nil {

				ctx := &Context{
					Content:   strings.TrimPrefix(message.Content, prefix+name),
					ChannelID: message.ChannelID,
					GuildID:   channel.Guild.ID,
					Type:      channel.Type,
					HasPrefix: true,
					Args:      args,
				}
				if len(message.Mentions) > 0 {
					ctx.HasMention = true
				}
				err = command(session, message, ctx)

				if err != nil {
					betterErr := errors.Wrap(err, 1)
					errStr := ""
					for _, frame := range betterErr.StackFrames() {
						line := fmt.Sprintf("[%s](%d)\n", frame.File, frame.LineNumber)
						source, err := frame.SourceLine()
						if err != nil {
							errStr += line
						} else {
							errStr += line + fmt.Sprintf(" %s: %s\n", frame.Name, source)
						}
						errStr = strings.Replace(errStr, build.Default.GOPATH + "/src/", "", -1)
						errStr = strings.Replace(errStr, "/usr/lib/", "", -1)

					}

					selfUserState := State.User(false)
					Discord.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
						Type: "rich",
						Title: "An error occured...",
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "KittehBotGo",
							IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", selfUserState.User.ID, selfUserState.User.Avatar),
						},
						Description:  "```md\n" + errStr + "\n```",
					})
				}

				return
			}
		}
	}
	return
}