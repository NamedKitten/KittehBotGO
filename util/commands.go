package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-errors/errors"
	"github.com/go-redis/redis"
	"log"
	"sort"
	"strings"
	"time"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

type Command struct {
	Name      string
	ShortHelp string
	Function  CommandFunction
}

type Commands struct {
	Commands []*Command
	Redis    *redis.Client
}

type Context struct {
	Args       []string
	Content    string
	ChannelID  string
	GuildID    string
	Type       discordgo.ChannelType
	HasPrefix  bool
	HasMention bool
	Commands   *Commands
}

func HelpCommand(session *discordgo.Session, message *discordgo.MessageCreate, ctx *Context) error {
	com := ctx.Commands
	prefix, _ := com.Redis.Get("prefix").Result()

	maxlen := 0
	keys := make([]string, 0, len(com.Commands))
	cmds := make(map[string]*Command)

	for _, command := range com.Commands {
		fmt.Println(command.Name)
		nameLen := len(command.Name)
		if nameLen > maxlen {
			maxlen = nameLen
		}
		cmds[command.Name] = command
		keys = append(keys, command.Name)
	}

	sort.Strings(keys)

	header := "KittehBotGO!"
	resp := "```md\n"
	resp += header + "\n" + strings.Repeat("-", len(header)) + "\n\n"

	for _, key := range keys {
		command := cmds[key]

		resp += fmt.Sprintf("<%s>\n", prefix+command.Name+strings.Repeat(" ", maxlen+1-len(command.Name))+command.ShortHelp)
	}

	resp += "```\n"

	session.ChannelMessageSend(message.ChannelID, resp)

	return nil
}

func New(r *redis.Client) *Commands {
	c := &Commands{}
	c.RegisterCommand("help", "Provides command help.", HelpCommand)
	c.Redis = r
	return c
}

func (com *Commands) RegisterCommand(Name, ShortHelp string, Function CommandFunction) {
	c := Command{}
	c.Name = Name
	c.ShortHelp = ShortHelp
	c.Function = Function
	com.Commands = append(com.Commands, &c)
}

func (com *Commands) GetCommand(msg string) (*Command, []string) {

	args := strings.Fields(msg)
	if len(args) == 0 {
		return nil, nil
	}

	for _, commandin := range com.Commands {

		if strings.HasPrefix(msg, commandin.Name) {
			return commandin, args[1:]
		}

	}

	return nil, nil
}

func (com *Commands) OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	var err error

	if message.Author.ID == session.State.User.ID {
		return
	}

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

	ctx := &Context{
		Content:   message.Content,
		ChannelID: message.ChannelID,
		GuildID:   channel.GuildID,
		Type:      channel.Type,
		Commands:  com,
	}

	for _, __ := range message.Mentions {
		_ = __
		ctx.HasMention = true
	}

	prefix, _ := com.Redis.Get("prefix").Result()

	if len(prefix) > 0 {
		if strings.HasPrefix(ctx.Content, prefix) {
			ctx.HasPrefix = true
			ctx.Content = strings.TrimPrefix(ctx.Content, prefix)
		}
	}

	if !ctx.HasPrefix {
		return
	}

	command, args := com.GetCommand(ctx.Content)
	if command != nil {
		ctx.Content = strings.TrimPrefix(ctx.Content, command.Name)
		ctx.Args = args
		start := time.Now()
		ret := command.Function(session, message, ctx)
		if ret != nil {
			session.ChannelMessageSend(message.ChannelID, "````go\n"+ret.(*errors.Error).ErrorStack()+"\n```")
		}
		elapsed := time.Since(start)
		log.Print("Command: " + command.Name + "took " + elapsed.String() + ".")
		return
	}
}
