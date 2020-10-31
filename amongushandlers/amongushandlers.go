package amongushandlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AttachHandlers(discordSession *discordgo.Session) error {
	discordSession.AddHandler(commandHandler)
	discordSession.AddHandler(messageReactionAddHandle)
	discordSession.AddHandler(messageReactionRemoveHandle)
	return nil
}

func messageReactionAddHandle(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Printf("discord user id: %s reacted", m.UserID)
}

func messageReactionRemoveHandle(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	fmt.Printf("discord user id: %s reacted", m.UserID)
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

		newMessage, err := s.ChannelMessageSendEmbed(
			m.ChannelID,
			&discordgo.MessageEmbed{
				Title: title,
				Color: 15105570,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Among Us Helper Bot",
					IconURL: "https://i.imgur.com/Mf4Rj0T.png",
				},
				Description: "Test For Among Us Helper",
				Footer:      &discordgo.MessageEmbedFooter{Text: "TEST Footer"},
			},
		)
		if err != nil {
			fmt.Printf("Error sending message err: %s", err)
		}
		err = addBaseReactions(s, newMessage)

		if err != nil {
			fmt.Printf("Error sending message err: %s", err)
		}
	}
}

func addBaseReactions(s *discordgo.Session, m *discordgo.Message) error {
	s.MessageReactionAdd(m.ChannelID, m.ID, "💯")
	return s.MessageReactionAdd(m.ChannelID, m.ID, "🙅‍♀️")
}
