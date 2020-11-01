package amongusevents

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type amongUsEventState struct {
	eventTitle               string
	eventAttendees           []string
	eventCantAttend          []string
	eventRequestedTimeChange []string
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
			Description: "\u200B\n",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{Name: "💯 **__Crew Mates__ (0) :**", Value: "\u200B\n\u200B\n", Inline: false},
				&discordgo.MessageEmbedField{Name: "🙅‍♀️ **__Can't Attend__ (0) :**", Value: "\u200B\n\u200B\n", Inline: false},
				&discordgo.MessageEmbedField{Name: "⏰ **__Requested Time Change__ (0) :**", Value: "\u200B\n\u200B\n", Inline: false},
			},
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
	currentState, _ := extractEventState(session, message)
	currentState.updateEmbedMessageFromState(session, message)
	return nil
}

func extractEventState(session *discordgo.Session, message *discordgo.Message) (*amongUsEventState, error) {
	rsvpYes, _ := session.MessageReactions(message.ChannelID, message.ID, "💯", 100, "", "")

	var attendingUsers []string
	for _, user := range rsvpYes {
		if !user.Bot {
			attendingUsers = append(attendingUsers, user.Username)
		}
	}

	rsvpNo, _ := session.MessageReactions(message.ChannelID, message.ID, "🙅‍♀️", 100, "", "")
	var notAttendingUsers []string
	for _, user := range rsvpNo {
		if !user.Bot {
			notAttendingUsers = append(notAttendingUsers, user.Username)
		}
	}

	timeChangeRequested, _ := session.MessageReactions(message.ChannelID, message.ID, "⏰", 100, "", "")
	var timeChangeRequestedUsers []string
	for _, user := range timeChangeRequested {
		if !user.Bot {
			timeChangeRequestedUsers = append(timeChangeRequestedUsers, user.Username)
		}
	}

	return &amongUsEventState{
		eventTitle:               message.Embeds[0].Title,
		eventAttendees:           attendingUsers,
		eventCantAttend:          notAttendingUsers,
		eventRequestedTimeChange: timeChangeRequestedUsers,
	}, nil
}

func (s *amongUsEventState) updateEmbedMessageFromState(session *discordgo.Session, message *discordgo.Message) error {
	var eventAttendeesText, eventCantAttendText, eventRequestedTimeChangeText string
	if len(s.eventAttendees) < 1 {
		eventAttendeesText = "\u200B\n\u200B\n"
	} else {
		for i, user := range s.eventAttendees {
			eventAttendeesText += fmt.Sprintf("\u200B    %s ``%d``\n", user, i+1)
		}
		eventAttendeesText += "\u200B\n"
	}

	if len(s.eventCantAttend) < 1 {
		eventCantAttendText = "\u200B\n\u200B\n"
	} else {
		for i, user := range s.eventCantAttend {
			eventCantAttendText += fmt.Sprintf("\u200B    %s ``%d``\n", user, i+1)
		}
		eventCantAttendText += "\u200B\n"
	}

	if len(s.eventRequestedTimeChange) < 1 {
		eventRequestedTimeChangeText = "\u200B\n\u200B\n"
	} else {
		for i, user := range s.eventRequestedTimeChange {
			eventRequestedTimeChangeText += fmt.Sprintf("\u200B    %s ``%d``\n", user, i+1)
		}
		eventRequestedTimeChangeText += "\u200B\n"
	}
	session.ChannelMessageEditEmbed(
		message.ChannelID,
		message.ID,
		&discordgo.MessageEmbed{
			Title: s.eventTitle,
			Color: 15105570,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    "Among Us Helper Bot",
				IconURL: "https://i.imgur.com/Mf4Rj0T.png",
			},
			Description: "\u200B\n",
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{Name: fmt.Sprintf("💯 **__Crew Mates__ (%d) :**", len(s.eventAttendees)), Value: eventAttendeesText, Inline: false},
				&discordgo.MessageEmbedField{Name: fmt.Sprintf("🙅‍♀️ **__Can't Attend__ (%d) :**", len(s.eventCantAttend)), Value: eventCantAttendText, Inline: false},
				&discordgo.MessageEmbedField{Name: fmt.Sprintf("⏰ **__Requested Time Change__ (%d) :**", len(s.eventRequestedTimeChange)), Value: eventRequestedTimeChangeText, Inline: false},
			},
		},
	)
	return nil
}

func applyBaseReactions(session *discordgo.Session, channelID string, messageID string) error {
	err := session.MessageReactionAdd(channelID, messageID, "💯")
	if err != nil {

	}
	err = session.MessageReactionAdd(channelID, messageID, "🙅‍♀️")
	if err != nil {

	}
	err = session.MessageReactionAdd(channelID, messageID, "⏰")
	if err != nil {

	}
	return nil
}
