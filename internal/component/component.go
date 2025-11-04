package component

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/service"
)

type Component interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var components = map[string]Component{}

func RegisterAll(
	s *discordgo.Session,
	movieRepository *repository.Movie,
	userRepository *repository.User,
	movieService *service.Movie,
	logger *slog.Logger,
) {
	components[CustomIDSetPreferredTimeModal] = NewSetPreferredEventTime(userRepository, logger)
	components[CustomIdAnnounceMovie] = NewAnnounceMovie(movieRepository)
	components[CustomIdCreateEventModal] = NewEventDetails(movieRepository, movieService)
	components[CustomIdCreateEventPreferredTime] = NewCreateEventPreferredTime(movieRepository, userRepository, movieService)
	components[CustomIdCreateEvent] = NewCreateEvent(movieRepository)
	components[CustomIdDeleteMovie] = NewDeleteMovie(movieRepository)
	components[CustomIdBonusMovieModal] = NewBonusMovie(movieService, logger)

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
