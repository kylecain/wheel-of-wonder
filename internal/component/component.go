package component

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Component interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var components = map[string]Component{}

func RegisterAll(
	s *discordgo.Session,
	movieRepository *repository.Movie,
	userRepository *repository.User,
) {
	components[CustomIDSetPreferredTimeModal] = NewSetPreferredEventTime(userRepository)
	components[CustomIdAnnounceMovie] = NewAnnounceMovie(movieRepository)
	components[CustomIdCreateEventModal] = NewEventDetails()
	components[CustomIdCreateEventPreferredTime] = NewCreateEventPreferredTime(movieRepository, userRepository)
	components[CustomIdCreateEvent] = NewCreateEvent(movieRepository)
	components[CustomIdDeleteMovie] = NewDeleteMovie(movieRepository)
	components[CustomIdSetActiveMovie] = NewSetActive(movieRepository, userRepository)

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var customId string

		switch i.Type {
		case discordgo.InteractionMessageComponent:
			customId = i.MessageComponentData().CustomID
		case discordgo.InteractionModalSubmit:
			customId = i.ModalSubmitData().CustomID
		default:
			return
		}

		key := strings.Split(customId, ":")[0]
		if len(key) == 0 {
			return
		}

		if c, ok := components[key]; ok {
			c.Handler(s, i)
		}
	})
}
