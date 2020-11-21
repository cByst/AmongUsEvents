package amongushandlers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cbyst/AmongUsEvents/amongusevents"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AttachHandlers will attach all among us bot handlers to the discord session
func AttachHandlers(discordSession *discordgo.Session) {
	discordSession.AddHandler(commandHandler)
	discordSession.AddHandler(messageReactionAddHandle)
	discordSession.AddHandler(messageReactionRemoveHandle)
	discordSession.AddHandler(serverBotAddHandler)
	discordSession.AddHandler(serverBotRemoveHandler)
}

func serverBotAddHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	log.Infof("Discord Server %s added the AmongUsEvents bot", g.Name)
}

func serverBotRemoveHandler(s *discordgo.Session, g *discordgo.GuildDelete) {
	log.Infof("Discord Server %s removed AmongUsEvents bot", g.Name)
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
		reactionID := m.MessageReaction.Emoji.Name
		if m.MessageReaction.Emoji.ID != "" {
			reactionID = url.QueryEscape(fmt.Sprintf("<:%s:%s", m.MessageReaction.Emoji.Name, m.MessageReaction.Emoji.ID))
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, reactionID, m.MessageReaction.UserID)
		if err != nil {
			log.Error(errors.WithMessage(err, fmt.Sprintf("Error removing unsupported reaction in message reaction add handler for %s reaction event", m.MessageReaction.Emoji.Name)))
		}
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages written by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check message for for command prefix to determine if the message is relevant to the bot
	if strings.HasPrefix(strings.ToLower(m.Content), "!createamongevent ") {
		// Check if user is privileged to command the bot
		userIsPrivledged, err := isUserPrivleged(s, m.Author.ID, m.GuildID)
		if err != nil {
			log.Error(errors.WithMessage(err, "Issue checking if user is privileged in command handler"))
		}

		// Ignore message if user is not privileged to command bot
		if !userIsPrivledged {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> You do not have access to create Among Us Events. To create events you need the amongusbot role.", m.Author.ID))
			// tell users there not permissioned for this
			return
		} else {
			title := strings.Trim(m.Content[18:len(m.Content)], "\"")

			log.Infof("Creating new among event with title: %s for user: %s", title, m.Author.Username)
			err = amongusevents.CreateEvent(s, title, m.ChannelID)
			if err != nil {
				log.Error(errors.WithMessage(err, "Error creating event in create event command handler"))
			}
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
