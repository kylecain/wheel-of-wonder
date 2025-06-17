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
		Name:        "setactive",
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
		response := fmt.Sprintf("Invalid movie ID: %s", input)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
		return
	}

	c.MovieRepository.UpdateActive(movieID, true)

	response := fmt.Sprintf("You set %s as active", input)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
