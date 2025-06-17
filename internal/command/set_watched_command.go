package command

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type SetWatchedCommand struct {
	MovieRepository *repository.MovieRepository
}

func NewSetWatchedCommand(movieRepository *repository.MovieRepository) *SetWatchedCommand {
	return &SetWatchedCommand{
		MovieRepository: movieRepository,
	}
}

func (c *SetWatchedCommand) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "setwatched",
		Description: "Set a movie as watched by ID",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "movie_id",
				Description: "Movie ID that will be set as watched",
				Required:    true,
			},
		},
	}
}

func (c *SetWatchedCommand) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	c.MovieRepository.UpdateWatched(movieID, true)

	response := fmt.Sprintf("You set %s as watched", input)
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
