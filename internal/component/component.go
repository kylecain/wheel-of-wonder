package component

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/config"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type Component interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var components = map[string]Component{}

func RegisterAll(
	s *discordgo.Session,
	config *config.Config,
	movieRepository *repository.Movie,
	userRepository *repository.User,
) {
	components[CustomIdSetActiveMovie] = NewSetActive(movieRepository, userRepository)
	components[CustomIdCreateEvent] = NewCreateEvent(movieRepository)
	components[CustomIdCreateEventModal] = NewEventDetails(config)
	components[CustomIDSetPreferredTimeModal] = NewSetPreferredEventTime(userRepository)
	components[CustomIdCreateEventPreferredTime] = NewCreateEventPreferredTime(movieRepository, userRepository, config)
	components[CustomIdAnnounceMovie] = NewAnnounceMovie(movieRepository)
	components[CustomIdDeleteMovie] = NewDeleteMovie(movieRepository)

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
