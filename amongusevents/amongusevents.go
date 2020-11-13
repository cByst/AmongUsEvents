package amongusevents

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type amongUsEventState struct {
	eventTitle               string
	eventAttendees           []string
	eventCantAttend          []string
	eventRequestedTimeChange []string
}

// CreateEvent creates new among us event and post to specified channel
func CreateEvent(session *discordgo.Session, title string, channelID string) error {
	newMessage, err := session.ChannelMessageSendEmbed(
		channelID,
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
		return errors.Wrapf(err, "Error creating initial ammong us event for channel:%s\nevent title:%s", channelID, title)
	}

	err = applyBaseReactions(session, newMessage.ChannelID, newMessage.ID)

	if err != nil {
		return errors.Wrapf(err, "Error adding base reactions to inital among us event in channel:%s\nevent title:%s", channelID, title)
	}
	return nil
}

// ReSyncEvent resyncs specified event message with current reactions state of the specified message id
func ReSyncEvent(session *discordgo.Session, message *discordgo.Message) error {
	currentState, err := extractEventState(session, message)
	if err != nil {
		return errors.Wrap(err, "Error extracting event state during resync")
	}

	err = currentState.updateEmbedMessageFromState(session, message)
	return errors.Wrap(err, "Error updating embeded message in rsync")
}

func extractEventState(session *discordgo.Session, message *discordgo.Message) (*amongUsEventState, error) {
	rsvpYes, err := session.MessageReactions(message.ChannelID, message.ID, "💯", 100, "", "")
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting message reaction info for message: %s in channel: %s reaction: 💯", message.ID, message.ChannelID)
	}

	var attendingUsers []string
	for _, user := range rsvpYes {
		if !user.Bot {
			attendingUsers = append(attendingUsers, user.Username)
		}
	}

	rsvpNo, err := session.MessageReactions(message.ChannelID, message.ID, "🙅‍♀️", 100, "", "")
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting message reaction info for message: %s in channel: %s reaction: 🙅‍♀️", message.ID, message.ChannelID)
	}

	var notAttendingUsers []string
	for _, user := range rsvpNo {
		if !user.Bot {
			notAttendingUsers = append(notAttendingUsers, user.Username)
		}
	}

	timeChangeRequested, err := session.MessageReactions(message.ChannelID, message.ID, "⏰", 100, "", "")
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting message reaction info for message: %s in channel: %s reaction: ⏰", message.ID, message.ChannelID)
	}

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
	_, err := session.ChannelMessageEditEmbed(
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

	return errors.Wrap(err, "Error updating embed message in updateEmbedMessageFromState")
}

func applyBaseReactions(session *discordgo.Session, channelID string, messageID string) error {
	err := session.MessageReactionAdd(channelID, messageID, "💯")
	if err != nil {
		return errors.Wrap(err, "Error adding base message reaction 💯")
	}
	err = session.MessageReactionAdd(channelID, messageID, "🙅‍♀️")
	if err != nil {
		return errors.Wrap(err, "Error adding base message reaction 🙅‍♀️")
	}
	err = session.MessageReactionAdd(channelID, messageID, "⏰")
	if err != nil {
		return errors.Wrap(err, "Error adding base message reaction ⏰")
	}
	return nil
}
