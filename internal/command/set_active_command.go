package command

import (
	"fmt"
	"log/slog"
	"strconv"

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
	input := i.ApplicationCommandData().Options[0].StringValue()
	movieID, err := strconv.Atoi(input)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Invalid movie ID: %s", input),
			},
		})
		return
	}

	err = c.MovieRepository.UpdateActive(movieID, true)
	if err != nil {
		slog.Error("failed to update active movie", "error", err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %s as active", input),
		},
	})

	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
