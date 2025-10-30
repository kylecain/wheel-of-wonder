package component

import (
	"log/slog"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type DeleteMovie struct {
	MovieRepository *repository.Movie
}

func NewDeleteMovie(movieRepository *repository.Movie) *DeleteMovie {
	return &DeleteMovie{
		MovieRepository: movieRepository,
	}
}

func (c *DeleteMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selected := i.MessageComponentData().Values[0]
	movieId, err := strconv.Atoi(selected)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to convert movie ID")
		return
	}

	err = c.MovieRepository.DeleteMovie(int64(movieId))
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to delete movie")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Movie deleted successfully.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("Failed to update user on deleted movie", "error", err)
	}
}
