package component

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type CreateEvent struct {
	MovieRepository *repository.Movie
}

func NewCreateEvent(movieRepository *repository.Movie) *CreateEvent {
	return &CreateEvent{
		MovieRepository: movieRepository,
	}
}

func CreateEventButton(movieID, movieTitle string) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Create Event",
		CustomID: fmt.Sprintf("%s:%s:%s", CustomIdCreateEvent, movieID, movieTitle),
	}
}

func (c *CreateEvent) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIDStr, movieTitle := args[1], args[2]

	components, err := util.ScheduleEventModal(false)
	if err != nil {
		slog.Error("Failed to create event modal", "error", err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   fmt.Sprintf("%s:%s:%s", CustomIdCreateEventModal, movieIDStr, movieTitle),
			Title:      "Create Event",
			Components: components,
		},
	})
	if err != nil {
		slog.Error("Failed to create event modal", "error", err)
	}
}
