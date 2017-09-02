package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

type Context struct {
	Args       []string
	Content    string
	ChannelID  string
	GuildID    string
	Type       discordgo.ChannelType
	HasPrefix  bool
	HasMention bool
}

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

func New() *Commands {
	c := &Commands{}
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
	log.Println(args)
	log.Print(msg)
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
	}

	for _, __ := range message.Mentions {
		_ = __
		ctx.HasMention = true
	}

	if len(com.Prefix) > 0 {
		if strings.HasPrefix(ctx.Content, com.Prefix) {
			log.Print(ctx.Content)
			ctx.HasPrefix = true
			ctx.Content = strings.TrimPrefix(ctx.Content, com.Prefix)
		}
	}

	if !ctx.HasPrefix {
		return
	}

	command, args := com.GetCommand(ctx.Content)
	if command != nil {
		log.Print(ctx.Content)
		ctx.Content = strings.TrimPrefix(ctx.Content, command.Name)
		ctx.Args = args
		command.Function(session, message, ctx)
		return
	}
}
