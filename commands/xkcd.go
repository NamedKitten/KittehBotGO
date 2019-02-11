package BotCommands

import (
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/go-errors/errors"
	"github.com/NamedKitten/discordgo"
	"github.com/nishanths/go-xkcd"
	"runtime/debug"
	"strconv"
)

var client = xkcd.NewClient()

func init() {
	commands.RegisterCommand("xkcd", xkcdCommand)
	commands.RegisterHelp("xkcd", "Shows the latest xkcd comic.")
}

func xkcdCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	defer debug.FreeOSMemory()

	var comic xkcd.Comic
	var err error

	if len(ctx.Args) > 0 {
		fmt.Println(ctx.Args)
		if ctx.Args[0] == "random" {
			comic, err = client.Random()
		} else {
			num, _ := strconv.ParseInt(ctx.Args[0], 10, 64)
			comic, err = client.Get(int(num))
		}

	} else {
		comic, err = client.Latest()
	}
	if err != nil {
		fmt.Println(err)

		return errors.Wrap(err, 1)
	}
	fmt.Println(comic)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("https://xkcd.com/%d/", comic.Number))

	return nil

}
