package bot

import (
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/NamedKitten/KittehBotGo/util/database"
	log "github.com/sirupsen/logrus"
)

// Start sets up and starts the discord connection.
func Start() {
	token := database.Get("token")
	log.Error(token)
	if len(token) <= 0 {
		log.Fatal("Token not found, please run with -runSetup to enter setup.")
		panic("")
	}
	commands.Discord.Token = "Bot " + token

	err := commands.Discord.Open()
	if err != nil {
		panic(err)
	}
}
