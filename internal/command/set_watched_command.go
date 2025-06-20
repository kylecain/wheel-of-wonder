package command

import (
	"fmt"
	"log/slog"

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
		Name:        commandNameSetWatched,
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
	input := i.ApplicationCommandData().Options[0].IntValue()

	err := c.MovieRepository.UpdateWatched(int(input), true)
	if err != nil {
		InteractionResponseError(s, i, err, "failed to update watched movie")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You set %d as watched", input),
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
