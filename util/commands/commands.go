package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"runtime/debug"
	"sort"
	"strings"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

type Command struct {
	Name      string
	ShortHelp string
	Function  CommandFunction
}

type Commands struct {
	Commands  map[string]*Command
	Redis     *redis.Client
	HelpCache string
}

type Context struct {
	T          i18n.TranslateFunc
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
	defer debug.FreeOSMemory()

	com := ctx.Commands
	if true {
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

			resp += fmt.Sprintf("<%s>\n", prefix+command.Name+strings.Repeat(" ", maxlen+1-len(command.Name))+ctx.T("command_"+command.Name+"_help"))
		}

		resp += "```\n"
		com.HelpCache = resp
	}

	session.ChannelMessageSend(message.ChannelID, com.HelpCache)

	return nil
}

func New(r *redis.Client) *Commands {
	c := &Commands{Commands: make(map[string]*Command)}
	c.RegisterCommand("help", "Provides command help.", HelpCommand)
	c.Redis = r
	return c
}

func (com *Commands) RegisterCommand(Name, ShortHelp string, Function CommandFunction) {
	c := Command{}
	c.Name = Name
	c.ShortHelp = ShortHelp
	c.Function = Function
	com.Commands[Name] = &c
}

func (com *Commands) GetCommand(msg string) (*Command, []string) {

	args := strings.Fields(msg)
	if len(args) == 0 {
		return nil, nil
	}
	if command, init := com.Commands[args[0]]; init {
		return command, args[1:]
	}

	/*
		for _, commandin := range com.Commands {

			if strings.HasPrefix(msg, commandin.Name) {
				return commandin, args[1:]
			}

		}
	*/

	return nil, nil
}

func (com *Commands) OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	//defer debug.FreeOSMemory()

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

	prefix, _ := com.Redis.Get("prefix").Result()

	if len(prefix) > 0 {

		if strings.HasPrefix(message.Content, prefix) {
			origMessage := message.Content

			message.Content = strings.TrimPrefix(message.Content, prefix)

			command, args := com.GetCommand(message.Content)
			if command != nil {

				language, err := com.Redis.Get("language_" + channel.GuildID).Result()
				if err != nil {
					com.Redis.Set("language_"+channel.GuildID, "en-GB", 0)
					language = "en-GB"
				}

				T, _ := i18n.Tfunc(language)

				ctx := &Context{
					Content:   strings.TrimPrefix(message.Content, prefix+command.Name),
					ChannelID: message.ChannelID,
					GuildID:   channel.GuildID,
					Type:      channel.Type,
					Commands:  com,
					HasPrefix: true,
					Args:      args,
					T:         T,
				}
				if len(message.Mentions) > 0 {
					ctx.HasMention = true
				}

				//start := time.Now()
				/*ret := */
				go command.Function(session, message, ctx)
				//if ret != nil {
				//	session.ChannelMessageSend(message.ChannelID, "````go\n"+ret.(*errors.Error).ErrorStack()+"\n```")
				//}
				//elapsed := time.Since(start)
				//log.Print("Command: " + command.Name + " took " + elapsed.String() + ".")
				guild, err := session.State.Guild(channel.GuildID)
				if err != nil {
					log.Printf("Can't find guild...")
					return
				}
				member, memerr := session.State.Member(channel.GuildID, message.Author.ID)
				if memerr != nil {
					log.Printf("Can't find member...")
					return
				}
				log.Printf("User %s used command \"%s\" in channel \"#%s\" (%s) and guild \"%s\" (%s)", member.User.Username, origMessage, channel.Name, channel.ID, guild.Name, channel.GuildID)
				//debug.FreeOSMemory()

				return
			}
		}
	}
	return
}

func (com *Commands) MotdEvent(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	motd, err := com.Redis.Get("motd_" + m.GuildID).Result()
	motdchannel, channelerr := com.Redis.Get("motd_" + m.GuildID + "_channel").Result()

	if err != nil || channelerr != nil {
		return
	} else {
		go s.ChannelMessageSend(motdchannel, motd)
	}
}

/*
func (com *Commands) EmojiiEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	 f, err := static.ReadFile("emotes/" + m.Content + ".png")
	 if err != nil {
		 return
	 }
	 fmt.Println(f)
	 n := bytes.NewReader(f)

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{File: &discordgo.File {Name:  m.Content + ".png", Reader: n, ContentType: "image/png"}} )

}
*/

//FS.OpenFile(CTX, "emotes/data.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
