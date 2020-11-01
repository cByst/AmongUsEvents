package amongushandlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cbyst/AmongUsHelper/amongusevents"
)

func AttachHandlers(discordSession *discordgo.Session) error {
	discordSession.AddHandler(commandHandler)
	discordSession.AddHandler(messageReactionAddHandle)
	return nil
}

func messageReactionAddHandle(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, _ := s.ChannelMessage(m.MessageReaction.ChannelID, m.MessageReaction.MessageID)
	for _, x := range message.Embeds {
		fmt.Printf("~~~~%+v", x)
	}

	if m.MessageReaction.UserID == s.State.User.ID || message.Author.ID != s.State.User.ID {
		return
	}

	if m.MessageReaction.Emoji.Name == "💯" {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "🙅‍♀️", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = amongusevents.ReSyncEvent(s, message)

	} else if m.MessageReaction.Emoji.Name == "🙅‍♀️" {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "💯", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = amongusevents.ReSyncEvent(s, message)
	} else {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, m.MessageReaction.Emoji.Name, m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Need to implement permissions over who can run commands on the bot

	// permission, err := s.State.UserChannelPermissions(m.Author.ID, m.ChannelID)
	// if err != nil {
	// 	log.Error(errors.WithMessage(err, "Error checking users permissions"))
	// 	permission = 0
	// }

	// if permission < 1 {
	// 	return
	// }

	if strings.HasPrefix(m.Content, "!CreateEvent") {
		title := strings.Trim(strings.TrimPrefix(m.Content, "!CreateEvent "), "\"")

		err := amongusevents.CreateEvent(s, title, m.ChannelID)
		if err != nil {
			fmt.Printf("Error sending message err: %s", err)
		}
	}

	return
}
