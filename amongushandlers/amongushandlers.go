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
	discordSession.AddHandler(messageReactionRemoveHandle)
	return nil
}

func messageReactionRemoveHandle(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	message, _ := s.ChannelMessage(m.MessageReaction.ChannelID, m.MessageReaction.MessageID)
	amongusevents.ReSyncEvent(s, message)
}

func messageReactionAddHandle(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, _ := s.ChannelMessage(m.MessageReaction.ChannelID, m.MessageReaction.MessageID)

	if m.MessageReaction.UserID == s.State.User.ID || message.Author.ID != s.State.User.ID {
		return
	}

	if m.MessageReaction.Emoji.Name == "💯" {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "🙅‍♀️", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "⏰", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = amongusevents.ReSyncEvent(s, message)

	} else if m.MessageReaction.Emoji.Name == "🙅‍♀️" {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "💯", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "⏰", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = amongusevents.ReSyncEvent(s, message)
	} else if m.MessageReaction.Emoji.Name == "⏰" {
		err := s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "💯", m.MessageReaction.UserID)
		if err != nil {
			fmt.Printf("Error removing unsupported reaction %s", err)
		}
		err = s.MessageReactionRemove(m.MessageReaction.ChannelID, m.MessageReaction.MessageID, "🙅‍♀️", m.MessageReaction.UserID)
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

	userIsPrivledged, err := isUserPrivleged(s, m.Author.ID, m.GuildID)
	if err != nil {

	}
	if !userIsPrivledged {
		return
	}

	if strings.HasPrefix(m.Content, "!CreateAmongEvent ") {
		title := strings.Trim(strings.TrimPrefix(m.Content, "!CreateAmongEvent "), "\"")

		err := amongusevents.CreateEvent(s, title, m.ChannelID)
		if err != nil {
			fmt.Printf("Error sending message err: %s", err)
		}
	}

	return
}

func isUserPrivleged(s *discordgo.Session, userId, guildId string) (bool, error) {
	amongUsRoleID, err := getAmongUsRoleID(s, guildId)
	if err != nil {
		return false, err
	}

	member, err := s.GuildMember(guildId, userId)
	if err != nil {
		return false, err
	} else {
		for _, role := range member.Roles {
			if role == amongUsRoleID {
				return true, nil
			}
		}
		return false, nil
	}
}

func getAmongUsRoleID(s *discordgo.Session, guildId string) (string, error) {
	roles, err := s.GuildRoles(guildId)
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
