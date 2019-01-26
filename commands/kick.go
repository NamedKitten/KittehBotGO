package BotCommands

import (
	//"fmt"
	"github.com/NamedKitten/KittehBotGo/util/commands"
	"github.com/bwmarrin/discordgo"
	//"runtime/debug"
	//"strings"
	//"github.com/go-errors/errors"
)

func init() {
	commands.RegisterCommand("kick", KickCommand)
}

func KickCommand(s *discordgo.Session, m *discordgo.MessageCreate, ctx *commands.Context) error {

	// Also define the string that we will return that lists all kicked and not kicked users (pulled from above slices)
	kickedUsersStr := "Kicked successfully: "
	notKickedUsersStr := "Kicked unsuccessfully: "
	// Get author permissions in current channel
	authorPermissions, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	// Check if author has permissions to kick in that channel
	if authorPermissions&discordgo.PermissionKickMembers > 0 {
		// Check if author mentioned at least one user
		if len(m.Mentions) > 0 {
			// Save all mentioned users to 'usersToKick'
			usersToKick := m.Mentions
			// Loop through all users in usersToKick
			for _, user := range usersToKick {
				// Try to kick the user
				err := s.GuildMemberDeleteWithReason(ctx.GuildID, user.ID, "Kick requested by: "+m.Author.Username)
				// If we succedeed, we can add that user to []kickedUsers table
				// If not, we add that user to []notKickedUsers
				if err != nil {
					notKickedUsersStr += "\n" + user.String()
				} else {
					kickedUsersStr += "\n" + user.String()
				}

			}
			s.ChannelMessageSend(m.ChannelID, kickedUsersStr)
			s.ChannelMessageSend(m.ChannelID, notKickedUsersStr)
		} else {
			//If we haven't mentioned any users, we don't kick anyone - exit
			s.ChannelMessageSend(m.ChannelID, "You haven't mentioned anyone")
		}
	} else {
		// If someone tried to kick others without having KickPermission, tell him, send a warning to action-log and exit
		s.ChannelMessageSend(m.ChannelID, m.Author.Username+" is not in the sudoers file. This incident will be reported.")
		return nil
	}
	return nil
}
