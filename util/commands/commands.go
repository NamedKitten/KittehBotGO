package commands

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/NamedKitten/discordgo"
	"github.com/jonas747/dstate"
	"github.com/NamedKitten/KittehBotGo/util/database"
	log "github.com/sirupsen/logrus"
	"go/build"
	"runtime/debug"
	"strings"
	"time"
)

// CommandFunction is a type which commands should follow.
type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

// Commands is a map of strings to CommandFunctions which contains all the registered commands.
var Commands map[string]CommandFunction

// HelpStrings is a map of command names to the string used in the help command.
var HelpStrings map[string]string

var helpCache string

// Discord session,
var Discord *discordgo.Session

// State is a alternative state which is much better then discordgo's defaults.
var State *dstate.State

func init() {
	var err error
	Commands = make(map[string]CommandFunction)
	HelpStrings = make(map[string]string)
	Discord, err = discordgo.New()
	if err != nil {
		log.Fatal(err)
	}
	State = dstate.NewState()
	Discord.StateEnabled = false
	Discord.SyncEvents = false
	State.MaxChannelMessages = 1000
	State.MaxMessageAge = time.Hour
	State.ThrowAwayDMMessages = true
	State.TrackPrivateChannels = true
	State.CacheExpirey = time.Minute * 10
	Discord.AddHandler(State.HandleEvent)
	Discord.AddHandler(onMessageCreate)
	RegisterCommand("help", helpCommand)
	RegisterHelp("help", "Shows you all the commands this bot has.")
}

// Context given to commands and contains useful information that commands often require.
type Context struct {
	Args       []string
	Content    string
	ChannelID  int64
	GuildID    int64
	Type       discordgo.ChannelType
	HasPrefix  bool
	HasMention bool
}

func helpCommand(session *discordgo.Session, message *discordgo.MessageCreate, ctx *Context) error {
	defer debug.FreeOSMemory()

	com := Commands
	if true {
		prefix := database.Get("prefix")

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
		helpCache = resp
	}

	session.ChannelMessageSend(message.ChannelID, helpCache)

	return nil
}

// RegisterCommand registers a bot command.
func RegisterCommand(Name string, Function CommandFunction) {
	Commands[Name] = Function
}

// RegisterHelp registers a command help string.
func RegisterHelp(Name string, Help string) {
	HelpStrings[Name] = Help
}

// GetCommand returns a command, command name and arguments from a message string.
// It also splits the args to account for Unix style command line arguments.
func GetCommand(msg string) (CommandFunction, string, []string) {
	var args []string
	oldArgs := strings.Fields(msg)
	if strings.Count(msg, "=") != 0 {

		tmpStr := ""
		insideQuote := false
		for _, arg := range oldArgs {
			count := strings.Count(arg, "=")
			if insideQuote {
				tmpStr += strings.Replace(arg, "\"", "", -1) + " "
			}
			if count == 1 {
				if insideQuote {
					insideQuote = false
					args = append(args, tmpStr)
					tmpStr = ""
				} else {
					insideQuote = true
					tmpStr += strings.Replace(arg, "\"", "", -1) + " "
				}
			}
			if !insideQuote && count != 1 {
				args = append(args, arg)
			}
		}
		if insideQuote {
			args = append(args, tmpStr)
		}
	} else {
		args = oldArgs
	}

	if len(args) == 0 {
		return nil, "", nil
	}
	return Commands[args[0]], args[0], args[1:]

}

func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	var err error
	defer debug.FreeOSMemory()

	var channel *dstate.ChannelState
	channel = State.Channel(false, message.ChannelID)
	if err != nil {
		log.Printf("Can't fetch channel.")
		return
	}

	prefix := database.Get("prefix")

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
						errStr = strings.Replace(errStr, build.Default.GOPATH+"/src/", "", -1)
						errStr = strings.Replace(errStr, "/usr/lib/", "", -1)

					}

					selfUserState := State.User(false)
					Discord.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
						Type:  "rich",
						Title: "An error occurred...",
						Author: &discordgo.MessageEmbedAuthor{
							Name:    "KittehBotGo",
							IconURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%s.jpg", selfUserState.User.ID, selfUserState.User.Avatar),
						},
						Description: "```md\n" + errStr + "\n```",
					})
				}

				return
			}
		}
	}
	return
}
