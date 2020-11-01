package amongusevents

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type amongUsEventState struct {
	eventTitle      string
	eventAttendees  []string
	eventCantAttend []string
}

func CreateEvent(session *discordgo.Session, title string, channelId string) error {
	newMessage, err := session.ChannelMessageSendEmbed(
		channelId,
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

	err = applyBaseReactions(session, newMessage.ChannelID, newMessage.ID)

	if err != nil {
		fmt.Printf("Error sending message err: %s", err)
	}
	return nil
}

func ReSyncEvent(session *discordgo.Session, message *discordgo.Message) error {
	return nil
}

func extractEventState(session *discordgo.Session, message *discordgo.Message) *amongUsEventState {
	// rsvpYes, _ := session.MessageReactions(message.ChannelID, message.ID, "💯", 100, "", "")
	// rsvpNo, _ := session.MessageReactions(message.ChannelID, message.ID, "💯", 100, "", "")
	return &amongUsEventState{}
}

func (s *amongUsEventState) getEmbedMessageFromState(messageId string, channelId string) error {
	return nil
}

func applyBaseReactions(session *discordgo.Session, channelID string, messageID string) error {
	err := session.MessageReactionAdd(channelID, messageID, "💯")
	if err != nil {

	}
	err = session.MessageReactionAdd(channelID, messageID, "🙅‍♀️")
	if err != nil {

	}
	return nil
}
