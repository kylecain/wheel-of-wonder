package command

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
)

type ActiveMovie struct {
	MovieRepository *repository.Movie
}

func NewActiveMovie(movieRepository *repository.Movie) *ActiveMovie {
	return &ActiveMovie{
		MovieRepository: movieRepository,
	}
}

func (c *ActiveMovie) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        commandNameSetActive,
		Description: "Show the active movie",
	}
}

func (c *ActiveMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	activeMovie, err := c.MovieRepository.GetActive(i.GuildID)
	if err != nil || activeMovie == nil {
		InteractionResponseError(s, i, err, "Failed to retrieve active movie.")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The active movie is: %s", activeMovie.Title),
		},
	})
	if err != nil {
		slog.Error("failed to respond to add command", "error", err)
	}
}
