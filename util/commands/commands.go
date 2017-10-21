package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"runtime/debug"
	"strings"
)

type CommandFunction func(*discordgo.Session, *discordgo.MessageCreate, *Context) error

var Redis *redis.Client
var Commands map[string]CommandFunction
var HelpCache string
var Discord *discordgo.Session

func init() {
	Commands = make(map[string]CommandFunction)
	Discord, _ = discordgo.New()
	Discord.AddHandler(OnMessageCreate)
	RegisterCommand("help", HelpCommand)
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
	Language   string
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
			resp += fmt.Sprintf("<%s>\n", prefix+name+strings.Repeat(" ", maxlen+1-len(name))+ctx.T("command_"+name+"_help"))
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

func GetCommand(msg string) (CommandFunction, string, []string) {

	args := strings.Fields(msg)
	if len(args) == 0 {
		return nil, "", nil
	}
	return Commands[args[0]], args[0], args[1:]

	/*
		for _, commandin := range com.Commands {

			if strings.HasPrefix(msg, commandin.Name) {
				return commandin, args[1:]
			}

		}
	*/

}

func OnMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
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

	prefix, _ := Redis.Get("prefix").Result()

	if len(prefix) > 0 {

		if strings.HasPrefix(message.Content, prefix) {
			origMessage := message.Content

			message.Content = strings.TrimPrefix(message.Content, prefix)

			command, name, args := GetCommand(message.Content)
			if command != nil {

				language, err := Redis.Get("language_" + channel.GuildID).Result()
				if err != nil {
					Redis.Set("language_"+channel.GuildID, "en-GB", 0)
					language = "en-GB"
				}

				T, _ := i18n.Tfunc(language)

				ctx := &Context{
					Content:   strings.TrimPrefix(message.Content, prefix+name),
					ChannelID: message.ChannelID,
					GuildID:   channel.GuildID,
					Type:      channel.Type,
					HasPrefix: true,
					Args:      args,
					T:         T,
					Language:  language,
				}
				if len(message.Mentions) > 0 {
					ctx.HasMention = true
				}

				//start := time.Now()
				/*ret := */
				command(session, message, ctx)
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
