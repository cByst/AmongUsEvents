package amongushandlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cbyst/AmongUsHelper/amongusevents"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AttachHandlers will attach all among us bot handlers to the discord session
func AttachHandlers(discordSession *discordgo.Session) {
	discordSession.AddHandler(commandHandler)
	discordSession.AddHandler(messageReactionAddHandle)
	discordSession.AddHandler(messageReactionRemoveHandle)
}

func messageReactionRemoveHandle(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	message, err := s.ChannelMessage(m.MessageReaction.ChannelID, m.MessageReaction.MessageID)
	if err != nil {
		log.Error(errors.WithMessage(err, "Error finding message in message reaction remove handler"))
	}
	err = amongusevents.ReSyncEvent(s, message)
	if err != nil {
		log.Error(errors.WithMessage(err, "Error resyncing event state in reaction remove handler"))
	}
}

func messageReactionAddHandle(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(m.MessageReaction.ChannelID, m.MessageReaction.MessageID)
	if err != nil {
		log.Error(errors.WithMessage(err, "Error finding message in message reaction add handler"))
	}

	// Ignore if action was performed by the bot
	if m.MessageReaction.UserID == s.State.User.ID || message.Author.ID != s.State.User.ID {
		return
	}

	if m.MessageReaction.Emoji.Name == "💯" {
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "🙅‍♀️", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing decline reaction in message reaction add handler for accept reaction event"))
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "⏰", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing change time reaction in message reaction add handler for accept reaction event"))
		}
		err = amongusevents.ReSyncEvent(s, message)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error resyncing event state in reaction add handler for accept reaction event"))
		}

	} else if m.MessageReaction.Emoji.Name == "🙅‍♀️" {
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "💯", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing accept reaction in message reaction add handler for decline reaction event"))
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "⏰", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing change time reaction in message reaction add handler for decline reaction event"))
		}
		err = amongusevents.ReSyncEvent(s, message)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error resyncing event state in reaction add handler for decline reaction event"))
		}
	} else if m.MessageReaction.Emoji.Name == "⏰" {
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "💯", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing accept reaction in message reaction add handler for change time reaction event"))
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "🙅‍♀️", m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error removing decline reaction in message reaction add handler for change time reaction event"))
		}
		err = amongusevents.ReSyncEvent(s, message)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error resyncing event state in reaction add handler for change time reaction event"))
		}
	} else {
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, m.MessageReaction.Emoji.Name, m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, fmt.Sprintf("Error removing unsupported reactio in message reaction add handler for %s reaction event", m.MessageReaction.Emoji.Name)))
		}
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages written by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if user is privileged to command the bot
	userIsPrivledged, err := isUserPrivleged(s, m.Author.ID, m.GuildID)
	if err != nil {
		log.Error(errors.WithMessage(err, "Issue checking if user is privileged in command handler"))
	}

	// Ignore message if user is not privileged to command bot
	if !userIsPrivledged {
		return
	}

	// Check message for for command prefix to determine if the message is relevant to the bot
	if strings.HasPrefix(m.Content, "!CreateAmongEvent ") {
		title := strings.Trim(strings.TrimPrefix(m.Content, "!CreateAmongEvent "), "\"")

		err = amongusevents.CreateEvent(s, title, m.ChannelID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Error creating event in create event command handler"))
		}
	}
}

// Check if user has a role called "amongusbot" on the discord server
func isUserPrivleged(s *discordgo.Session, userID, guildID string) (bool, error) {
	amongUsRoleID, err := getAmongUsRoleID(s, guildID)
	if err != nil {
		return false, err
	}

	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	for _, role := range member.Roles {
		if role == amongUsRoleID {
			return true, nil
		}
	}

	return false, nil
}

// Look up the role id for the "amongusbot" bot role on a server
func getAmongUsRoleID(s *discordgo.Session, guildID string) (string, error) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "-1", err
	}
	for _, role := range roles {
		if "amongusbot" == strings.ToLower(role.Name) {
			return role.ID, nil
		}
	}
	return "-1", nil
}
