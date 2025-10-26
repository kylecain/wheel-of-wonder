package component

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
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
		InteractionResponseError(s, i, err, "failed to convert movie ID")
		return
	}
	c.MovieRepository.DeleteMovie(int64(movieId))

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    "Movie deleted successfully.",
			Components: []discordgo.MessageComponent{},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		InteractionResponseError(s, i, err, "failed to create event modal")
		return
	}
}
