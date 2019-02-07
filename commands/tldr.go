package BotCommands

import (
	"errors"
	"flag"
	"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/jonas747/discordgo"
	"github.com/kennygrant/sanitize"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getTLDR(name string, language string, variant string) string {

	tldrString := "Command not found!"
	safeName := sanitize.Path(name)
	safeLanguage := sanitize.Path(language)
	safeVariant := sanitize.Path(variant)
	log.Info(safeName, safeLanguage, safeVariant)

	pagesPath := "/tmp/tldr/pages"

	if len(safeLanguage) > 0 && safeLanguage != "." {
		pagesPath += "." + safeLanguage
	}
	if len(safeVariant) > 0 && safeVariant != "." {
		pagesPath += "/" + safeVariant
	}
	log.Info(pagesPath)

	filepath.Walk(pagesPath, func(path string, f os.FileInfo, err error) error {

		if strings.Index(path, "/"+safeName+".md") != -1 {
			str, _ := ioutil.ReadFile(path)
			tldrString = string(str)
			return errors.New("Found!")
		}
		return nil
	})

	return tldrString
}

func updateCache() {
	if _, err := os.Stat("/tmp/tldr"); os.IsNotExist(err) {
		log.Info("Cloning tldr repo.")
		_, err := git.PlainClone("/tmp/tldr", false, &git.CloneOptions{
			URL:      "https://github.com/tldr-pages/tldr",
			Progress: os.Stdout,
		})
		if err != nil {
			log.Info("Failed to clone repo.")
		}
	} else {
		r, err := git.PlainOpen("/tmp/tldr")
		if err != nil {
			log.Info("Failed to open repo")
		}
		w, _ := r.Worktree()
		w.Pull(&git.PullOptions{RemoteName: "origin"})
	}
}

func init() {
	go updateCache()
	commands.RegisterCommand("tldr", TldrCommand)
}

func TldrCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {
	flagSet := flag.NewFlagSet("tldr", flag.ContinueOnError)

	flagSet.Usage = func() {
		usageString := "```sh\n"
		usageString += "Usage: tldr [-variant variant] [-language language] [command]\n"
		flagSet.VisitAll(func(fl *flag.Flag) {
			usageString += fmt.Sprintf("  --%s %s (default: %s)\n", fl.Name, fl.Usage, fl.DefValue)
		})
		usageString += "```"
		s.ChannelMessageSend(ctx.ChannelID, usageString)
	}
	flagSet.String("variant", "", "OS variant of command.")
	flagSet.String("language", "", "Language of TLDR page.")
	flagSet.Bool("refresh", false, "Gets the latest version.")

	flagSet.Parse(ctx.Args)
	command := flagSet.Arg(0)
	variant := flagSet.Lookup("variant").Value.(flag.Getter).Get().(string)
	language := flagSet.Lookup("language").Value.(flag.Getter).Get().(string)
	log.Error(variant)
	s.ChannelMessageSend(m.ChannelID, "```md\n"+getTLDR(command, language, variant)+"\n```")
	return nil

}
