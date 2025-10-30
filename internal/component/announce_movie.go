package component

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kylecain/wheel-of-wonder/internal/db/repository"
	"github.com/kylecain/wheel-of-wonder/internal/util"
)

type AnnounceMovie struct {
	MovieRepository *repository.Movie
}

func NewAnnounceMovie(movieRepository *repository.Movie) *AnnounceMovie {
	return &AnnounceMovie{
		MovieRepository: movieRepository,
	}
}

func AnnounceMovieButton(movieID, movieTitle string) discordgo.Button {
	return discordgo.Button{
		Style:    discordgo.PrimaryButton,
		Label:    "Announce Movie",
		CustomID: fmt.Sprintf("%s:%s:%s", CustomIdAnnounceMovie, movieID, movieTitle),
	}
}

func (c *AnnounceMovie) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	args := strings.Split(i.MessageComponentData().CustomID, ":")
	movieIdStr := args[1]
	movieId, err := strconv.Atoi(movieIdStr)
	if err != nil {
		util.InteractionResponseError(s, i, err, "Failed to convert movieID")
		return
	}

	selectedMovie, err := c.MovieRepository.GetMovieByID(movieId)
	if err != nil {
		util.InteractionResponseError(s, i, err, "failed to get movie by ID")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:  []*discordgo.MessageEmbed{util.MovieEmbed(selectedMovie)},
			Content: "Now Showing:",
		},
	})
	if err != nil {
		slog.Error("Failed to respond to interaction", "error", err)
	}
}
