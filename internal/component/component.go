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
	components[CustomIdSetActiveMovie] = NewSetActive(movieRepository)
	components[CustomIdCreateEvent] = NewCreateEvent(movieRepository)
	components[CustomIdCreateEventModal] = NewEventDetails(config)
	components[CustomIDSetPreferredTimeModal] = NewSetPreferredEventTime(userRepository)

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var customId string

		if i.Type == discordgo.InteractionMessageComponent {
			customId = i.MessageComponentData().CustomID
		} else if i.Type == discordgo.InteractionModalSubmit {
			customId = i.ModalSubmitData().CustomID
		} else {
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
