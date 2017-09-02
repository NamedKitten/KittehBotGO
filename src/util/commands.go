package commands

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"sort"
	"log"
	"strings"
	"time"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context)

type Command struct {
	Name      string
	ShortHelp string
	Function  CommandFunction
}

type Commands struct {
	Commands []*Command
	Prefix   string
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

func HelpCommand(session *discordgo.Session, message *discordgo.MessageCreate, ctx *Context) {
	com := ctx.Commands
	prefix := com.Prefix

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

		resp += fmt.Sprintf("<%s>\n", prefix + command.Name +  strings.Repeat(" ", maxlen +1 - len(command.Name)) + command.ShortHelp)
	}

	resp += "```\n"

	session.ChannelMessageSend(message.ChannelID, resp)

	return
}

func New() *Commands {
	c := &Commands{}
	ch := Command{}
	ch.Name = "help"
	ch.ShortHelp = "Display this message."
	ch.Function = HelpCommand
	c.Commands = append(c.Commands, &ch)

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
		Commands: com,
	}

	for _, __ := range message.Mentions {
		_ = __
		ctx.HasMention = true
	}

	if len(com.Prefix) > 0 {
		if strings.HasPrefix(ctx.Content, com.Prefix) {
			ctx.HasPrefix = true
			ctx.Content = strings.TrimPrefix(ctx.Content, com.Prefix)
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
		command.Function(session, message, ctx)
		elapsed := time.Since(start)
		log.Print("Command: " + command.Name + "took " + elapsed.String() + ".")

		return
	}
}
