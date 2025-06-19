package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SetActiveCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewSetActiveCommand(movieRepository *repository.MovieRepository) *SetActiveCommand {
	return &SetActiveCommand{
		MovieRepository: movieRepository,
	}
}

func (c *SetActiveCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSetActive,
		Description: "Set a movie as active by ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "movie_id",
				Description: "Movie ID that will be set as active",
				Required:    true,
			},
		},
	}
}

func (c *SetActiveCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].IntValue()

	err := c.MovieRepository.UpdateActive(int(input), true)
	if err != nil {
		InteractionResponseError(s, i, err, "Failed to update active movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %d as active", input),
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
